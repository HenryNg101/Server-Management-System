package server

import (
	"context"
	"errors"

	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	FindAll() ([]model.Server, error)
	Create(ctx context.Context, server *model.Server) (*model.Server, error)
	FindByID(ctx context.Context, id uint, server *model.Server) (*model.Server, error)
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
	var servers []model.Server

	err := r.db.Model(&model.Server{}).Find(&servers).Error
	return servers, err
}

// Create if not exist, otherwise, returns error
func (r *serverRepository) Create(ctx context.Context, server *model.Server) (*model.Server, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(server)

	if result.Error != nil {
		return nil, result.Error
	}
	// Handle a failure to insert due to record already exist, not validations or something
	if result.RowsAffected == 0 {
		return nil, errors.New("The server with this name is already existed")
	}
	return server, nil
}

func (r *serverRepository) FindByID(ctx context.Context, id uint, server *model.Server) (*model.Server, error) {
	err := r.db.Model(&model.Server{}).
		Where("id = ?", id).
		Find(&server).
		Error
	return server, err
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
