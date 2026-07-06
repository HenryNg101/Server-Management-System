package data_transfer

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	// Servers
	CreateServer(ctx context.Context, server *model.Server) error
	BulkUpsertServers(ctx context.Context, servers []*model.Server) error

	// Jobs system
	CreateJob(ctx context.Context, job *model.ImportJob) error
	UpdateJobStatus(ctx context.Context, jobID string, status string, errMsg string) error
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

func (r *dataTransferRepository) CreateJob(ctx context.Context, job *model.ImportJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *dataTransferRepository) UpdateJobStatus(ctx context.Context, jobID string, status string, errMsg string) error {
	return r.db.WithContext(ctx).
		Model(&model.ImportJob{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status": status,
			"error":  errMsg,
		}).Error
}
