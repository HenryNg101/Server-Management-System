package runner

import (
	"strings"
	"time"

	"github.com/HenryNg101/docker-monitoring-agent/internal/metrics"
	"github.com/moby/moby/api/types/container"
)

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
