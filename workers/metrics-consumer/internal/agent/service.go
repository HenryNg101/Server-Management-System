package agent

import (
	"context"
)

type Service interface {
	PushMetricsToElastic(ctx context.Context, messages []MetricMessage) error
}

type agentService struct {
	elasticRepo ElasticAgentRepository
}

func NewService(e ElasticAgentRepository) Service {
	return &agentService{elasticRepo: e}
}

func (s *agentService) PushMetricsToElastic(ctx context.Context, messages []MetricMessage) error {
	return s.elasticRepo.BulkInsertStatus(ctx, messages)
}
