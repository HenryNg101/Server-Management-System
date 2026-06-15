package main

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

// --------------------
// Calculate IO throughput of the container
// We also use delta here, to calculate average read/write, to see how many bytes per second got read/written
// --------------------
func (t *IOTracker) GetIO(cgroupPath string) (float64, float64) {
	data, err := os.ReadFile(filepath.Join(cgroupPath, "io.stat"))
	if err != nil {
		return 0, 0
	}

	// Calculate read and write bytes
	var readBytes, writeBytes float64

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		for _, p := range parts {
			if strings.HasPrefix(p, "rbytes=") {
				val := strings.TrimPrefix(p, "rbytes=")
				readBytes, _ = strconv.ParseFloat(val, 64)
			}
			if strings.HasPrefix(p, "wbytes=") {
				val := strings.TrimPrefix(p, "wbytes=")
				writeBytes, _ = strconv.ParseFloat(val, 64)
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

	return deltaRead / deltaTime, deltaWrite / deltaTime // bytes/sec
}

func getIOPressure(cgroupPath string) float64 {
	return getResourcePressure(cgroupPath, "io")
}
