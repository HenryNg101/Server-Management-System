package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	ID   int
	Host string
	Port int
}

type Result struct {
	ServerID int
	Status   int // 1 = up, 0 = down
	Latency  time.Duration
}

// --- MOCK: replace with your DB fetch later ---
func fetchServers() []Server {
	servers := make([]Server, 10000)
	for i := 0; i < 10000; i++ {
		servers[i] = Server{
			ID:   i,
			Host: "localhost", // just for demo
			Port: 10000 + i,
		}
	}
	return servers
}

// --- Core worker pool ---
func runCheckCycle(ctx context.Context, servers []Server) []Result {
	const workerCount = 200 // tune this

	// Buffered channels to store all jobs and outputs
	jobs := make(chan Server, len(servers))
	results := make(chan Result, len(servers))

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
		jobs <- s
	}
	close(jobs)

	// Wait for workers to finish
	wg.Wait()
	close(results)

	// Collect results
	collected := make([]Result, 0, len(servers))
	for r := range results {
		collected = append(collected, r)
	}

	return collected
}

// --- Scheduler ---
func main() {
	// Create a ticker with a channel to tick every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Starting check cycle...")

			// TODO: Use the actual server query functionality here
			servers := fetchServers()

			start := time.Now()
			results := runCheckCycle(ctx, servers)
			elapsed := time.Since(start)

			fmt.Printf("Checked %d servers in %v\n", len(results), elapsed)

			// Just demo output
			up := 0
			for _, r := range results {
				if r.Status == 1 {
					up++
				}
			}
			fmt.Printf("UP: %d, DOWN: %d\n", up, len(results)-up)

		case <-ctx.Done():
			return
		}
	}
}
