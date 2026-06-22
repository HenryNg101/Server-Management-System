package agent

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/HenryNg101/server-management-system/internal/shared/auth"
)

type Service interface {
	AgentExist(ctx context.Context, key string, server *model.Agent) error
	PushMetrics(ctx context.Context, messages []MetricMessage) error
}

type agentService struct {
	repo        Repository
	elasticRepo ElasticAgentRepository
}

func NewService(r Repository, e ElasticAgentRepository) Service {
	return &agentService{repo: r, elasticRepo: e}
}

// TODO: Push to Kafka instead of directly push into Elasticsearch like this
func (s *agentService) PushMetrics(ctx context.Context, messages []MetricMessage) error {
	return s.elasticRepo.BulkInsertStatus(ctx, messages)
}

func (s *agentService) AgentExist(ctx context.Context, key string, server *model.Agent) error {
	hashedKey := auth.HashAPIKey(key)
	_, err := s.repo.FindByKey(ctx, hashedKey)

	if err != nil {
		return err
	}
	return nil
}
