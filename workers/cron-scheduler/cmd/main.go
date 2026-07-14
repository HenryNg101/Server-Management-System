package main

import (
	"context"
	"log"
	"time"

	"github.com/HenryNg101/cron-scheduler/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// run once on startup (important)
	// runEmailJob(application, ctx)
	// runCleanupJob(application, ctx)

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runEmailJob(application, ctx)
			runCleanupJob(application, ctx)

		case <-ctx.Done():
			return
		}
	}
}
