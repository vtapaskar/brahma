package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type DeviceServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Unregister(context.Context, *UnregisterRequest) (*UnregisterResponse, error)
	GetDevice(context.Context, *GetDeviceRequest) (*DeviceResponse, error)
	ListDevices(context.Context, *ListDevicesRequest) (*ListDevicesResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	mustEmbedUnimplementedDeviceServiceServer()
}

type UnimplementedDeviceServiceServer struct{}

func (UnimplementedDeviceServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, nil
}
func (UnimplementedDeviceServiceServer) Unregister(context.Context, *UnregisterRequest) (*UnregisterResponse, error) {
	return nil, nil
}
func (UnimplementedDeviceServiceServer) GetDevice(context.Context, *GetDeviceRequest) (*DeviceResponse, error) {
	return nil, nil
}
func (UnimplementedDeviceServiceServer) ListDevices(context.Context, *ListDevicesRequest) (*ListDevicesResponse, error) {
	return nil, nil
}
func (UnimplementedDeviceServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, nil
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}

func RegisterDeviceServiceServer(s *grpc.Server, srv DeviceServiceServer) {
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "brahma.v1.DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _DeviceService_Register_Handler,
		},
		{
			MethodName: "Unregister",
			Handler:    _DeviceService_Unregister_Handler,
		},
		{
			MethodName: "GetDevice",
			Handler:    _DeviceService_GetDevice_Handler,
		},
		{
			MethodName: "ListDevices",
			Handler:    _DeviceService_ListDevices_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _DeviceService_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "brahma/v1/device.proto",
}

func _DeviceService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.DeviceService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Unregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnregisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Unregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.DeviceService/Unregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Unregister(ctx, req.(*UnregisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.DeviceService/GetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetDevice(ctx, req.(*GetDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_ListDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDevicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).ListDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.DeviceService/ListDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).ListDevices(ctx, req.(*ListDevicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.DeviceService/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

type MetricsServiceServer interface {
	ReportCPUStats(context.Context, *CPUStatsRequest) (*MetricsResponse, error)
	ReportProcessStats(context.Context, *ProcessStatsRequest) (*MetricsResponse, error)
	ReportMgmtNetworkStats(context.Context, *MgmtNetworkStatsRequest) (*MetricsResponse, error)
	ReportRouterBaseState(context.Context, *RouterBaseStateRequest) (*MetricsResponse, error)
	StreamMetrics(MetricsService_StreamMetricsServer) error
	mustEmbedUnimplementedMetricsServiceServer()
}

type UnimplementedMetricsServiceServer struct{}

func (UnimplementedMetricsServiceServer) ReportCPUStats(context.Context, *CPUStatsRequest) (*MetricsResponse, error) {
	return nil, nil
}
func (UnimplementedMetricsServiceServer) ReportProcessStats(context.Context, *ProcessStatsRequest) (*MetricsResponse, error) {
	return nil, nil
}
func (UnimplementedMetricsServiceServer) ReportMgmtNetworkStats(context.Context, *MgmtNetworkStatsRequest) (*MetricsResponse, error) {
	return nil, nil
}
func (UnimplementedMetricsServiceServer) ReportRouterBaseState(context.Context, *RouterBaseStateRequest) (*MetricsResponse, error) {
	return nil, nil
}
func (UnimplementedMetricsServiceServer) StreamMetrics(MetricsService_StreamMetricsServer) error {
	return nil
}
func (UnimplementedMetricsServiceServer) mustEmbedUnimplementedMetricsServiceServer() {}

type MetricsService_StreamMetricsServer interface {
	SendAndClose(*MetricsResponse) error
	Recv() (*MetricsStreamRequest, error)
	grpc.ServerStream
}

type metricsServiceStreamMetricsServer struct {
	grpc.ServerStream
}

func (x *metricsServiceStreamMetricsServer) SendAndClose(m *MetricsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *metricsServiceStreamMetricsServer) Recv() (*MetricsStreamRequest, error) {
	m := new(MetricsStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func RegisterMetricsServiceServer(s *grpc.Server, srv MetricsServiceServer) {
	s.RegisterService(&MetricsService_ServiceDesc, srv)
}

var MetricsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "brahma.v1.MetricsService",
	HandlerType: (*MetricsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReportCPUStats",
			Handler:    _MetricsService_ReportCPUStats_Handler,
		},
		{
			MethodName: "ReportProcessStats",
			Handler:    _MetricsService_ReportProcessStats_Handler,
		},
		{
			MethodName: "ReportMgmtNetworkStats",
			Handler:    _MetricsService_ReportMgmtNetworkStats_Handler,
		},
		{
			MethodName: "ReportRouterBaseState",
			Handler:    _MetricsService_ReportRouterBaseState_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamMetrics",
			Handler:       _MetricsService_StreamMetrics_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "brahma/v1/metrics.proto",
}

func _MetricsService_ReportCPUStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CPUStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).ReportCPUStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.MetricsService/ReportCPUStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).ReportCPUStats(ctx, req.(*CPUStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_ReportProcessStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).ReportProcessStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.MetricsService/ReportProcessStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).ReportProcessStats(ctx, req.(*ProcessStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_ReportMgmtNetworkStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MgmtNetworkStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).ReportMgmtNetworkStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.MetricsService/ReportMgmtNetworkStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).ReportMgmtNetworkStats(ctx, req.(*MgmtNetworkStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_ReportRouterBaseState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouterBaseStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).ReportRouterBaseState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.MetricsService/ReportRouterBaseState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).ReportRouterBaseState(ctx, req.(*RouterBaseStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_StreamMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MetricsServiceServer).StreamMetrics(&metricsServiceStreamMetricsServer{stream})
}

type LogServiceServer interface {
	UploadCrashReport(LogService_UploadCrashReportServer) error
	UploadBacktrace(LogService_UploadBacktraceServer) error
	GetLogMetadata(context.Context, *GetLogMetadataRequest) (*LogMetadataResponse, error)
	ListLogs(context.Context, *ListLogsRequest) (*ListLogsResponse, error)
	mustEmbedUnimplementedLogServiceServer()
}

type UnimplementedLogServiceServer struct{}

func (UnimplementedLogServiceServer) UploadCrashReport(LogService_UploadCrashReportServer) error {
	return nil
}
func (UnimplementedLogServiceServer) UploadBacktrace(LogService_UploadBacktraceServer) error {
	return nil
}
func (UnimplementedLogServiceServer) GetLogMetadata(context.Context, *GetLogMetadataRequest) (*LogMetadataResponse, error) {
	return nil, nil
}
func (UnimplementedLogServiceServer) ListLogs(context.Context, *ListLogsRequest) (*ListLogsResponse, error) {
	return nil, nil
}
func (UnimplementedLogServiceServer) mustEmbedUnimplementedLogServiceServer() {}

type LogService_UploadCrashReportServer interface {
	SendAndClose(*LogUploadResponse) error
	Recv() (*CrashReportChunk, error)
	grpc.ServerStream
}

type logServiceUploadCrashReportServer struct {
	grpc.ServerStream
}

func (x *logServiceUploadCrashReportServer) SendAndClose(m *LogUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *logServiceUploadCrashReportServer) Recv() (*CrashReportChunk, error) {
	m := new(CrashReportChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type LogService_UploadBacktraceServer interface {
	SendAndClose(*LogUploadResponse) error
	Recv() (*BacktraceChunk, error)
	grpc.ServerStream
}

type logServiceUploadBacktraceServer struct {
	grpc.ServerStream
}

func (x *logServiceUploadBacktraceServer) SendAndClose(m *LogUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *logServiceUploadBacktraceServer) Recv() (*BacktraceChunk, error) {
	m := new(BacktraceChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func RegisterLogServiceServer(s *grpc.Server, srv LogServiceServer) {
	s.RegisterService(&LogService_ServiceDesc, srv)
}

var LogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "brahma.v1.LogService",
	HandlerType: (*LogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLogMetadata",
			Handler:    _LogService_GetLogMetadata_Handler,
		},
		{
			MethodName: "ListLogs",
			Handler:    _LogService_ListLogs_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadCrashReport",
			Handler:       _LogService_UploadCrashReport_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadBacktrace",
			Handler:       _LogService_UploadBacktrace_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "brahma/v1/logs.proto",
}

func _LogService_UploadCrashReport_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LogServiceServer).UploadCrashReport(&logServiceUploadCrashReportServer{stream})
}

func _LogService_UploadBacktrace_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LogServiceServer).UploadBacktrace(&logServiceUploadBacktraceServer{stream})
}

func _LogService_GetLogMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).GetLogMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.LogService/GetLogMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).GetLogMetadata(ctx, req.(*GetLogMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogService_ListLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).ListLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brahma.v1.LogService/ListLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).ListLogs(ctx, req.(*ListLogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}
