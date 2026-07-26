package agent

type ServerPushStats struct {
	CPUUsageAvg      float64 `json:"cpu_usage_avg"`
	CPUThrottlingAvg float64 `json:"cpu_throttling_avg"`
	CPUPressureAvg   float64 `json:"cpu_pressure_avg"`

	MemoryUsageAvg      float64 `json:"memory_usage_avg"`
	MemoryWorkingSetAvg float64 `json:"memory_working_set_avg"`
	MemoryRSSAvg        float64 `json:"memory_rss_avg"`
	MemoryPressureAvg   float64 `json:"memory_pressure_avg"`

	ReadBPSAvg    float64 `json:"read_bps_avg"`
	WriteBPSAvg   float64 `json:"write_bps_avg"`
	IOPressureAvg float64 `json:"io_pressure_avg"`

	PIDsAvg float64 `json:"pids_avg"`

	OOMEventsTotal float64 `json:"oom_events_total"`
	OOMKillsTotal  float64 `json:"oom_kills_total"`

	Uptime float64 `json:"uptime"`
}
