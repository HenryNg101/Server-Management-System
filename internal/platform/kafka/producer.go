package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Send(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}

	err = p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(time.Now().String()),
			Value: data,
		},
	)

	if err != nil {
		log.Println("kafka write error:", err)
	}
}
