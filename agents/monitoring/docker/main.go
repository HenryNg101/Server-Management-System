package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/HenryNg101/docker-monitoring-agent/internal/bootstrap"
	"github.com/HenryNg101/docker-monitoring-agent/metrics"
)

func main() {
	agentConfigs := bootstrap.InitiateAgent()

	clientHTTP := &http.Client{Timeout: 5 * time.Second}

	for {
		containers, err := listContainers()
		if err != nil {
			log.Println("error listing containers:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		var messages []metrics.MetricMessage
		cpuTrackers := make(map[string]*metrics.CPUTracker)
		ioTrackers := make(map[string]*metrics.IOTracker)

		for _, c := range containers {
			name := "unknown"
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}

			cgroupPath := metrics.GetCgroupPath(c.ID)
			if cgroupPath == "" {
				continue
			}

			cpuTracker := metrics.AddCPUTracker(cpuTrackers, c.ID)
			ioTracker := metrics.AddIOTracker(ioTrackers, c.ID)

			oomEvents, oomKills := metrics.GetOOMEvents(cgroupPath)
			readBPS, writeBPS := ioTracker.GetIO(cgroupPath)

			msg := metrics.MetricMessage{
				Timestamp:     time.Now().UTC(),
				ServerID:      agentConfigs.ServerID,
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

		body, _ := json.Marshal(messages)

		log.Println("Sending:", string(body))

		req, err := http.NewRequest(http.MethodPost, agentConfigs.APIURL+"/agent/metrics", bytes.NewBuffer(body))
		if err != nil {
			log.Println("request error:", err)
			continue
		}

		req.Header.Set("X-Agent-API-Key", agentConfigs.APIKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := clientHTTP.Do(req)
		if err != nil {
			log.Println("send failed:", err)
		} else {
			resp.Body.Close()
		}

		time.Sleep(10 * time.Second)
	}
}
