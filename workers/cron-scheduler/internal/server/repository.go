package server

import (
	"context"

	"github.com/HenryNg101/cron-scheduler/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetStats(ctx context.Context) (total, up, down int64, err error)
}

type serverRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &serverRepository{db: db}
}

func (r *serverRepository) GetStats(ctx context.Context) (total, up, down int64, err error) {
	err = r.db.WithContext(ctx).Model(&model.Server{}).Count(&total).Error
	if err != nil {
		return 0, 0, 0, err
	}

	err = r.db.WithContext(ctx).Model(&model.Server{}).
		Where("status = ?", true).
		Count(&up).Error
	if err != nil {
		return 0, 0, 0, err
	}

	down = total - up
	return
}
