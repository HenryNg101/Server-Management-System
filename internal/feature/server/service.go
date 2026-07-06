package server

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/HenryNg101/server-management-system/internal/shared/auth"
)

type Service interface {
	GetServers(ctx context.Context, q GetServersQuery) (*PaginatedServers, error)
	CreateServer(ctx context.Context, req CreateServerRequest) (*model.Server, error)
	GetServer(ctx context.Context, id uint, server *model.Server) (*model.Server, error)
	UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error)
	DeleteServer(ctx context.Context, id uint) error
	BulkUpdateServersStatuses(ctx context.Context, servers []*model.Server) error

	// Elastic services
	ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error

	// Monitoring agent
	CreateAgent(ctx context.Context, serverID uint) (*model.Agent, string, error)
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

func (s *serverService) CreateServer(ctx context.Context, req CreateServerRequest) (*model.Server, error) {
	server := model.Server{
		Name:        req.Name,
		Status:      req.Status,
		IPv4Address: req.IPv4Address,
		Port:        req.Port,
		Protocol:    req.Protocol,
	}
	return s.repo.Create(ctx, &server)
}

// TODO: Use this in server's creation
func (s *serverService) CreateAgent(ctx context.Context, serverID uint) (*model.Agent, string, error) {
	rawKey, hash, err := auth.GenerateAPIKey()
	if err != nil {
		return nil, "", err
	}

	agent := model.Agent{
		ServerID: serverID,
		APIKey:   hash,
	}
	createdAgent, err := s.repo.CreateAgent(ctx, &agent)

	return createdAgent, rawKey, err
}

func (s *serverService) GetServer(ctx context.Context, id uint, server *model.Server) (*model.Server, error) {
	return s.repo.FindByID(ctx, id, server)
}

func (s *serverService) UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error) {
	server := &model.Server{}

	server, err := s.repo.FindByID(ctx, id, server)
	if err != nil {
		return nil, err
	}

	// Apply updates ONLY if provided
	isUpdated := false
	if req.Name != nil {
		isUpdated = true
		server.Name = *req.Name
	}

	if req.Status != nil {
		isUpdated = true
		server.Status = *req.Status
	}

	if req.IPv4Address != nil {
		isUpdated = true
		server.IPv4Address = *req.IPv4Address
	}

	if req.Port != nil {
		isUpdated = true
		server.Port = *req.Port
	}

	if req.Protocol != nil {
		isUpdated = true
		server.Protocol = *req.Protocol
	}

	// If anything changes, it means that, it's actually updated
	if !isUpdated {
		return nil, errors.New("Nothing is updated")
	}

	server.LastUpdated = time.Now()
	updated, err := s.repo.Update(ctx, server)
	if err != nil {
		return nil, err
	}

	return updated, nil
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

var ErrNotFound = errors.New("server not found")

func (s *serverService) DeleteServer(ctx context.Context, id uint) error {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *serverService) ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error {
	return s.elasticRepo.BulkInsertStatus(ctx, serversResults)
}
