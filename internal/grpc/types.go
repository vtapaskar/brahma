package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RegisterRequest struct {
	ForeignKey string            `protobuf:"bytes,1,opt,name=foreign_key,json=foreignKey,proto3" json:"foreign_key,omitempty"`
	Hostname   string            `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	IpAddress  string            `protobuf:"bytes,3,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	DeviceType string            `protobuf:"bytes,4,opt,name=device_type,json=deviceType,proto3" json:"device_type,omitempty"`
	Platform   string            `protobuf:"bytes,5,opt,name=platform,proto3" json:"platform,omitempty"`
	Version    string            `protobuf:"bytes,6,opt,name=version,proto3" json:"version,omitempty"`
	Labels     map[string]string `protobuf:"bytes,7,rep,name=labels,proto3" json:"labels,omitempty"`
}

type RegisterResponse struct {
	Uid          string                 `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	ForeignKey   string                 `protobuf:"bytes,2,opt,name=foreign_key,json=foreignKey,proto3" json:"foreign_key,omitempty"`
	Status       string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	RegisteredAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=registered_at,json=registeredAt,proto3" json:"registered_at,omitempty"`
}

type UnregisterRequest struct {
	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

type UnregisterResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Uid    string `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
}

type GetDeviceRequest struct {
	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

type DeviceResponse struct {
	Uid          string                 `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	ForeignKey   string                 `protobuf:"bytes,2,opt,name=foreign_key,json=foreignKey,proto3" json:"foreign_key,omitempty"`
	Hostname     string                 `protobuf:"bytes,3,opt,name=hostname,proto3" json:"hostname,omitempty"`
	IpAddress    string                 `protobuf:"bytes,4,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	DeviceType   string                 `protobuf:"bytes,5,opt,name=device_type,json=deviceType,proto3" json:"device_type,omitempty"`
	Platform     string                 `protobuf:"bytes,6,opt,name=platform,proto3" json:"platform,omitempty"`
	Version      string                 `protobuf:"bytes,7,opt,name=version,proto3" json:"version,omitempty"`
	Labels       map[string]string      `protobuf:"bytes,8,rep,name=labels,proto3" json:"labels,omitempty"`
	RegisteredAt *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=registered_at,json=registeredAt,proto3" json:"registered_at,omitempty"`
	LastSeen     *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=last_seen,json=lastSeen,proto3" json:"last_seen,omitempty"`
}

type ListDevicesRequest struct {
	Limit  int32 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset int32 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
}

type ListDevicesResponse struct {
	Devices []*DeviceResponse `protobuf:"bytes,1,rep,name=devices,proto3" json:"devices,omitempty"`
	Total   int32             `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

type HeartbeatRequest struct {
	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

type HeartbeatResponse struct {
	Status     string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	ServerTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=server_time,json=serverTime,proto3" json:"server_time,omitempty"`
}

type MetricsResponse struct {
	Status  string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Uid     string `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

type CPUStatsRequest struct {
	Uid           string    `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	UsagePercent  float64   `protobuf:"fixed64,2,opt,name=usage_percent,json=usagePercent,proto3" json:"usage_percent,omitempty"`
	UserPercent   float64   `protobuf:"fixed64,3,opt,name=user_percent,json=userPercent,proto3" json:"user_percent,omitempty"`
	SystemPercent float64   `protobuf:"fixed64,4,opt,name=system_percent,json=systemPercent,proto3" json:"system_percent,omitempty"`
	IdlePercent   float64   `protobuf:"fixed64,5,opt,name=idle_percent,json=idlePercent,proto3" json:"idle_percent,omitempty"`
	IowaitPercent float64   `protobuf:"fixed64,6,opt,name=iowait_percent,json=iowaitPercent,proto3" json:"iowait_percent,omitempty"`
	LoadAvg_1Min  float64   `protobuf:"fixed64,7,opt,name=load_avg_1min,json=loadAvg1min,proto3" json:"load_avg_1min,omitempty"`
	LoadAvg_5Min  float64   `protobuf:"fixed64,8,opt,name=load_avg_5min,json=loadAvg5min,proto3" json:"load_avg_5min,omitempty"`
	LoadAvg_15Min float64   `protobuf:"fixed64,9,opt,name=load_avg_15min,json=loadAvg15min,proto3" json:"load_avg_15min,omitempty"`
	NumCores      int32     `protobuf:"varint,10,opt,name=num_cores,json=numCores,proto3" json:"num_cores,omitempty"`
	PerCoreUsage  []float64 `protobuf:"fixed64,11,rep,packed,name=per_core_usage,json=perCoreUsage,proto3" json:"per_core_usage,omitempty"`
}

type ProcessStatsRequest struct {
	Uid           string         `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	TotalCount    int32          `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	RunningCount  int32          `protobuf:"varint,3,opt,name=running_count,json=runningCount,proto3" json:"running_count,omitempty"`
	SleepingCount int32          `protobuf:"varint,4,opt,name=sleeping_count,json=sleepingCount,proto3" json:"sleeping_count,omitempty"`
	ZombieCount   int32          `protobuf:"varint,5,opt,name=zombie_count,json=zombieCount,proto3" json:"zombie_count,omitempty"`
	Processes     []*ProcessInfo `protobuf:"bytes,6,rep,name=processes,proto3" json:"processes,omitempty"`
}

type ProcessInfo struct {
	Pid        int32   `protobuf:"varint,1,opt,name=pid,proto3" json:"pid,omitempty"`
	Name       string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	State      string  `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
	CpuPercent float64 `protobuf:"fixed64,4,opt,name=cpu_percent,json=cpuPercent,proto3" json:"cpu_percent,omitempty"`
	MemPercent float64 `protobuf:"fixed64,5,opt,name=mem_percent,json=memPercent,proto3" json:"mem_percent,omitempty"`
	MemoryRss  uint64  `protobuf:"varint,6,opt,name=memory_rss,json=memoryRss,proto3" json:"memory_rss,omitempty"`
	Threads    int32   `protobuf:"varint,7,opt,name=threads,proto3" json:"threads,omitempty"`
	StartTime  int64   `protobuf:"varint,8,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	Command    string  `protobuf:"bytes,9,opt,name=command,proto3" json:"command,omitempty"`
}

type MgmtNetworkStatsRequest struct {
	Uid           string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	InterfaceName string   `protobuf:"bytes,2,opt,name=interface_name,json=interfaceName,proto3" json:"interface_name,omitempty"`
	Status        string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	IpAddress     string   `protobuf:"bytes,4,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	Netmask       string   `protobuf:"bytes,5,opt,name=netmask,proto3" json:"netmask,omitempty"`
	Gateway       string   `protobuf:"bytes,6,opt,name=gateway,proto3" json:"gateway,omitempty"`
	MacAddress    string   `protobuf:"bytes,7,opt,name=mac_address,json=macAddress,proto3" json:"mac_address,omitempty"`
	Speed         int64    `protobuf:"varint,8,opt,name=speed,proto3" json:"speed,omitempty"`
	Duplex        string   `protobuf:"bytes,9,opt,name=duplex,proto3" json:"duplex,omitempty"`
	RxBytes       uint64   `protobuf:"varint,10,opt,name=rx_bytes,json=rxBytes,proto3" json:"rx_bytes,omitempty"`
	TxBytes       uint64   `protobuf:"varint,11,opt,name=tx_bytes,json=txBytes,proto3" json:"tx_bytes,omitempty"`
	RxPackets     uint64   `protobuf:"varint,12,opt,name=rx_packets,json=rxPackets,proto3" json:"rx_packets,omitempty"`
	TxPackets     uint64   `protobuf:"varint,13,opt,name=tx_packets,json=txPackets,proto3" json:"tx_packets,omitempty"`
	RxErrors      uint64   `protobuf:"varint,14,opt,name=rx_errors,json=rxErrors,proto3" json:"rx_errors,omitempty"`
	TxErrors      uint64   `protobuf:"varint,15,opt,name=tx_errors,json=txErrors,proto3" json:"tx_errors,omitempty"`
	RxDropped     uint64   `protobuf:"varint,16,opt,name=rx_dropped,json=rxDropped,proto3" json:"rx_dropped,omitempty"`
	TxDropped     uint64   `protobuf:"varint,17,opt,name=tx_dropped,json=txDropped,proto3" json:"tx_dropped,omitempty"`
	DnsServers    []string `protobuf:"bytes,18,rep,name=dns_servers,json=dnsServers,proto3" json:"dns_servers,omitempty"`
	NtpServers    []string `protobuf:"bytes,19,rep,name=ntp_servers,json=ntpServers,proto3" json:"ntp_servers,omitempty"`
}

type RouterBaseStateRequest struct {
	Uid               string                 `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Hostname          string                 `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Platform          string                 `protobuf:"bytes,3,opt,name=platform,proto3" json:"platform,omitempty"`
	HardwareVersion   string                 `protobuf:"bytes,4,opt,name=hardware_version,json=hardwareVersion,proto3" json:"hardware_version,omitempty"`
	SoftwareVersion   string                 `protobuf:"bytes,5,opt,name=software_version,json=softwareVersion,proto3" json:"software_version,omitempty"`
	SonicVersion      string                 `protobuf:"bytes,6,opt,name=sonic_version,json=sonicVersion,proto3" json:"sonic_version,omitempty"`
	KernelVersion     string                 `protobuf:"bytes,7,opt,name=kernel_version,json=kernelVersion,proto3" json:"kernel_version,omitempty"`
	UptimeSeconds     int64                  `protobuf:"varint,8,opt,name=uptime_seconds,json=uptimeSeconds,proto3" json:"uptime_seconds,omitempty"`
	BootTime          *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=boot_time,json=bootTime,proto3" json:"boot_time,omitempty"`
	SerialNumber      string                 `protobuf:"bytes,10,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	MgmtNetworkStatus *MgmtStatus            `protobuf:"bytes,11,opt,name=mgmt_network_status,json=mgmtNetworkStatus,proto3" json:"mgmt_network_status,omitempty"`
	DnsStatus         *DNSStatus             `protobuf:"bytes,12,opt,name=dns_status,json=dnsStatus,proto3" json:"dns_status,omitempty"`
	DhcpStatus        *DHCPStatus            `protobuf:"bytes,13,opt,name=dhcp_status,json=dhcpStatus,proto3" json:"dhcp_status,omitempty"`
	LldpStatus        *LLDPStatus            `protobuf:"bytes,14,opt,name=lldp_status,json=lldpStatus,proto3" json:"lldp_status,omitempty"`
	MemoryTotal       uint64                 `protobuf:"varint,15,opt,name=memory_total,json=memoryTotal,proto3" json:"memory_total,omitempty"`
	MemoryUsed        uint64                 `protobuf:"varint,16,opt,name=memory_used,json=memoryUsed,proto3" json:"memory_used,omitempty"`
	MemoryFree        uint64                 `protobuf:"varint,17,opt,name=memory_free,json=memoryFree,proto3" json:"memory_free,omitempty"`
	DiskTotal         uint64                 `protobuf:"varint,18,opt,name=disk_total,json=diskTotal,proto3" json:"disk_total,omitempty"`
	DiskUsed          uint64                 `protobuf:"varint,19,opt,name=disk_used,json=diskUsed,proto3" json:"disk_used,omitempty"`
	DiskFree          uint64                 `protobuf:"varint,20,opt,name=disk_free,json=diskFree,proto3" json:"disk_free,omitempty"`
}

type MgmtStatus struct {
	Status    string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	IpAddress string `protobuf:"bytes,2,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	Gateway   string `protobuf:"bytes,3,opt,name=gateway,proto3" json:"gateway,omitempty"`
	Reachable bool   `protobuf:"varint,4,opt,name=reachable,proto3" json:"reachable,omitempty"`
}

type DNSStatus struct {
	Enabled     bool     `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Servers     []string `protobuf:"bytes,2,rep,name=servers,proto3" json:"servers,omitempty"`
	Domain      string   `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	SearchList  []string `protobuf:"bytes,4,rep,name=search_list,json=searchList,proto3" json:"search_list,omitempty"`
	Operational bool     `protobuf:"varint,5,opt,name=operational,proto3" json:"operational,omitempty"`
}

type DHCPStatus struct {
	Enabled   bool   `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	State     string `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	LeaseTime int64  `protobuf:"varint,3,opt,name=lease_time,json=leaseTime,proto3" json:"lease_time,omitempty"`
	RenewTime int64  `protobuf:"varint,4,opt,name=renew_time,json=renewTime,proto3" json:"renew_time,omitempty"`
	ServerIp  string `protobuf:"bytes,5,opt,name=server_ip,json=serverIp,proto3" json:"server_ip,omitempty"`
}

type LLDPStatus struct {
	Enabled       bool            `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	ChassisId     string          `protobuf:"bytes,2,opt,name=chassis_id,json=chassisId,proto3" json:"chassis_id,omitempty"`
	SystemName    string          `protobuf:"bytes,3,opt,name=system_name,json=systemName,proto3" json:"system_name,omitempty"`
	NeighborCount int32           `protobuf:"varint,4,opt,name=neighbor_count,json=neighborCount,proto3" json:"neighbor_count,omitempty"`
	Neighbors     []*LLDPNeighbor `protobuf:"bytes,5,rep,name=neighbors,proto3" json:"neighbors,omitempty"`
}

type LLDPNeighbor struct {
	LocalPort        string `protobuf:"bytes,1,opt,name=local_port,json=localPort,proto3" json:"local_port,omitempty"`
	RemoteChassisId  string `protobuf:"bytes,2,opt,name=remote_chassis_id,json=remoteChassisId,proto3" json:"remote_chassis_id,omitempty"`
	RemotePortId     string `protobuf:"bytes,3,opt,name=remote_port_id,json=remotePortId,proto3" json:"remote_port_id,omitempty"`
	RemoteSystemName string `protobuf:"bytes,4,opt,name=remote_system_name,json=remoteSystemName,proto3" json:"remote_system_name,omitempty"`
	RemotePortDesc   string `protobuf:"bytes,5,opt,name=remote_port_desc,json=remotePortDesc,proto3" json:"remote_port_desc,omitempty"`
	Ttl              int32  `protobuf:"varint,6,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

type MetricsStreamRequest struct {
	Metrics isMetricsStreamRequest_Metrics
}

type isMetricsStreamRequest_Metrics interface {
	isMetricsStreamRequest_Metrics()
}

type MetricsStreamRequest_CpuStats struct {
	CpuStats *CPUStatsRequest
}

type MetricsStreamRequest_ProcessStats struct {
	ProcessStats *ProcessStatsRequest
}

type MetricsStreamRequest_MgmtNetworkStats struct {
	MgmtNetworkStats *MgmtNetworkStatsRequest
}

type MetricsStreamRequest_RouterBaseState struct {
	RouterBaseState *RouterBaseStateRequest
}

func (*MetricsStreamRequest_CpuStats) isMetricsStreamRequest_Metrics()        {}
func (*MetricsStreamRequest_ProcessStats) isMetricsStreamRequest_Metrics()    {}
func (*MetricsStreamRequest_MgmtNetworkStats) isMetricsStreamRequest_Metrics() {}
func (*MetricsStreamRequest_RouterBaseState) isMetricsStreamRequest_Metrics() {}

type CrashReportChunk struct {
	Data isCrashReportChunk_Data
}

type isCrashReportChunk_Data interface {
	isCrashReportChunk_Data()
}

type CrashReportChunk_Metadata struct {
	Metadata *CrashReportMetadata
}

type CrashReportChunk_Chunk struct {
	Chunk []byte
}

func (*CrashReportChunk_Metadata) isCrashReportChunk_Data() {}
func (*CrashReportChunk_Chunk) isCrashReportChunk_Data()    {}

type CrashReportMetadata struct {
	Uid        string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	ProcessTag string `protobuf:"bytes,2,opt,name=process_tag,json=processTag,proto3" json:"process_tag,omitempty"`
	Version    string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Filename   string `protobuf:"bytes,4,opt,name=filename,proto3" json:"filename,omitempty"`
}

type BacktraceChunk struct {
	Data isBacktraceChunk_Data
}

type isBacktraceChunk_Data interface {
	isBacktraceChunk_Data()
}

type BacktraceChunk_Metadata struct {
	Metadata *BacktraceMetadata
}

type BacktraceChunk_Chunk struct {
	Chunk []byte
}

func (*BacktraceChunk_Metadata) isBacktraceChunk_Data() {}
func (*BacktraceChunk_Chunk) isBacktraceChunk_Data()    {}

type BacktraceMetadata struct {
	Uid        string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	ProcessTag string `protobuf:"bytes,2,opt,name=process_tag,json=processTag,proto3" json:"process_tag,omitempty"`
	Version    string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Filename   string `protobuf:"bytes,4,opt,name=filename,proto3" json:"filename,omitempty"`
}

type LogUploadResponse struct {
	Status  string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	LogId   string `protobuf:"bytes,2,opt,name=log_id,json=logId,proto3" json:"log_id,omitempty"`
	Uid     string `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	S3Key   string `protobuf:"bytes,4,opt,name=s3_key,json=s3Key,proto3" json:"s3_key,omitempty"`
	Message string `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
}

type GetLogMetadataRequest struct {
	LogId string `protobuf:"bytes,1,opt,name=log_id,json=logId,proto3" json:"log_id,omitempty"`
}

type LogMetadataResponse struct {
	LogId      string                 `protobuf:"bytes,1,opt,name=log_id,json=logId,proto3" json:"log_id,omitempty"`
	DeviceUid  string                 `protobuf:"bytes,2,opt,name=device_uid,json=deviceUid,proto3" json:"device_uid,omitempty"`
	LogType    string                 `protobuf:"bytes,3,opt,name=log_type,json=logType,proto3" json:"log_type,omitempty"`
	ProcessTag string                 `protobuf:"bytes,4,opt,name=process_tag,json=processTag,proto3" json:"process_tag,omitempty"`
	Version    string                 `protobuf:"bytes,5,opt,name=version,proto3" json:"version,omitempty"`
	Filename   string                 `protobuf:"bytes,6,opt,name=filename,proto3" json:"filename,omitempty"`
	S3Key      string                 `protobuf:"bytes,7,opt,name=s3_key,json=s3Key,proto3" json:"s3_key,omitempty"`
	Timestamp  *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

type ListLogsRequest struct {
	Uid     string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	LogType string `protobuf:"bytes,2,opt,name=log_type,json=logType,proto3" json:"log_type,omitempty"`
	Limit   int32  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset  int32  `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
}

type ListLogsResponse struct {
	Logs  []*LogMetadataResponse `protobuf:"bytes,1,rep,name=logs,proto3" json:"logs,omitempty"`
	Total int32                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}
