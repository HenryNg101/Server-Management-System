package main

import (
	"context"
	"log"
	"time"

	"github.com/HenryNg101/server-management-system/internal/app"
)

func runCleanupJob(app *app.Application, ctx context.Context) {
	log.Println("Running import jobs cleanup...")

	err := app.DataTransferService.CleanupOldImportJobs(ctx, 7*24*time.Hour)
	if err != nil {
		log.Println("cleanup job error:", err)
	}
}
