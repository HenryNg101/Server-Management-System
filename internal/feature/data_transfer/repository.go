package data_transfer

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, server *model.Server) (*model.Server, error)
	BulkUpsert(ctx context.Context, servers []*model.Server) error
}

type dataTransferRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &dataTransferRepository{db: db}
}

// Create if not exist, otherwise, returns error
func (r *dataTransferRepository) Create(ctx context.Context, server *model.Server) (*model.Server, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(server)

	if result.Error != nil {
		return nil, result.Error
	}
	return server, nil
}

// Bulk upsert servers, used for importing servers from CSV file
func (r *dataTransferRepository) BulkUpsert(ctx context.Context, servers []*model.Server) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(&servers).Error
}
