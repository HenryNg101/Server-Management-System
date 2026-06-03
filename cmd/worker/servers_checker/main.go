package main

import (
	"context"
	"log"
	"time"

	"github.com/HenryNg101/server-management-system/internal/app"
	"github.com/HenryNg101/server-management-system/internal/feature/server"
	"github.com/HenryNg101/server-management-system/internal/model"
)

func fetchServers(application *app.Application, ctx context.Context) ([]model.Server, error) {
	serverService := *application.ServerService
	paginatedServers, err := serverService.GetServers(ctx, server.GetServersQuery{})
	if err != nil {
		return nil, err
	}
	return paginatedServers.Servers, nil
}

func updateServers(application *app.Application, ctx context.Context, resultServers []*model.Server) error {
	serverService := *application.ServerService

	err := serverService.BulkUpdateServersStatuses(ctx, resultServers)
	if err != nil {
		return err
	}
	return nil
}

// --- Scheduler ---
func main() {
	// Create app
	newApplication, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	// Create a ticker with a channel to tick every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			log.Println("Starting check cycle...")

			// Servers statuses checks
			servers, err := fetchServers(newApplication, ctx)
			if err != nil {
				log.Fatal(err)
			}
			if len(servers) == 0 {
				continue
			}

			start := time.Now()
			results := runCheckCycle(ctx, servers)
			elapsed := time.Since(start)

			log.Printf("Checked %d servers in %v\n", len(results), elapsed)

			//
			// Bulk update statuses to Postgres
			start = time.Now()
			err = updateServers(newApplication, ctx, results)
			if err != nil {
				log.Fatal(err)
			}
			elapsed = time.Since(start)
			log.Printf("Updated %d servers statuses to DB in %v\n", len(results), elapsed)

			//
			// Bulk insert to Elasticsearch data stream
			start = time.Now()
			serverService := *newApplication.ServerService
			err = serverService.ElasticBulkInsert(ctx, results)
			if err != nil {
				log.Fatal(err)
			}
			elapsed = time.Since(start)
			log.Printf("Logged %d servers statuses to Elasticsearch in %v\n", len(results), elapsed)

		case <-ctx.Done():
			return
		}
	}
}
