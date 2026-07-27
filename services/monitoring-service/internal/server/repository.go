package server

import (
	"context"
	"time"

	"github.com/HenryNg101/monitoring-service/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetStats(ctx context.Context) (total int64, err error)
	GetCreatedAtMap(ctx context.Context) (map[uint]time.Time, error)
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

func (r *serverRepository) GetCreatedAtMap(ctx context.Context) (map[uint]time.Time, error) {
	type Row struct {
		ID        uint      `gorm:"column:id"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	var rows []Row
	err := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Select("id, created_at").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]time.Time, len(rows))
	for _, row := range rows {
		result[row.ID] = row.CreatedAt.UTC()
	}

	return result, nil
}
