package main

import (
	"path/filepath"
	"time"
)

type CPUTracker struct {
	prevUsage float64
	prevTime  time.Time
}

// --------------------
// Calculate the percentage of CPU time that was used, by this formula:
// x = Total CPU time used
// y = The time where this was checked
// -> (x_now - x_prev) / (y_now - y_prev) * 100
// --------------------
func (t *CPUTracker) GetCPUPercent(cgroupPath string) float64 {
	stats := readKVFile(filepath.Join(cgroupPath, "cpu.stat"))

	usage := stats["usage_usec"] // microseconds
	now := time.Now()

	if t.prevTime.IsZero() {
		t.prevUsage = usage
		t.prevTime = now
		return 0
	}

	deltaUsage := usage - t.prevUsage
	deltaTime := now.Sub(t.prevTime).Seconds()

	t.prevUsage = usage
	t.prevTime = now

	if deltaTime == 0 {
		return 0
	}

	// convert microseconds → seconds
	cpuSeconds := deltaUsage / 1e6

	return (cpuSeconds / deltaTime) * 100
}

// --------------------
// Calculate the CPU throttling rate of the container
// CPU throttling is an event that, when a container try to use more CPU than it's allowed, kernel forces it to wait
// So in cgroup, resources budget are allocated to each container (through settings), for example, when checking a cpu.max file, and you see this:
// cat cpu.max -> 50000 100000 -> Quota are 50,000 µs (50ms), and for the period 100,000 µs (100ms), which can be translated as "This container can use 50ms of CPU only, for every 100ms"
// So, if during the time inside a period, a container used up all of CPU time, it will have to be waiting -> One CPU throttling state in that period
// Basically, we don't want this happening too many times
// --------------------
func getCPUThrottling(cgroupPath string) float64 {
	stats := readKVFile(filepath.Join(cgroupPath, "cpu.stat"))

	nrPeriods := stats["nr_periods"]     // Number of CPU periods that had gone passed
	nrThrottled := stats["nr_throttled"] // Number of CPU periods that the container hits the CPU usage limit and led to CPU throttling

	if nrPeriods == 0 {
		return 0
	}

	return (nrThrottled / nrPeriods) * 100
}

func getCPUPressure(cgroupPath string) float64 {
	return getResourcePressure(cgroupPath, "cpu")
}
