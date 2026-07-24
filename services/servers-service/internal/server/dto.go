package server

import (
	"time"

	"github.com/HenryNg101/servers-service/internal/model"
)

type GetServersQuery struct {
	Name *string

	Page     *int
	PageSize *int

	SortBy string
	Order  string
	UserID *uint
}

type GetServerResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	IPv4      string    `json:"ipv4"`
	Status    string    `json:"status"` // derived: UP / NO_DATA
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedServers struct {
	Servers    []model.Server `json:"servers"`
	Total      int64          `json:"total"`
	Page       *int           `json:"page"`
	PageSize   *int           `json:"page_size"`
	TotalPages *int           `json:"total_pages"`
}

type UpdateServerRequest struct {
	Name *string `json:"name"`
	IPv4 *string `json:"ipv4"`
}

type CreateServerRequest struct {
	Name string `json:"name" binding:"required"`
	IPv4 string `json:"ipv4" binding:"required,ip"`
}

type CreateServerResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	IPv4      string    `json:"ipv4"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
