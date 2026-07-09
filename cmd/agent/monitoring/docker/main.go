package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HenryNg101/server-management-system/cmd/agent/monitoring/docker/metrics"
	"github.com/HenryNg101/server-management-system/internal/feature/agent"
)

// --------------------
// MAIN
// --------------------
func main() {
	serverIDStr := os.Getenv("SERVER_ID")
	apiURL := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")

	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		log.Fatal(err)
	}

	clientHTTP := &http.Client{Timeout: 5 * time.Second}

	for {
		containers, err := listContainers()
		if err != nil {
			log.Println("error listing containers:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		var messages []agent.MetricMessage
		cpuTrackers := make(map[string]*metrics.CPUTracker)
		ioTrackers := make(map[string]*metrics.IOTracker)

		for _, c := range containers {
			name := "unknown"
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}

			// ---- cgroup path
			cgroupPath := metrics.GetCgroupPath(c.ID)
			if cgroupPath == "" {
				continue
			}

			cpuTracker := metrics.AddCPUTracker(cpuTrackers, c.ID)
			ioTracker := metrics.AddIOTracker(ioTrackers, c.ID)

			// Dealing with metrics
			oomEvents, oomKills := metrics.GetOOMEvents(cgroupPath)
			readBPS, writeBPS := ioTracker.GetIO(cgroupPath)
			msg := agent.MetricMessage{
				Timestamp:     time.Now().UTC(),
				ServerID:      serverID,
				ContainerName: name,

				PIDs: int(metrics.GetPIDs(cgroupPath)),

				CPU: struct {
					Usage      float64 `json:"usage"`
					Throttling float64 `json:"throttling"`
					Pressure   float64 `json:"pressure"`
				}{
					Usage:      cpuTracker.GetCPUPercent(cgroupPath),
					Throttling: metrics.GetCPUThrottling(cgroupPath),
					Pressure:   metrics.GetCPUPressure(cgroupPath),
				},

				Memory: struct {
					Usage      float64 `json:"usage"`
					WorkingSet float64 `json:"working_set"`
					RSS        float64 `json:"rss"`
					Pressure   float64 `json:"pressure"`
				}{
					Usage:      metrics.GetContainerMemory(cgroupPath),
					WorkingSet: metrics.GetMemoryWorkingSet(cgroupPath),
					RSS:        metrics.GetMemoryRSS(cgroupPath),
					Pressure:   metrics.GetMemoryPressure(cgroupPath),
				},

				IO: struct {
					ReadBPS  float64 `json:"read_bps"`
					WriteBPS float64 `json:"write_bps"`
					Pressure float64 `json:"pressure"`
				}{
					ReadBPS:  readBPS,
					WriteBPS: writeBPS,
					Pressure: metrics.GetIOPressure(cgroupPath),
				},

				OOM: struct {
					Events int `json:"events"`
					Kills  int `json:"kills"`
				}{
					Events: int(oomEvents),
					Kills:  int(oomKills),
				},
			}
			messages = append(messages, msg)
		}

		//
		// Sending metrics to HTTP server
		body, _ := json.Marshal(messages)

		fmt.Println("Sending:", string(body))

		req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(body))
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("X-Agent-API-Key", apiKey)

		resp, err := clientHTTP.Do(req)
		if err != nil {
			log.Println("failed:", err)
		} else {
			resp.Body.Close()
		}

		time.Sleep(10 * time.Second)
	}
}
