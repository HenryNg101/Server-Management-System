package agent

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/HenryNg101/server-management-system/internal/shared/auth"
)

type Service interface {
	AgentExist(ctx context.Context, key string, server *model.Agent) error
	PushMetrics(ctx context.Context, msg MetricMessage) error
}

type agentService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &agentService{repo: r}
}

// TODO: Push to Kafka
func (s *agentService) PushMetrics(ctx context.Context, msg MetricMessage) error {
	return nil
}

func (s *agentService) AgentExist(ctx context.Context, key string, server *model.Agent) error {
	hashedKey := auth.HashAPIKey(key)
	_, err := s.repo.FindByKey(ctx, hashedKey)

	if err != nil {
		return err
	}
	return nil
}
