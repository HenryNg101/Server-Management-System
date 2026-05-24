package service

import (
	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/HenryNg101/server-management-system/internal/repository"
)

type UserService interface {
	GetUsers() ([]model.User, error)
	CreateUser(user model.User) (model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) CreateUser(user model.User) (model.User, error) {
	return s.repo.Create(user)
}
