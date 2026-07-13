package server

import (
	"context"
	"math"

	"github.com/HenryNg101/servers-ping-checker/internal/model"
)

type Service interface {
	GetServers(ctx context.Context, q GetServersQuery) (*PaginatedServers, error)
	BulkUpdateServersStatuses(ctx context.Context, servers []*model.Server) error

	// Elastic services
	ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error
}

type serverService struct {
	repo        Repository
	elasticRepo ElasticServerRepository
}

func NewService(r Repository, elastic ElasticServerRepository) Service {
	return &serverService{repo: r, elasticRepo: elastic}
}

func (s *serverService) GetServers(ctx context.Context, q GetServersQuery) (*PaginatedServers, error) {
	servers, total, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return nil, err
	}

	var totalPages int
	if q.PageSize != nil {
		totalPages = int(math.Ceil(float64(total) / float64(*q.PageSize)))
	}

	return &PaginatedServers{
		Servers:    servers,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: &totalPages,
	}, nil
}

// Bulk update statuses to be used by worker
func (s *serverService) BulkUpdateServersStatuses(ctx context.Context, servers []*model.Server) error {
	const chunkSize = 1000

	for i := 0; i < len(servers); i += chunkSize {
		end := i + chunkSize
		if end > len(servers) {
			end = len(servers)
		}

		err := s.repo.BulkUpdateStatus(ctx, servers[i:end])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serverService) ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error {
	return s.elasticRepo.BulkInsertStatus(ctx, serversResults)
}
