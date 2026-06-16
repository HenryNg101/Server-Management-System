package metrics

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/mem"
)

// --------------------
// Memory from cgroup, to track the percentage of memory used, compared to the maximum limit allowed
// --------------------
func GetContainerMemory(cgroupPath string) float64 {
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

	// fallback to host memory's total memory usage
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
// Calculate the memory working set (Working Set Size - WSS), which is the amount of RAM actively being used, exclude reclaimable cache
// This is done based on the same way that cAdvisor does it, by subtracting total memory usage to inactive file cache
// --------------------
func GetMemoryWorkingSet(cgroupPath string) float64 {
	stats := readKVFile(filepath.Join(cgroupPath, "memory.stat"))

	total := stats["anon"] + stats["file"]
	inactive := stats["inactive_file"]

	return total - inactive
}

// --------------------
// Calculate Resident Set Size, which is the actual portion of a process's memory held in physical RAM (i.e the real amount of RAM that's being used)
// This is different from WSS, because this doesn't include active file caches, only direct usage
// In cgroup, this can be calculated approximately by anon memory, which consists of anonymous mappings like heaps and stacks memory
// --------------------
func GetMemoryRSS(cgroupPath string) float64 {
	stats := readKVFile(filepath.Join(cgroupPath, "memory.stat"))
	return stats["anon"]
}

// --------------------
// Track the number of times where Out Of Memory (OOM) events happened, which is good for alerting
// "oom" stats means the times memory pressure occurred, and "oom_kill" means the times processes were actually killed
// --------------------
func GetOOMEvents(cgroupPath string) (float64, float64) {
	stats := readKVFile(filepath.Join(cgroupPath, "memory.events"))

	return stats["oom"], stats["oom_kill"]
}

func GetMemoryPressure(cgroupPath string) float64 {
	return getResourcePressure(cgroupPath, "memory")
}
