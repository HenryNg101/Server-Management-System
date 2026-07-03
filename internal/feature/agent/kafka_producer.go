package agent

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/HenryNg101/server-management-system/internal/platform/kafka"
	kgo "github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	PublishMetrics(ctx context.Context, messages []MetricMessage) error
}

type kafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(p *kafka.Producer) KafkaProducer {
	return &kafkaProducer{producer: p}
}

func (k *kafkaProducer) PublishMetrics(ctx context.Context, messages []MetricMessage) error {
	kmsgs := make([]kgo.Message, 0, len(messages))

	for _, m := range messages {
		val, err := json.Marshal(m)
		if err != nil {
			return err
		}

		kmsgs = append(kmsgs, kgo.Message{
			Key:   []byte(strconv.Itoa(m.ServerID)),
			Value: val,
			Time:  m.Timestamp,
		})
	}

	return k.producer.WriteMany(ctx, kmsgs)
}
