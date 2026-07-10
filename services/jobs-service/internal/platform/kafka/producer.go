package kafka

import (
	"context"

	kgo "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kgo.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kgo.Writer{
			Addr:     kgo.TCP(brokers...),
			Topic:    topic,
			Balancer: &kgo.Hash{},
		},
	}
}

func (p *Producer) WriteOne(ctx context.Context, msg kgo.Message) error {
	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) WriteMany(ctx context.Context, msgs []kgo.Message) error {
	return p.writer.WriteMessages(ctx, msgs...)
}
