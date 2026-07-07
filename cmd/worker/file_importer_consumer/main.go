package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/HenryNg101/server-management-system/internal/app"
	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/feature/data_transfer"
	kafkaClient "github.com/HenryNg101/server-management-system/internal/platform/kafka"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	application, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.LoadKafka()

	consumer := kafkaClient.NewConsumer(
		cfg.Brokers,
		cfg.ServersImportTopic,
		cfg.ServersImportConsumerGroup,
	)

	log.Println("import worker started")

	for {
		msg, err := consumer.Fetch(ctx)
		if err != nil {
			log.Println("fetch error:", err)
			continue
		}

		var jobMsg data_transfer.ImportJobMessage
		// Skip poison message to not waste time
		if err := json.Unmarshal(msg.Value, &jobMsg); err != nil {
			log.Println("invalid message")
			_ = consumer.Commit(ctx, msg)
			continue
		}

		err = application.DataTransferService.ProcessImportJob(ctx, jobMsg)
		if err != nil {
			log.Println("job failed:", err)
		}

		_ = consumer.Commit(ctx, msg)
	}
}
