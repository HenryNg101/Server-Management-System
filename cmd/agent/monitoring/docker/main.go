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

	"github.com/HenryNg101/server-management-system/internal/feature/agent"
	"github.com/HenryNg101/server-management-system/internal/monitoring/metrics"
)

// TODO: Add API usage here, and claiming API key through config somehow
// --------------------
// MAIN
// --------------------
func main() {
	serverIDStr := os.Getenv("SERVER_ID")
	apiURL := os.Getenv("API_URL")

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

		var containerMetrics []agent.ContainerMetric
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

			// ---- simple health check (best effort)
			status := false
			if len(c.Ports) > 0 {
				port := strconv.Itoa(int(c.Ports[0].PrivatePort))
				addr := name + ":" + port
				status = isPortOpen(addr)
			}

			// Dealing with metrics
			oomEvents, oomKills := metrics.GetOOMEvents(cgroupPath)
			readBPS, writeBPS := ioTracker.GetIO(cgroupPath)
			containerMetrics = append(containerMetrics, agent.ContainerMetric{
				Name:   name,
				Status: status,

				CPUUsage:      cpuTracker.GetCPUPercent(cgroupPath),
				CPUThrottling: metrics.GetCPUThrottling(cgroupPath),
				CPUPressure:   metrics.GetCPUPressure(cgroupPath),

				MemoryUsage:      metrics.GetContainerMemory(cgroupPath),
				MemoryWorkingSet: metrics.GetMemoryWorkingSet(cgroupPath),
				MemoryRSS:        metrics.GetMemoryRSS(cgroupPath),
				MemoryPressure:   metrics.GetMemoryPressure(cgroupPath),

				OOMEvents: oomEvents,
				OOMKills:  oomKills,

				PIDs: metrics.GetPIDs(cgroupPath),

				ReadBPS:    readBPS,
				WriteBPS:   writeBPS,
				IOPressure: metrics.GetIOPressure(cgroupPath),
			})
		}

		//
		// Sending metrics to HTTP server
		payload := agent.Payload{
			ServerID:  serverID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Metrics:   containerMetrics,
		}

		body, _ := json.Marshal(payload)

		fmt.Println("Sending:", string(body))

		resp, err := clientHTTP.Post(apiURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Println("failed:", err)
		} else {
			resp.Body.Close()
		}

		time.Sleep(10 * time.Second)
	}
}
