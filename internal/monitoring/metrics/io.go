package metrics

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type IOTracker struct {
	prevRead  float64
	prevWrite float64
	prevTime  time.Time
}

func AddIOTracker(ioTrackers map[string]*IOTracker, id string) *IOTracker {
	if t, ok := ioTrackers[id]; ok {
		return t
	}
	t := &IOTracker{}
	ioTrackers[id] = t
	return t
}

// --------------------
// Calculate IO throughput of the container, by adding all read and write throughput of block devices in the container
// The only exception is loop device, which has major number start off with 7:, to prevent double counting, as loop devices are virtual filesystems that maps back to file on real disk
// We also use delta here, to calculate average read/write, to see how many bytes per second got read/written
// --------------------
func (t *IOTracker) GetIO(cgroupPath string) (float64, float64) {
	data, err := os.ReadFile(filepath.Join(cgroupPath, "io.stat"))
	if err != nil {
		return 0, 0
	}

	var readBytes, writeBytes float64

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)

		// skip loop devices
		if len(parts) > 0 && strings.HasPrefix(parts[0], "7:") {
			continue
		}

		for _, p := range parts {
			if strings.HasPrefix(p, "rbytes=") {
				val := strings.TrimPrefix(p, "rbytes=")
				v, _ := strconv.ParseFloat(val, 64)
				readBytes += v
			}
			if strings.HasPrefix(p, "wbytes=") {
				val := strings.TrimPrefix(p, "wbytes=")
				v, _ := strconv.ParseFloat(val, 64)
				writeBytes += v
			}
		}
	}

	now := time.Now()

	if t.prevTime.IsZero() {
		t.prevRead = readBytes
		t.prevWrite = writeBytes
		t.prevTime = now
		return 0, 0
	}

	deltaRead := readBytes - t.prevRead
	deltaWrite := writeBytes - t.prevWrite
	deltaTime := now.Sub(t.prevTime).Seconds()

	t.prevRead = readBytes
	t.prevWrite = writeBytes
	t.prevTime = now

	if deltaTime == 0 {
		return 0, 0
	}

	return deltaRead / deltaTime, deltaWrite / deltaTime // bytes / sec
}

func GetIOPressure(cgroupPath string) float64 {
	return getResourcePressure(cgroupPath, "io")
}
