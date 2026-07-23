package server

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/HenryNg101/server-service/internal/model"
)

type Service interface {
	GetServers(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error)
	CreateServer(ctx context.Context, req CreateServerRequest, userID uint) (*model.Server, error)
	GetServer(ctx context.Context, serverID uint, userCtx *UserContext, server *model.Server) (*model.Server, error)
	UpdateServer(ctx context.Context, serverID uint, userCtx *UserContext, req UpdateServerRequest) (*model.Server, error)
	DeleteServer(ctx context.Context, serverID uint, userCtx *UserContext) error
}

type serverService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &serverService{repo: r}
}

func (s *serverService) GetServers(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error) {
	if !(userCtx.Role == "admin") {
		q.UserID = &userCtx.UserID
	}

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

func (s *serverService) CreateServer(ctx context.Context, req CreateServerRequest, userID uint) (*model.Server, error) {
	server := model.Server{
		Name:   req.Name,
		IPv4:   req.IPv4,
		UserID: userID,
	}

	return s.repo.Create(ctx, &server)
}

func (s *serverService) GetServer(ctx context.Context, serverID uint, userCtx *UserContext, server *model.Server) (*model.Server, error) {
	server, err := s.repo.FindByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// If not admin, user must have the right to see it
	if !(userCtx.Role == "admin") && server.UserID != userCtx.UserID {
		return nil, errors.New("forbidden")
	}
	return server, nil
}

func (s *serverService) UpdateServer(ctx context.Context, serverID uint, userCtx *UserContext, req UpdateServerRequest) (*model.Server, error) {
	server, err := s.repo.FindByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// Basic role check. If not admin, you can only update your own server
	if userCtx.Role != "admin" && server.UserID != userCtx.UserID {
		return nil, errors.New("forbidden: cannot update this server")
	}

	// Apply updates (only if provided)
	updated := false

	if req.Name != nil {
		server.Name = *req.Name
		updated = true
	}

	if req.IPv4 != nil {
		server.IPv4 = *req.IPv4
		updated = true
	}

	if !updated {
		return nil, errors.New("nothing to update")
	}

	server.UpdatedAt = time.Now()

	return s.repo.Update(ctx, server)
}

var ErrNotFound = errors.New("server not found")

func (s *serverService) DeleteServer(ctx context.Context, serverID uint, userCtx *UserContext) error {
	server, err := s.repo.FindByID(ctx, serverID)
	if err != nil {
		return err
	}

	// Basic role check again
	if userCtx.Role != "admin" && server.UserID != userCtx.UserID {
		return errors.New("forbidden: cannot delete this server")
	}

	return s.repo.Delete(ctx, serverID)
}
