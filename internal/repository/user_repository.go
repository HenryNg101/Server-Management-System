package repository

import "github.com/HenryNg101/server-management-system/internal/model"

type UserRepository interface {
	FindAll() ([]model.User, error)
	Create(user model.User) (model.User, error)
}
