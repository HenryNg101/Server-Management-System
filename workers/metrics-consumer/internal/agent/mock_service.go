package agent

import (
	"context"
)

type MockAgentService struct {
	pushMetricsFn func(ctx context.Context, msg MetricMessage) error
}

func (m *MockAgentService) PushMetrics(ctx context.Context, msg MetricMessage) error {
	return m.pushMetricsFn(ctx, msg)
}
