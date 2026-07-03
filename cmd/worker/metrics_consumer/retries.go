package main

import (
	"context"
	"log"
	"time"
)

// Retries ES bulk insert for N times before send to DLQ. This is to handle transient ES failures, e.g. network issues, ES cluster being busy, etc.
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
