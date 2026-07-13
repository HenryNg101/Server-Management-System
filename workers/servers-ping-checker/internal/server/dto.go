package server

import "github.com/HenryNg101/servers-ping-checker/internal/model"

type GetServersQuery struct {
	Status   *bool
	Protocol *string
	Name     *string

	Page     *int
	PageSize *int

	SortBy string
	Order  string
}

type PaginatedServers struct {
	Servers    []model.Server `json:"servers"`
	Total      int64          `json:"total"`
	Page       *int           `json:"page"`
	PageSize   *int           `json:"page_size"`
	TotalPages *int           `json:"total_pages"`
}
