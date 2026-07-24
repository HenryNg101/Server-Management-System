package agent

import (
	"context"

	"github.com/HenryNg101/agents-service/internal/model"
)

type MockAgentService struct {
	existFn       func(ctx context.Context, key string, server *model.Agent) error
	pushMetricsFn func(ctx context.Context, msg MetricMessage) error
}

func (m *MockAgentService) AgentExist(ctx context.Context, key string, server *model.Agent) error {
	return m.existFn(ctx, key, server)
}

func (m *MockAgentService) PushMetrics(ctx context.Context, msg MetricMessage) error {
	return m.pushMetricsFn(ctx, msg)
}
