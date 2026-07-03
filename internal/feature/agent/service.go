package agent

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/HenryNg101/server-management-system/internal/shared/auth"
)

type Service interface {
	AgentExist(ctx context.Context, key string, server *model.Agent) error
	PushMetrics(ctx context.Context, messages []MetricMessage) error
	PushMetricsToElastic(ctx context.Context, messages []MetricMessage) error
}

type agentService struct {
	repo          Repository
	kafkaProducer KafkaProducer
	elasticRepo   ElasticAgentRepository
}

func NewService(r Repository, k KafkaProducer, e ElasticAgentRepository) Service {
	return &agentService{repo: r, kafkaProducer: k, elasticRepo: e}
}

func (s *agentService) PushMetrics(ctx context.Context, messages []MetricMessage) error {
	return s.kafkaProducer.PublishMetrics(ctx, messages)
}

func (s *agentService) PushMetricsToElastic(ctx context.Context, messages []MetricMessage) error {
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
