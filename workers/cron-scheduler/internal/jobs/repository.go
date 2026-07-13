package data_transfer

import (
	"context"
	"time"

	"github.com/HenryNg101/cron-scheduler/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetOldJobs(ctx context.Context, olderThan time.Duration) ([]model.ImportJob, error)
	DeleteJob(ctx context.Context, jobID string) error
}

type dataTransferRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &dataTransferRepository{db: db}
}

func (r *dataTransferRepository) GetOldJobs(ctx context.Context, olderThan time.Duration) ([]model.ImportJob, error) {
	var jobs []model.ImportJob

	cutoff := time.Now().Add(-olderThan)

	err := r.db.WithContext(ctx).
		Where("created_at < ?", cutoff).
		Find(&jobs).Error

	return jobs, err
}

func (r *dataTransferRepository) DeleteJob(ctx context.Context, jobID string) error {
	return r.db.WithContext(ctx).
		Delete(&model.ImportJob{}, "id = ?", jobID).Error
}
