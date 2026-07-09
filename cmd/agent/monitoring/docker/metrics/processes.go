package metrics

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// --------------------
// Get the number of running processes inside the container. Good to debug in case of process explosions
// --------------------
func GetPIDs(cgroupPath string) float64 {
	data, err := os.ReadFile(filepath.Join(cgroupPath, "pids.current"))
	if err != nil {
		return 0
	}

	val, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	return val
}
