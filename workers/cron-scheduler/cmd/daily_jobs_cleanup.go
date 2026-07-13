package main

import (
	"context"
	"log"
	"time"

	"github.com/HenryNg101/cron-scheduler/internal/app"
)

func runCleanupJob(app *app.Application, ctx context.Context) {
	log.Println("Running import jobs cleanup...")

	err := app.JobService.CleanupOldImportJobs(ctx, 7*24*time.Hour)
	if err != nil {
		log.Println("cleanup job error:", err)
	}
}
