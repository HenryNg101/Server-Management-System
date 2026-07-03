package kafka

import (
	"context"

	kgo "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kgo.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kgo.NewReader(kgo.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *Consumer) Read(ctx context.Context) (kgo.Message, error) {
	return c.reader.ReadMessage(ctx)
}
