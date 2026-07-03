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

	buffer   []agent.MetricMessage // Contains documents data for ES bulk insert
	messages []kgo.Message         // Track Kafka messages for offset commitment + DLQ
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

	// We use two checks for flushes (i.e retry, commit and reset buffer): Time and buffer size
	// So, if the write throughput is too high, the batch size of the buffer will trigger flushes, and if the write throughput is too low, the flush timeout will trigger flushes
	// Either way, we ensure that the buffer is flushed in a timely manner, and we don't block on waiting for the timer to fire if we have enough messages to flush, or waiting too long without flushing
	for {
		select {
		case <-ticker.C:
			c.flush(ctx)

		default:
			c.consumeOne(ctx)
		}
	}
}

// Fetches a single message and appends to buffer
func (c *MetricsConsumer) consumeOne(ctx context.Context) {
	msg, err := c.consumer.Fetch(ctx)
	if err != nil {
		log.Printf("fetch error: %v", err)
		return
	}

	var metric agent.MetricMessage
	if err := json.Unmarshal(msg.Value, &metric); err != nil {
		log.Printf("invalid message: %v", err)

		// Commit invalid message to avoid blocking the partition by skipping the poison messages
		_ = c.consumer.Commit(ctx, msg)
		return
	}

	c.buffer = append(c.buffer, metric)
	c.messages = append(c.messages, msg)

	if len(c.buffer) >= batchSize {
		c.flush(ctx)
	}
}

// Flush handles:
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

		c.buffer = c.buffer[:0]
		c.messages = c.messages[:0]
		return
	}

	// Commit only after successful insert
	// This is committed after documents are inserted to ES, however, this can fail -> Potential duplicate insert of document to Elasticsearch
	if err := c.consumer.Commit(ctx, c.messages...); err != nil {
		log.Printf("commit failed: %v", err)
		return
	}

	log.Printf("flushed %d messages", len(c.buffer))
	c.buffer = c.buffer[:0]
	c.messages = c.messages[:0]
}
