package server

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]model.Server, error)
	Create(user model.Server) (model.Server, error)
	FindByID(ctx context.Context, id uint) (model.Server, error)
	Update(ctx context.Context, server *model.Server) (*model.Server, error)
	ExistsByID(ctx context.Context, id uint) (bool, error)
	Delete(ctx context.Context, id uint) error
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

func (r *serverRepository) FindByID(ctx context.Context, id uint) (model.Server, error) {
	var result model.Server

	err := r.db.Model(&model.Server{}).
		Where("id = ?", id).
		Find(&result).
		Error
	return result, err
}

func (r *serverRepository) Update(ctx context.Context, server *model.Server) (*model.Server, error) {
	err := r.db.WithContext(ctx).
		Save(server).Error

	return server, err
}

func (r *serverRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *serverRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Delete(&model.Server{}, id).Error
}
