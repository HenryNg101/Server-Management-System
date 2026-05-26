package server

import (
	"context"
	"errors"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type Service interface {
	GetServers() ([]model.Server, error)
	CreateServer(req CreateServerRequest) (model.Server, error)
	GetServer(ctx context.Context, id uint) (model.Server, error)
	UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error)
	DeleteServer(ctx context.Context, id uint) error
}

type serverService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &serverService{repo: r}
}

func (s *serverService) GetServers() ([]model.Server, error) {
	return s.repo.FindAll()
}

func (s *serverService) CreateServer(req CreateServerRequest) (model.Server, error) {
	server := model.Server{
		Name:        req.Name,
		Status:      req.Status,
		IPv4Address: req.IPv4Address,
		Port:        req.Port,
		Protocol:    req.Protocol,
	}
	return s.repo.Create(server)
}

func (s *serverService) GetServer(ctx context.Context, id uint) (model.Server, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *serverService) UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error) {
	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates ONLY if provided
	if req.Name != nil {
		server.Name = *req.Name
	}

	if req.Status != nil {
		server.Status = *req.Status
	}

	if req.IPv4Address != nil {
		server.IPv4Address = *req.IPv4Address
	}

	if req.Port != nil {
		server.Port = *req.Port
	}

	if req.Protocol != nil {
		server.Protocol = *req.Protocol
	}

	updated, err := s.repo.Update(ctx, &server)
	if err != nil {
		return nil, err
	}

	return updated, nil
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
