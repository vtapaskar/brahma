package metrics

import "time"

type CPUStats struct {
	UID            string    `json:"uid"`
	Timestamp      time.Time `json:"timestamp"`
	UsagePercent   float64   `json:"usage_percent"`
	UserPercent    float64   `json:"user_percent"`
	SystemPercent  float64   `json:"system_percent"`
	IdlePercent    float64   `json:"idle_percent"`
	IOWaitPercent  float64   `json:"iowait_percent"`
	LoadAvg1Min    float64   `json:"load_avg_1min"`
	LoadAvg5Min    float64   `json:"load_avg_5min"`
	LoadAvg15Min   float64   `json:"load_avg_15min"`
	NumCores       int       `json:"num_cores"`
	PerCoreUsage   []float64 `json:"per_core_usage,omitempty"`
}

type ProcessStats struct {
	UID          string        `json:"uid"`
	Timestamp    time.Time     `json:"timestamp"`
	TotalCount   int           `json:"total_count"`
	RunningCount int           `json:"running_count"`
	SleepingCount int          `json:"sleeping_count"`
	ZombieCount  int           `json:"zombie_count"`
	Processes    []ProcessInfo `json:"processes,omitempty"`
}

type ProcessInfo struct {
	PID         int     `json:"pid"`
	Name        string  `json:"name"`
	State       string  `json:"state"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemPercent  float64 `json:"mem_percent"`
	MemoryRSS   uint64  `json:"memory_rss"`
	Threads     int     `json:"threads"`
	StartTime   int64   `json:"start_time"`
	Command     string  `json:"command,omitempty"`
}

type MgmtNetworkStats struct {
	UID             string              `json:"uid"`
	Timestamp       time.Time           `json:"timestamp"`
	InterfaceName   string              `json:"interface_name"`
	Status          string              `json:"status"`
	IPAddress       string              `json:"ip_address"`
	Netmask         string              `json:"netmask"`
	Gateway         string              `json:"gateway"`
	MACAddress      string              `json:"mac_address"`
	Speed           int64               `json:"speed"`
	Duplex          string              `json:"duplex"`
	RxBytes         uint64              `json:"rx_bytes"`
	TxBytes         uint64              `json:"tx_bytes"`
	RxPackets       uint64              `json:"rx_packets"`
	TxPackets       uint64              `json:"tx_packets"`
	RxErrors        uint64              `json:"rx_errors"`
	TxErrors        uint64              `json:"tx_errors"`
	RxDropped       uint64              `json:"rx_dropped"`
	TxDropped       uint64              `json:"tx_dropped"`
	DNSServers      []string            `json:"dns_servers,omitempty"`
	NTPServers      []string            `json:"ntp_servers,omitempty"`
}

type RouterBaseState struct {
	UID              string           `json:"uid"`
	Timestamp        time.Time        `json:"timestamp"`
	Hostname         string           `json:"hostname"`
	Platform         string           `json:"platform"`
	HardwareVersion  string           `json:"hardware_version"`
	SoftwareVersion  string           `json:"software_version"`
	SONiCVersion     string           `json:"sonic_version"`
	KernelVersion    string           `json:"kernel_version"`
	UptimeSeconds    int64            `json:"uptime_seconds"`
	BootTime         time.Time        `json:"boot_time"`
	SerialNumber     string           `json:"serial_number"`
	MgmtNetworkStatus MgmtStatus      `json:"mgmt_network_status"`
	DNSStatus        DNSStatus        `json:"dns_status"`
	DHCPStatus       DHCPStatus       `json:"dhcp_status"`
	LLDPStatus       LLDPStatus       `json:"lldp_status"`
	MemoryTotal      uint64           `json:"memory_total"`
	MemoryUsed       uint64           `json:"memory_used"`
	MemoryFree       uint64           `json:"memory_free"`
	DiskTotal        uint64           `json:"disk_total"`
	DiskUsed         uint64           `json:"disk_used"`
	DiskFree         uint64           `json:"disk_free"`
}

type MgmtStatus struct {
	Status      string `json:"status"`
	IPAddress   string `json:"ip_address"`
	Gateway     string `json:"gateway"`
	Reachable   bool   `json:"reachable"`
}

type DNSStatus struct {
	Enabled     bool     `json:"enabled"`
	Servers     []string `json:"servers"`
	Domain      string   `json:"domain,omitempty"`
	SearchList  []string `json:"search_list,omitempty"`
	Operational bool     `json:"operational"`
}

type DHCPStatus struct {
	Enabled     bool   `json:"enabled"`
	State       string `json:"state"`
	LeaseTime   int64  `json:"lease_time,omitempty"`
	RenewTime   int64  `json:"renew_time,omitempty"`
	ServerIP    string `json:"server_ip,omitempty"`
}

type LLDPStatus struct {
	Enabled       bool           `json:"enabled"`
	ChassisID     string         `json:"chassis_id"`
	SystemName    string         `json:"system_name"`
	NeighborCount int            `json:"neighbor_count"`
	Neighbors     []LLDPNeighbor `json:"neighbors,omitempty"`
}

type LLDPNeighbor struct {
	LocalPort      string `json:"local_port"`
	RemoteChassisID string `json:"remote_chassis_id"`
	RemotePortID   string `json:"remote_port_id"`
	RemoteSystemName string `json:"remote_system_name"`
	RemotePortDesc string `json:"remote_port_desc,omitempty"`
	TTL            int    `json:"ttl"`
}
