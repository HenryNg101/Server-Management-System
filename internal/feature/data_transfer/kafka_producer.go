package data_transfer

import (
	"context"

	kafkaClient "github.com/HenryNg101/server-management-system/internal/platform/kafka"
	kgo "github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	WriteOne(ctx context.Context, key []byte, value []byte) error
}

type kafkaProducer struct {
	producer *kafkaClient.Producer
}

func NewKafkaProducer(p *kafkaClient.Producer) KafkaProducer {
	return &kafkaProducer{producer: p}
}

func (k *kafkaProducer) WriteOne(ctx context.Context, key []byte, value []byte) error {
	return k.producer.WriteOne(ctx, kgo.Message{
		Key:   key,
		Value: value,
	})
}

// func (k *kafkaProducer) PublishImportJob(ctx context.Context, msg ImportJobMessage) error {
// 	payload, err := json.Marshal(msg)
// 	if err != nil {
// 		return err
// 	}

// 	return k.producer.WriteOne(ctx, kgo.Message{
// 		Key:   []byte(msg.JobID), // important for partitioning
// 		Value: payload,
// 	})
// }
