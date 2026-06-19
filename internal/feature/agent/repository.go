package agent

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, agent *model.Agent) (*model.Agent, error)
	FindByKey(ctx context.Context, key string) (*model.Agent, error)
}

type agentRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &agentRepository{db: db}
}

func (r *agentRepository) Create(ctx context.Context, agent *model.Agent) (*model.Agent, error) {
	err := r.db.WithContext(ctx).Create(&agent).Error
	return agent, err
}

func (r *agentRepository) FindByKey(ctx context.Context, key string) (*model.Agent, error) {
	var result model.Agent
	err := r.db.WithContext(ctx).
		Where("api_key = ?", key).
		First(&result).
		Error
	return &result, err
}
