package jobs

import (
	"context"

	"github.com/HenryNg101/jobs-service/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	// Jobs system
	CreateJob(ctx context.Context, job *model.ImportJob) error
	GetJob(ctx context.Context, id string) (*model.ImportJob, error)
}

type dataTransferRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &dataTransferRepository{db: db}
}

func (r *dataTransferRepository) CreateJob(ctx context.Context, job *model.ImportJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *dataTransferRepository) GetJob(ctx context.Context, id string) (*model.ImportJob, error) {
	var job model.ImportJob
	if err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
