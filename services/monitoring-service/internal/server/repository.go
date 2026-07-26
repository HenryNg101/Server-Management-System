package server

import (
	"context"

	"github.com/HenryNg101/monitoring-service/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetStats(ctx context.Context) (total int64, err error)
}

type serverRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &serverRepository{db: db}
}

func (r *serverRepository) GetStats(ctx context.Context) (total int64, err error) {
	err = r.db.WithContext(ctx).Model(&model.Server{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return
}
