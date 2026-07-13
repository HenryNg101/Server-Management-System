package jobs

import (
	"context"

	"github.com/HenryNg101/files-import-consumer/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	// Servers
	CreateServer(ctx context.Context, server *model.Server) error
	BulkUpsertServers(ctx context.Context, servers []*model.Server) error

	// Jobs system
	GetJob(ctx context.Context, id string) (*model.ImportJob, error)
	UpdateJobStatus(ctx context.Context, jobID string, status model.JobStatus, errMsg *string, resultPath *string) error
	UpdateJobProgress(ctx context.Context, id string, processed int, success int, failed int) error
}

type dataTransferRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &dataTransferRepository{db: db}
}

// Create if not exist, otherwise, returns error
func (r *dataTransferRepository) CreateServer(ctx context.Context, server *model.Server) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(server).Error
}

// Bulk upsert servers, used for importing servers from CSV file
func (r *dataTransferRepository) BulkUpsertServers(ctx context.Context, servers []*model.Server) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(&servers).Error
}

func (r *dataTransferRepository) GetJob(ctx context.Context, id string) (*model.ImportJob, error) {
	var job model.ImportJob
	if err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *dataTransferRepository) UpdateJobStatus(ctx context.Context, jobID string, status model.JobStatus, errMsg *string, resultPath *string) error {
	return r.db.WithContext(ctx).
		Model(&model.ImportJob{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":      status,
			"error":       errMsg,
			"result_path": resultPath,
		}).Error
}

func (r *dataTransferRepository) UpdateJobProgress(ctx context.Context, id string, processed int, success int, failed int) error {
	return r.db.WithContext(ctx).
		Model(&model.ImportJob{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed_rows":     processed,
			"success_rows_count": success,
			"failed_rows_count":  failed,
		}).Error
}
