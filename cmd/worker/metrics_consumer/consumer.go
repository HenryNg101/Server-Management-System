package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/HenryNg101/server-management-system/internal/app"
	"github.com/HenryNg101/server-management-system/internal/feature/agent"
	kafkaClient "github.com/HenryNg101/server-management-system/internal/platform/kafka"
)

func readMessages(ctx context.Context, consumer *kafkaClient.Consumer) (*agent.MetricMessage, error) {
	var metric agent.MetricMessage
	msg, err := consumer.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}

	if err := json.Unmarshal(msg.Value, &metric); err != nil {
		return nil, fmt.Errorf("invalid message: %v", err)
	}
	return &metric, nil
}

func flushBuffer(ctx context.Context, application *app.Application, buffer *[]agent.MetricMessage) {
	if len(*buffer) == 0 {
		return
	}

	err := application.AgentService.PushMetricsToElastic(ctx, *buffer)
	if err != nil {
		log.Printf("bulk insert failed: %v", err)
		return
	}

	log.Printf("flushed %d metrics", len(*buffer))
	*buffer = (*buffer)[:0]
}
