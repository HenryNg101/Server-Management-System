package data_transfer

import (
	"context"
	"encoding/json"

	kafkaClient "github.com/HenryNg101/server-management-system/internal/platform/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	PublishImportJob(ctx context.Context, jobID string, objectKey string) error
}

type kafkaProducer struct {
	producer *kafkaClient.Producer
}

func NewKafkaProducer(p *kafkaClient.Producer) KafkaProducer {
	return &kafkaProducer{producer: p}
}

func (k *kafkaProducer) PublishImportJob(ctx context.Context, jobID string, objectKey string) error {
	msg := ImportJobMessage{
		JobID:     jobID,
		ObjectKey: objectKey,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return k.producer.WriteOne(ctx, kafka.Message{
		Key:   []byte(jobID),
		Value: payload,
	})
}
