package domain

import (
	"time"
)

type StatResources struct {
	Resources []Stats
}

type Stats struct {
	Type               string  `json:"type"`
	Index              int     `json:"index"`
	InstanceGuid       string  `json:"instance_guid"`
	State              string  `json:"state"`
	Routable           *bool   `json:"routable"`
	Host               string  `json:"host"`
	InstanceInternalIp string  `json:"instance_internal_ip"`
	Uptime             int     `json:"uptime"`
	MemQuota           *int    `json:"mem_quota"`
	DiskQuota          *int    `json:"disk_quota"`
	LogRateLimit       *int    `json:"log_rate_limit"`
	FdsQuota           int     `json:"fds_quota"`
	IsolationSegment   *string `json:"isolation_segment"`
	Details            *string `json:"details"`
	InstancePorts      []struct {
		External             int  `json:"external"`
		Internal             int  `json:"internal"`
		ExternalTlsProxyPort *int `json:"external_tls_proxy_port"`
		InternalTlsProxyPort int  `json:"internal_tls_proxy_port"`
	} `json:"instance_ports"`
	Usage struct {
		Time           time.Time `json:"time"`
		Cpu            float64   `json:"cpu"`
		CpuEntitlement float64   `json:"cpu_entitlement"`
		Mem            int       `json:"mem"`
		Disk           int       `json:"disk"`
		LogRate        int       `json:"log_rate"`
	} `json:"usage"`
}
