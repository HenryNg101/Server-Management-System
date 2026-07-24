package runner

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/HenryNg101/docker-monitoring-agent/internal/metrics"
	"github.com/HenryNg101/docker-monitoring-agent/internal/security"
	"github.com/moby/moby/api/types/container"
)

func (r *Runner) handleMetrics() {
	messages, err := collectMetrics(r.cfg.ServerID)
	if err != nil {
		log.Println("[Agent] collect error:", err)
		return
	}

	if len(messages) == 0 {
		return
	}

	body, _ := json.Marshal(messages)

	req, err := http.NewRequest(
		http.MethodPost,
		r.cfg.APIURL+"/agents/metrics",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Println("[Agent] request error:", err)
		return
	}

	// Decrypt before send
	key := security.DeriveKey(r.cfg.InstanceID)
	rawKey, err := security.Decrypt(r.sec.APIKeyEncrypted, key)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Agent-API-Key", rawKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		log.Println("[Agent] send failed:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("[Agent] metrics sent:", len(messages))
}

func collectMetrics(serverID int) ([]metrics.MetricMessage, error) {
	containers, err := metrics.ListContainers()
	if err != nil {
		return nil, err
	}

	var messages []metrics.MetricMessage
	cpuTrackers := make(map[string]*metrics.CPUTracker)
	ioTrackers := make(map[string]*metrics.IOTracker)

	for _, c := range containers {
		msg := buildMetricMessage(c, serverID, cpuTrackers, ioTrackers)
		if msg != nil {
			messages = append(messages, *msg)
		}
	}

	return messages, nil
}

func buildMetricMessage(
	c container.Summary,
	serverID int,
	cpuTrackers map[string]*metrics.CPUTracker,
	ioTrackers map[string]*metrics.IOTracker,
) *metrics.MetricMessage {

	name := "unknown"
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}

	cgroupPath := metrics.GetCgroupPath(c.ID)
	if cgroupPath == "" {
		return nil
	}

	cpuTracker := metrics.AddCPUTracker(cpuTrackers, c.ID)
	ioTracker := metrics.AddIOTracker(ioTrackers, c.ID)

	oomEvents, oomKills := metrics.GetOOMEvents(cgroupPath)
	readBPS, writeBPS := ioTracker.GetIO(cgroupPath)

	return &metrics.MetricMessage{
		Timestamp:     time.Now().UTC(),
		ServerID:      serverID,
		ContainerName: name,
		PIDs:          int(metrics.GetPIDs(cgroupPath)),

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
}
