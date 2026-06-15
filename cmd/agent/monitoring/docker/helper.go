package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func readKVFile(path string) map[string]float64 {
	result := make(map[string]float64)

	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		val, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}

		result[parts[0]] = val
	}

	return result
}

// --------------------
// Find cgroup path, where all information of CPU, memory usage, etc. are stored
// For context, /sys/fs/cgroup is the root directory of the virtual filesystem to cgroup subsystem, so it contains live info from the kernel
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
// Calculate the amount of time processes are stalled, waiting for new resource becomes available (Either Disk I/O, freed memory, or CPU core)
// Right now, for simplicity, we take only the average time in last 10 seconds (i.e avg10),
// that at least one process stalled ("some", instead of "all"), waiting for resource.
// These info are tracked via Linux kernel's Pressure Stall Information (PSI), a kernel feature to track those info
// --------------------
func getResourcePressure(cgroupPath string, resourceType string) float64 {
	data, err := os.ReadFile(filepath.Join(cgroupPath, fmt.Sprintf("%s.pressure", resourceType)))
	if err != nil {
		return 0
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "some") {
			// example: some avg10=0.10 avg60=0.05 avg300=0.02
			parts := strings.Fields(line)
			for _, p := range parts {
				if strings.HasPrefix(p, "avg10=") {
					val := strings.TrimPrefix(p, "avg10=")
					f, _ := strconv.ParseFloat(val, 64)
					return f
				}
			}
		}
	}

	return 0
}
