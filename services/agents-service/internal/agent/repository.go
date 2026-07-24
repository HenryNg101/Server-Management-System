package agent

import (
	"context"
	"errors"

	"github.com/HenryNg101/agents-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Upsert(ctx context.Context, agent *model.Agent) (*model.Agent, error)
	FindByKey(ctx context.Context, key string) (*model.Agent, error)
	FindByInstance(ctx context.Context, serverID uint, instanceID string) (*model.Agent, error)
	UpdateAPIKey(ctx context.Context, agentID uint, newHash string) error
}

type agentRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &agentRepository{db: db}
}

func (r *agentRepository) Upsert(ctx context.Context, agent *model.Agent) (*model.Agent, error) {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "server_id"}, // conflict target
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"api_key":      agent.APIKey,
				"instance_id":  agent.InstanceID,
				"status":       agent.Status,
				"last_seen_at": gorm.Expr("NOW()"),
			}),
		}).
		Create(agent).Error

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

func (r *agentRepository) FindByInstance(ctx context.Context, serverID uint, instanceID string) (*model.Agent, error) {
	var agent model.Agent
	err := r.db.WithContext(ctx).
		Where("server_id = ? AND instance_id = ?", serverID, instanceID).
		First(&agent).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *agentRepository) UpdateAPIKey(ctx context.Context, agentID uint, newHash string) error {
	return r.db.WithContext(ctx).
		Model(&model.Agent{}).
		Where("id = ?", agentID).
		Update("api_key", newHash).Error
}
