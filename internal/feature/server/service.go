package server

import (
	"github.com/HenryNg101/server-management-system/internal/model"
)

type Service interface {
	GetServers() ([]model.Server, error)
	CreateServer(user model.Server) (model.Server, error)
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

func (s *serverService) CreateServer(user model.Server) (model.Server, error) {
	return s.repo.Create(user)
}
