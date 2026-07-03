package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/HenryNg101/server-management-system/internal/app"
	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/feature/agent"
	kafkaClient "github.com/HenryNg101/server-management-system/internal/platform/kafka"
)

const (
	batchSize    = 500
	flushTimeout = 1 * time.Second
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	kafkaCfg := config.LoadKafka()

	consumer := kafkaClient.NewConsumer(
		kafkaCfg.Brokers,
		kafkaCfg.AgentMetricsTopic,
		kafkaCfg.AgentMetricsConsumerGroup,
	)

	log.Println("Metrics consumer started...")

	buffer := make([]agent.MetricMessage, 0, batchSize)
	ticker := time.NewTicker(flushTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			flushBuffer(ctx, application, &buffer)

		default:
			msg, err := consumer.Read(ctx)
			if err != nil {
				log.Printf("read error: %v", err.Error())
				continue
			}

			var metric agent.MetricMessage
			if err := json.Unmarshal(msg.Value, &metric); err != nil {
				log.Printf("invalid message: %v", err)
				continue
			}

			buffer = append(buffer, metric)

			if len(buffer) >= batchSize {
				flushBuffer(ctx, application, &buffer)
			}
		}
	}
}
