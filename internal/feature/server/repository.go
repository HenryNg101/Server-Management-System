package server

import (
	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]model.Server, error)
	Create(user model.Server) (model.Server, error)
}

type serverRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &serverRepository{db: db}
}

func (r *serverRepository) FindAll() ([]model.Server, error) {
	var result []model.Server

	err := r.db.Model(&model.Server{}).Find(&result).Error
	return result, err
}

func (r *serverRepository) Create(user model.Server) (model.Server, error) {
	err := r.db.Create(&user).Error
	return user, err
}
