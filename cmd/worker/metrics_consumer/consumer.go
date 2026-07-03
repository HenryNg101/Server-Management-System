package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/HenryNg101/server-management-system/internal/feature/agent"
	kafkaClient "github.com/HenryNg101/server-management-system/internal/platform/kafka"
	kgo "github.com/segmentio/kafka-go"
)

// TODO: Refactor code for this further. Maybe split into multiple files for better organization and readability.
type MetricsConsumer struct {
	consumer    *kafkaClient.Consumer
	dlqProducer *kafkaClient.Producer
	service     agent.Service

	buffer   []agent.MetricMessage
	messages []kgo.Message
}

func NewMetricsConsumer(consumer *kafkaClient.Consumer, dlq *kafkaClient.Producer, service agent.Service) *MetricsConsumer {
	return &MetricsConsumer{
		consumer:    consumer,
		dlqProducer: dlq,
		service:     service,
		buffer:      make([]agent.MetricMessage, 0, batchSize),
		messages:    make([]kgo.Message, 0, batchSize),
	}
}

// Run starts the main consumption loop
func (c *MetricsConsumer) Run(ctx context.Context) {
	ticker := time.NewTicker(flushTimeout)
	defer ticker.Stop()

	log.Println("metrics consumer started")

	for {
		select {
		case <-ticker.C:
			c.flush(ctx)

		default:
			c.consumeOne(ctx)
		}
	}
}

// consumeOne fetches a single message and appends to buffer
func (c *MetricsConsumer) consumeOne(ctx context.Context) {
	msg, err := c.consumer.Fetch(ctx)
	if err != nil {
		log.Printf("fetch error: %v", err)
		return
	}

	var metric agent.MetricMessage
	if err := json.Unmarshal(msg.Value, &metric); err != nil {
		log.Printf("invalid message: %v", err)

		// Commit invalid message to avoid blocking the partition
		_ = c.consumer.Commit(ctx, msg)
		return
	}

	c.buffer = append(c.buffer, metric)
	c.messages = append(c.messages, msg)

	if len(c.buffer) >= batchSize {
		c.flush(ctx)
	}
}

// flush handles:
// 1. retrying ES insert
// 2. sending to DLQ if failed
// 3. committing offsets only on final outcome
func (c *MetricsConsumer) flush(ctx context.Context) {
	if len(c.buffer) == 0 {
		return
	}

	err := c.retryInsert(ctx)

	if err != nil {
		log.Printf("bulk insert failed after retries, sending to DLQ")

		c.sendToDLQ(ctx, err.Error())

		// Commit even on failure so we don't retry forever
		if commitErr := c.consumer.Commit(ctx, c.messages...); commitErr != nil {
			log.Printf("commit after DLQ failed: %v", commitErr)
		}

		c.reset()
		return
	}

	// Commit only after successful insert
	if err := c.consumer.Commit(ctx, c.messages...); err != nil {
		log.Printf("commit failed: %v", err)
		return
	}

	log.Printf("flushed %d messages", len(c.buffer))
	c.reset()
}

// retryInsert retries ES bulk insert
func (c *MetricsConsumer) retryInsert(ctx context.Context) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = c.service.PushMetricsToElastic(ctx, c.buffer)
		if err == nil {
			return nil
		}

		log.Printf("retry %d failed: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return err
}

// sendToDLQ pushes failed messages to DLQ topic
func (c *MetricsConsumer) sendToDLQ(ctx context.Context, reason string) {
	for _, m := range c.messages {
		dlqMsg := kgo.Message{
			Key:   m.Key,
			Value: m.Value,
			Time:  time.Now(),
			Headers: []kgo.Header{
				{Key: "error", Value: []byte(reason)},
			},
		}

		if err := c.dlqProducer.WriteOne(ctx, dlqMsg); err != nil {
			log.Printf("failed to write DLQ: %v", err)
		}
	}
}

// reset clears buffer after flush
func (c *MetricsConsumer) reset() {
	c.buffer = c.buffer[:0]
	c.messages = c.messages[:0]
}
