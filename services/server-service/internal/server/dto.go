package server

import (
	"time"

	"github.com/HenryNg101/server-service/internal/model"
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

type CreateServerResponse struct {
	CreatedServer model.Server `json:"created_server"`
	CreatedAgent  model.Agent  `json:"created_agent"`
	ApiKey        string       `json:"agent_api_key"`
}

type ServerPullStats struct {
	Uptime float64 `json:"uptime"`

	// TODO: Add this later
	// LatencyP95 float64
}
