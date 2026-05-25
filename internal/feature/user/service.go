package user

import (
	"github.com/HenryNg101/server-management-system/internal/model"
)

type Service interface {
	GetUsers() ([]model.User, error)
	CreateUser(user model.User) (model.User, error)
}

type userService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &userService{repo: r}
}

func (s *userService) GetUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) CreateUser(user model.User) (model.User, error) {
	return s.repo.Create(user)
}
