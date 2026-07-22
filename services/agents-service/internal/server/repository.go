package server

import (
	"context"

	"github.com/HenryNg101/agent-service/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindByNameAndUser(ctx context.Context, serverName string) (*model.Server, error)
}

type serverRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &serverRepository{db: db}
}

func (r *serverRepository) FindByNameAndUser(ctx context.Context, serverName string) (*model.Server, error) {
	var server model.Server
	err := r.db.WithContext(ctx).
		Where("name = ?", serverName).
		First(&server).Error

	if err != nil {
		return nil, err
	}
	return &server, nil
}
