package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v4/mem"
)

type ContainerMetric struct {
	Name        string  `json:"name"`
	Status      bool    `json:"status"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

type Payload struct {
	ServerID  int               `json:"server_id"`
	Timestamp string            `json:"timestamp"`
	Metrics   []ContainerMetric `json:"metrics"`
}

// --------------------
// HEALTH CHECK (simple TCP)
// --------------------
func isPortOpen(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// --------------------
// FIND CGROUP PATH
// --------------------
func getCgroupPath(containerID string) string {
	paths := []string{
		"/host/sys/fs/cgroup/docker/" + containerID,
		"/host/sys/fs/cgroup/system.slice/docker-" + containerID + ".scope",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			fmt.Println("FOUND cgroup path:", p)
			return p
		}
	}

	fmt.Println("NO cgroup path for:", containerID)
	return ""
}

// --------------------
// MEMORY FROM CGROUP
// --------------------
func getContainerMemory(cgroupPath string) float64 {
	usageBytes, err := os.ReadFile(filepath.Join(cgroupPath, "memory.current"))
	if err != nil {
		return 0
	}

	usage, _ := strconv.ParseFloat(strings.TrimSpace(string(usageBytes)), 64)

	limitBytes, err := os.ReadFile(filepath.Join(cgroupPath, "memory.max"))
	if err != nil {
		return 0
	}

	limitStr := strings.TrimSpace(string(limitBytes))

	// fallback to host memory total
	if limitStr == "max" {
		memInfo, err := mem.VirtualMemory()
		if err != nil || memInfo.Total == 0 {
			return 0
		}
		return (usage / float64(memInfo.Total)) * 100
	}

	limit, _ := strconv.ParseFloat(limitStr, 64)
	if limit == 0 {
		return 0
	}

	return (usage / limit) * 100
}

// --------------------
// DISCOVER CONTAINERS
// --------------------
func listContainers() ([]container.Summary, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	f := filters.NewArgs()
	f.Add("label", "monitor=true")

	return cli.ContainerList(context.Background(), container.ListOptions{
		Filters: f,
	})
}

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

		var metrics []ContainerMetric

		for _, c := range containers {
			name := "unknown"
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}

			// ---- cgroup path
			cgroupPath := getCgroupPath(c.ID)
			if cgroupPath == "" {
				continue
			}

			mem := getContainerMemory(cgroupPath)

			// ---- simple health check (best effort)
			status := false
			if len(c.Ports) > 0 {
				port := strconv.Itoa(int(c.Ports[0].PrivatePort))
				addr := name + ":" + port
				status = isPortOpen(addr)
			}

			metrics = append(metrics, ContainerMetric{
				Name:        name,
				Status:      status,
				CPUUsage:    0, // TODO later
				MemoryUsage: mem,
			})
		}

		payload := Payload{
			ServerID:  serverID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Metrics:   metrics,
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
