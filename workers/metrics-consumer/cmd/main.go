package main

import (
	"context"
	"log"
	"time"

	"github.com/HenryNg101/metrics-consumer/internal/app"
	"github.com/HenryNg101/metrics-consumer/internal/config"
	kafkaClient "github.com/HenryNg101/metrics-consumer/internal/platform/kafka"
)

const (
	batchSize    = 500
	flushTimeout = 1 * time.Second
	maxRetries   = 3
)

func main() {
	// Load app and config
	ctx := context.Background()
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("init app failed: %v", err)
	}
	kafkaCfg := config.LoadKafka()

	// Create Kafka consumer and producer for dead-letter queue
	consumer := kafkaClient.NewConsumer(
		kafkaCfg.Brokers,
		kafkaCfg.AgentMetricsTopic,
		kafkaCfg.AgentMetricsConsumerGroup,
	)
	dlqProducer := kafkaClient.NewProducer(
		kafkaCfg.Brokers,
		kafkaCfg.AgentMetricsDLQTopic,
	)

	metricsConsumer := NewMetricsConsumer(
		consumer,
		dlqProducer,
		application.AgentService,
	)
	metricsConsumer.Run(ctx)
}
