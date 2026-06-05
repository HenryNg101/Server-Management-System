package server

import (
	"time"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type GetServersQuery struct {
	Status   *bool
	Protocol *string
	Name     *string

	Page     *int
	PageSize *int

	SortBy string
	Order  string
}

type GetServerResponse struct {
	Name        string    `json:"name"`
	Status      bool      `json:"status"`       // Server status at the time (Is it on or off)
	IPv4Address string    `json:"ipv4_address"` // IPv4 of the server
	Port        uint      `json:"port"`         // Port of the server
	Protocol    string    `json:"protocol"`     // Network protocol that the server use, could be TCP, FTP, SSH, etc.
	CreatedAt   time.Time `json:"created_at"`   // Automatically managed by GORM for creation time
	LastUpdated time.Time `json:"last_updated"` // Last time the server got updated
}

type PaginatedServers struct {
	Servers    []model.Server `json:"servers"`
	Total      int64          `json:"total"`
	Page       *int           `json:"page"`
	PageSize   *int           `json:"page_size"`
	TotalPages *int           `json:"total_pages"`
}

type UpdateServerRequest struct {
	Name        *string `json:"name"`
	Status      *bool   `json:"status"`
	IPv4Address *string `json:"ipv4"`
	Port        *uint   `json:"port"`
	Protocol    *string `json:"protocol"`
}

type CreateServerRequest struct {
	Name        string `json:"name" binding:"required"`
	Status      bool   `json:"status" binding:"required"`
	IPv4Address string `json:"ipv4_address" binding:"required,ip"`
	Port        uint   `json:"port" binding:"required"`
	Protocol    string `json:"protocol"`
}

type ImportServersResponse struct {
	SuccessCount int             `json:"success_count"`
	FailedCount  int             `json:"failed_count"`
	Successes    []model.Server  `json:"successes"`
	Failures     []ImportFailure `json:"failures"`
}

type ImportFailure struct {
	Row    int               `json:"row"`
	Error  string            `json:"error"`
	Record map[string]string `json:"record"`
}

type Report struct {
	TotalServers int64            `json:"total_servers"`
	ServersUp    int64            `json:"servers_up"`
	ServersDown  int64            `json:"servers_down"`
	Uptime       map[uint]float64 `json:"uptime_per_server"`
}

type SendReportRequest struct {
	Start  *string   `json:"start"`
	End    *string   `json:"end"`
	TopN   *int      `json:"count"`            // The number of top N worst servers you want to show in the report
	Emails *[]string `json:"emails,omitempty"` // Optional. It might sends emails or nah
}
