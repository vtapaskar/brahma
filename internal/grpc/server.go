package grpc

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/vtapaskar/brahma/internal/config"
	"github.com/vtapaskar/brahma/internal/metrics"
	"github.com/vtapaskar/brahma/internal/registry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	config    config.GRPCConfig
	collector *metrics.Collector
	registry  *registry.Registry
	logger    *zap.Logger
	server    *grpc.Server
	UnimplementedDeviceServiceServer
	UnimplementedMetricsServiceServer
	UnimplementedLogServiceServer
}

func NewServer(cfg config.GRPCConfig, collector *metrics.Collector, reg *registry.Registry, logger *zap.Logger) *Server {
	s := &Server{
		config:    cfg,
		collector: collector,
		registry:  reg,
		logger:    logger,
	}

	opts := []grpc.ServerOption{}
	s.server = grpc.NewServer(opts...)

	RegisterDeviceServiceServer(s.server, s)
	RegisterMetricsServiceServer(s.server, s)
	RegisterLogServiceServer(s.server, s)

	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info("gRPC server starting", zap.String("address", addr))
	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

func (s *Server) validateUID(uid string) bool {
	_, exists := s.registry.GetByUID(uid)
	return exists
}

func (s *Server) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	if req.ForeignKey == "" {
		return nil, status.Error(codes.InvalidArgument, "foreign_key is required")
	}

	regReq := registry.RegistrationRequest{
		ForeignKey: req.ForeignKey,
		Hostname:   req.Hostname,
		IPAddress:  req.IpAddress,
		DeviceType: req.DeviceType,
		Platform:   req.Platform,
		Version:    req.Version,
		Labels:     req.Labels,
	}

	device, err := s.registry.Register(regReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register device: %v", err)
	}

	return &RegisterResponse{
		Uid:          device.UID,
		ForeignKey:   device.ForeignKey,
		Status:       "registered",
		RegisteredAt: timestamppb.New(device.RegisteredAt),
	}, nil
}

func (s *Server) Unregister(ctx context.Context, req *UnregisterRequest) (*UnregisterResponse, error) {
	if !s.registry.Unregister(req.Uid) {
		return nil, status.Error(codes.NotFound, "device not found")
	}

	return &UnregisterResponse{
		Status: "unregistered",
		Uid:    req.Uid,
	}, nil
}

func (s *Server) GetDevice(ctx context.Context, req *GetDeviceRequest) (*DeviceResponse, error) {
	device, exists := s.registry.GetByUID(req.Uid)
	if !exists {
		return nil, status.Error(codes.NotFound, "device not found")
	}

	return &DeviceResponse{
		Uid:          device.UID,
		ForeignKey:   device.ForeignKey,
		Hostname:     device.Hostname,
		IpAddress:    device.IPAddress,
		DeviceType:   device.DeviceType,
		Platform:     device.Platform,
		Version:      device.Version,
		Labels:       device.Labels,
		RegisteredAt: timestamppb.New(device.RegisteredAt),
		LastSeen:     timestamppb.New(device.LastSeen),
	}, nil
}

func (s *Server) ListDevices(ctx context.Context, req *ListDevicesRequest) (*ListDevicesResponse, error) {
	devices := s.registry.ListDevices()

	resp := &ListDevicesResponse{
		Total: int32(len(devices)),
	}

	for _, d := range devices {
		resp.Devices = append(resp.Devices, &DeviceResponse{
			Uid:          d.UID,
			ForeignKey:   d.ForeignKey,
			Hostname:     d.Hostname,
			IpAddress:    d.IPAddress,
			DeviceType:   d.DeviceType,
			Platform:     d.Platform,
			Version:      d.Version,
			Labels:       d.Labels,
			RegisteredAt: timestamppb.New(d.RegisteredAt),
			LastSeen:     timestamppb.New(d.LastSeen),
		})
	}

	return resp, nil
}

func (s *Server) Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error) {
	if !s.validateUID(req.Uid) {
		return nil, status.Error(codes.NotFound, "device not registered")
	}

	s.registry.UpdateLastSeen(req.Uid)

	return &HeartbeatResponse{
		Status:     "ok",
		ServerTime: timestamppb.Now(),
	}, nil
}

func (s *Server) ReportCPUStats(ctx context.Context, req *CPUStatsRequest) (*MetricsResponse, error) {
	if !s.validateUID(req.Uid) {
		return nil, status.Error(codes.NotFound, "device not registered")
	}

	stats := &metrics.CPUStats{
		UID:           req.Uid,
		UsagePercent:  req.UsagePercent,
		UserPercent:   req.UserPercent,
		SystemPercent: req.SystemPercent,
		IdlePercent:   req.IdlePercent,
		IOWaitPercent: req.IowaitPercent,
		LoadAvg1Min:   req.LoadAvg_1Min,
		LoadAvg5Min:   req.LoadAvg_5Min,
		LoadAvg15Min:  req.LoadAvg_15Min,
		NumCores:      int(req.NumCores),
		PerCoreUsage:  req.PerCoreUsage,
	}

	if err := s.collector.CollectCPUStats(stats); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to collect CPU stats: %v", err)
	}

	s.registry.UpdateLastSeen(req.Uid)

	return &MetricsResponse{
		Status: "accepted",
		Uid:    req.Uid,
	}, nil
}

func (s *Server) ReportProcessStats(ctx context.Context, req *ProcessStatsRequest) (*MetricsResponse, error) {
	if !s.validateUID(req.Uid) {
		return nil, status.Error(codes.NotFound, "device not registered")
	}

	stats := &metrics.ProcessStats{
		UID:           req.Uid,
		TotalCount:    int(req.TotalCount),
		RunningCount:  int(req.RunningCount),
		SleepingCount: int(req.SleepingCount),
		ZombieCount:   int(req.ZombieCount),
	}

	for _, p := range req.Processes {
		stats.Processes = append(stats.Processes, metrics.ProcessInfo{
			PID:        int(p.Pid),
			Name:       p.Name,
			State:      p.State,
			CPUPercent: p.CpuPercent,
			MemPercent: p.MemPercent,
			MemoryRSS:  p.MemoryRss,
			Threads:    int(p.Threads),
			StartTime:  p.StartTime,
			Command:    p.Command,
		})
	}

	if err := s.collector.CollectProcessStats(stats); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to collect process stats: %v", err)
	}

	s.registry.UpdateLastSeen(req.Uid)

	return &MetricsResponse{
		Status: "accepted",
		Uid:    req.Uid,
	}, nil
}

func (s *Server) ReportMgmtNetworkStats(ctx context.Context, req *MgmtNetworkStatsRequest) (*MetricsResponse, error) {
	if !s.validateUID(req.Uid) {
		return nil, status.Error(codes.NotFound, "device not registered")
	}

	stats := &metrics.MgmtNetworkStats{
		UID:           req.Uid,
		InterfaceName: req.InterfaceName,
		Status:        req.Status,
		IPAddress:     req.IpAddress,
		Netmask:       req.Netmask,
		Gateway:       req.Gateway,
		MACAddress:    req.MacAddress,
		Speed:         req.Speed,
		Duplex:        req.Duplex,
		RxBytes:       req.RxBytes,
		TxBytes:       req.TxBytes,
		RxPackets:     req.RxPackets,
		TxPackets:     req.TxPackets,
		RxErrors:      req.RxErrors,
		TxErrors:      req.TxErrors,
		RxDropped:     req.RxDropped,
		TxDropped:     req.TxDropped,
		DNSServers:    req.DnsServers,
		NTPServers:    req.NtpServers,
	}

	if err := s.collector.CollectMgmtNetworkStats(stats); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to collect mgmt network stats: %v", err)
	}

	s.registry.UpdateLastSeen(req.Uid)

	return &MetricsResponse{
		Status: "accepted",
		Uid:    req.Uid,
	}, nil
}

func (s *Server) ReportRouterBaseState(ctx context.Context, req *RouterBaseStateRequest) (*MetricsResponse, error) {
	if !s.validateUID(req.Uid) {
		return nil, status.Error(codes.NotFound, "device not registered")
	}

	state := &metrics.RouterBaseState{
		UID:             req.Uid,
		Hostname:        req.Hostname,
		Platform:        req.Platform,
		HardwareVersion: req.HardwareVersion,
		SoftwareVersion: req.SoftwareVersion,
		SONiCVersion:    req.SonicVersion,
		KernelVersion:   req.KernelVersion,
		UptimeSeconds:   req.UptimeSeconds,
		SerialNumber:    req.SerialNumber,
		MemoryTotal:     req.MemoryTotal,
		MemoryUsed:      req.MemoryUsed,
		MemoryFree:      req.MemoryFree,
		DiskTotal:       req.DiskTotal,
		DiskUsed:        req.DiskUsed,
		DiskFree:        req.DiskFree,
	}

	if req.BootTime != nil {
		state.BootTime = req.BootTime.AsTime()
	}

	if req.MgmtNetworkStatus != nil {
		state.MgmtNetworkStatus = metrics.MgmtStatus{
			Status:    req.MgmtNetworkStatus.Status,
			IPAddress: req.MgmtNetworkStatus.IpAddress,
			Gateway:   req.MgmtNetworkStatus.Gateway,
			Reachable: req.MgmtNetworkStatus.Reachable,
		}
	}

	if req.DnsStatus != nil {
		state.DNSStatus = metrics.DNSStatus{
			Enabled:     req.DnsStatus.Enabled,
			Servers:     req.DnsStatus.Servers,
			Domain:      req.DnsStatus.Domain,
			SearchList:  req.DnsStatus.SearchList,
			Operational: req.DnsStatus.Operational,
		}
	}

	if req.DhcpStatus != nil {
		state.DHCPStatus = metrics.DHCPStatus{
			Enabled:   req.DhcpStatus.Enabled,
			State:     req.DhcpStatus.State,
			LeaseTime: req.DhcpStatus.LeaseTime,
			RenewTime: req.DhcpStatus.RenewTime,
			ServerIP:  req.DhcpStatus.ServerIp,
		}
	}

	if req.LldpStatus != nil {
		state.LLDPStatus = metrics.LLDPStatus{
			Enabled:       req.LldpStatus.Enabled,
			ChassisID:     req.LldpStatus.ChassisId,
			SystemName:    req.LldpStatus.SystemName,
			NeighborCount: int(req.LldpStatus.NeighborCount),
		}
		for _, n := range req.LldpStatus.Neighbors {
			state.LLDPStatus.Neighbors = append(state.LLDPStatus.Neighbors, metrics.LLDPNeighbor{
				LocalPort:        n.LocalPort,
				RemoteChassisID:  n.RemoteChassisId,
				RemotePortID:     n.RemotePortId,
				RemoteSystemName: n.RemoteSystemName,
				RemotePortDesc:   n.RemotePortDesc,
				TTL:              int(n.Ttl),
			})
		}
	}

	if err := s.collector.CollectRouterBaseState(state); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to collect router base state: %v", err)
	}

	s.registry.UpdateLastSeen(req.Uid)

	return &MetricsResponse{
		Status: "accepted",
		Uid:    req.Uid,
	}, nil
}

func (s *Server) StreamMetrics(stream MetricsService_StreamMetricsServer) error {
	var uid string
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&MetricsResponse{
				Status:  "accepted",
				Uid:     uid,
				Message: fmt.Sprintf("processed %d metrics", count),
			})
		}
		if err != nil {
			return err
		}

		switch m := req.Metrics.(type) {
		case *MetricsStreamRequest_CpuStats:
			uid = m.CpuStats.Uid
			if _, err := s.ReportCPUStats(stream.Context(), m.CpuStats); err != nil {
				s.logger.Warn("Failed to process CPU stats in stream", zap.Error(err))
			}
		case *MetricsStreamRequest_ProcessStats:
			uid = m.ProcessStats.Uid
			if _, err := s.ReportProcessStats(stream.Context(), m.ProcessStats); err != nil {
				s.logger.Warn("Failed to process process stats in stream", zap.Error(err))
			}
		case *MetricsStreamRequest_MgmtNetworkStats:
			uid = m.MgmtNetworkStats.Uid
			if _, err := s.ReportMgmtNetworkStats(stream.Context(), m.MgmtNetworkStats); err != nil {
				s.logger.Warn("Failed to process mgmt network stats in stream", zap.Error(err))
			}
		case *MetricsStreamRequest_RouterBaseState:
			uid = m.RouterBaseState.Uid
			if _, err := s.ReportRouterBaseState(stream.Context(), m.RouterBaseState); err != nil {
				s.logger.Warn("Failed to process router base state in stream", zap.Error(err))
			}
		}
		count++
	}
}

func (s *Server) UploadCrashReport(stream LogService_UploadCrashReportServer) error {
	var metadata *CrashReportMetadata
	var content []byte

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch data := chunk.Data.(type) {
		case *CrashReportChunk_Metadata:
			metadata = data.Metadata
		case *CrashReportChunk_Chunk:
			content = append(content, data.Chunk...)
		}
	}

	if metadata == nil {
		return status.Error(codes.InvalidArgument, "metadata is required")
	}

	if !s.validateUID(metadata.Uid) {
		return status.Error(codes.NotFound, "device not registered")
	}

	report := &metrics.LogReport{
		DeviceUID:  metadata.Uid,
		ProcessTag: metadata.ProcessTag,
		Version:    metadata.Version,
		Filename:   metadata.Filename,
		Content:    content,
	}

	logID, err := s.collector.CollectCrashReport(report)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to store crash report: %v", err)
	}

	s.registry.UpdateLastSeen(metadata.Uid)

	return stream.SendAndClose(&LogUploadResponse{
		Status: "created",
		LogId:  logID,
		Uid:    metadata.Uid,
		S3Key:  report.S3Key,
	})
}

func (s *Server) UploadBacktrace(stream LogService_UploadBacktraceServer) error {
	var metadata *BacktraceMetadata
	var content []byte

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch data := chunk.Data.(type) {
		case *BacktraceChunk_Metadata:
			metadata = data.Metadata
		case *BacktraceChunk_Chunk:
			content = append(content, data.Chunk...)
		}
	}

	if metadata == nil {
		return status.Error(codes.InvalidArgument, "metadata is required")
	}

	if !s.validateUID(metadata.Uid) {
		return status.Error(codes.NotFound, "device not registered")
	}

	report := &metrics.LogReport{
		DeviceUID:  metadata.Uid,
		ProcessTag: metadata.ProcessTag,
		Version:    metadata.Version,
		Filename:   metadata.Filename,
		Content:    content,
	}

	logID, err := s.collector.CollectBacktrace(report)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to store backtrace: %v", err)
	}

	s.registry.UpdateLastSeen(metadata.Uid)

	return stream.SendAndClose(&LogUploadResponse{
		Status: "created",
		LogId:  logID,
		Uid:    metadata.Uid,
		S3Key:  report.S3Key,
	})
}

func (s *Server) GetLogMetadata(ctx context.Context, req *GetLogMetadataRequest) (*LogMetadataResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *Server) ListLogs(ctx context.Context, req *ListLogsRequest) (*ListLogsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
