package monitoring

import (
	"time"

	"github.com/HenryNg101/monitoring-service/internal/agent"
)

type SendReportRequest struct {
	Start  *string   `json:"start"`
	End    *string   `json:"end"`
	TopN   *int      `json:"count"`            // The number of top N worst servers you want to show in the report
	Emails *[]string `json:"emails,omitempty"` // Optional. It might sends emails or nah
}

type Report struct {
	TotalServers int64             `json:"total_servers"`
	ServersUp    int64             `json:"servers_up"`
	ServersDown  int64             `json:"servers_down"`
	Stats        []*ServerOverview `json:"servers_stats"`
}

type ServerOverview struct {
	ServerID uint

	// The time window actually used for uptime calculation for this server.
	// This is clamped to max(report_start, server_created_at).
	ActualStart time.Time `json:"actual_start"`
	ActualEnd   time.Time `json:"actual_end"`

	// Push model (agent based)
	agent.ServerPushStats
}
