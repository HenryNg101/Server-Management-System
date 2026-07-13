package main

import (
	"context"
	"sync"
	"time"

	"github.com/HenryNg101/servers-ping-checker/internal/model"
)

// --- Core worker pool ---
func runCheckCycle(ctx context.Context, servers []model.Server) []*model.Server {
	const workerCount = 200 // tune this

	// Buffered channels to store all jobs and outputs
	jobs := make(chan *model.Server, len(servers))
	results := make(chan *model.Server, len(servers))

	var wg sync.WaitGroup

	// Start workers
	// Each worker keeps pulling jobs to process from the buffered channel until empty
	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for s := range jobs {
				res := ping(ctx, s)
				results <- res
			}
		}()
	}

	// Send jobs
	for _, s := range servers {
		jobs <- &s
	}
	close(jobs)

	// Wait for workers to finish
	wg.Wait()
	close(results)

	// Collect results
	collected := make([]*model.Server, 0, len(servers))
	for r := range results {
		r.LastUpdated = time.Now()
		collected = append(collected, r)
	}
	return collected
}
