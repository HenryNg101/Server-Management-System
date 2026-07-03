package main

import (
	"context"
	"log"
	"time"

	kgo "github.com/segmentio/kafka-go"
)

// Pushes failed messages to DLQ topic
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

		// Best effort write. If a DLQ write fails, we log it but continue to the next message
		// This prevents one bad message from blocking the rest, with the trade-off that some messages may be lost if the DLQ is down.
		if err := c.dlqProducer.WriteOne(ctx, dlqMsg); err != nil {
			log.Printf("failed to write DLQ: %v", err)
		}
	}
}
