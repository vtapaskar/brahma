package models

import "time"

type SonicDevice struct {
	ID         string            `json:"id"`
	Hostname   string            `json:"hostname"`
	IPAddress  string            `json:"ip_address"`
	DeviceType string            `json:"device_type"`
	Platform   string            `json:"platform"`
	Version    string            `json:"version"`
	Labels     map[string]string `json:"labels,omitempty"`
	LastSeen   time.Time         `json:"last_seen"`
}

type InterfaceMetric struct {
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	Speed      int64   `json:"speed"`
	MTU        int     `json:"mtu"`
	RxBytes    uint64  `json:"rx_bytes"`
	TxBytes    uint64  `json:"tx_bytes"`
	RxPackets  uint64  `json:"rx_packets"`
	TxPackets  uint64  `json:"tx_packets"`
	RxErrors   uint64  `json:"rx_errors"`
	TxErrors   uint64  `json:"tx_errors"`
	RxDropped  uint64  `json:"rx_dropped"`
	TxDropped  uint64  `json:"tx_dropped"`
	Utilization float64 `json:"utilization"`
}

type SystemMetric struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryTotal    uint64  `json:"memory_total"`
	MemoryUsed     uint64  `json:"memory_used"`
	MemoryFree     uint64  `json:"memory_free"`
	DiskTotal      uint64  `json:"disk_total"`
	DiskUsed       uint64  `json:"disk_used"`
	Uptime         int64   `json:"uptime_seconds"`
	Temperature    float64 `json:"temperature"`
	FanSpeed       int     `json:"fan_speed"`
	PSUStatus      string  `json:"psu_status"`
}

type BGPMetric struct {
	NeighborIP     string `json:"neighbor_ip"`
	NeighborAS     int    `json:"neighbor_as"`
	State          string `json:"state"`
	PrefixReceived int    `json:"prefix_received"`
	PrefixSent     int    `json:"prefix_sent"`
	Uptime         int64  `json:"uptime_seconds"`
	MessagesIn     uint64 `json:"messages_in"`
	MessagesOut    uint64 `json:"messages_out"`
}

type VXLANMetric struct {
	VNI           int    `json:"vni"`
	SourceIP      string `json:"source_ip"`
	RemoteVTEPs   int    `json:"remote_vteps"`
	MACCount      int    `json:"mac_count"`
	ARPCount      int    `json:"arp_count"`
}

type QoSMetric struct {
	Queue          string `json:"queue"`
	Priority       int    `json:"priority"`
	PacketsQueued  uint64 `json:"packets_queued"`
	PacketsDropped uint64 `json:"packets_dropped"`
	BytesQueued    uint64 `json:"bytes_queued"`
}

type HealthStatus struct {
	DeviceID    string    `json:"device_id"`
	Status      string    `json:"status"`
	LastCheck   time.Time `json:"last_check"`
	Components  []ComponentHealth `json:"components"`
}

type ComponentHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
