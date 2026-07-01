package agent

import "time"

type MetricMessage struct {
	Timestamp     time.Time `json:"@timestamp"`
	ServerID      int       `json:"server_id"`
	ContainerName string    `json:"container_name"`

	CPU struct {
		Usage      float64 `json:"usage"`
		Throttling float64 `json:"throttling"`
		Pressure   float64 `json:"pressure"`
	} `json:"cpu"`

	Memory struct {
		Usage      float64 `json:"usage"`
		WorkingSet float64 `json:"working_set"`
		RSS        float64 `json:"rss"`
		Pressure   float64 `json:"pressure"`
	} `json:"memory"`

	IO struct {
		ReadBPS  float64 `json:"read_bps"`
		WriteBPS float64 `json:"write_bps"`
		Pressure float64 `json:"pressure"`
	} `json:"io"`

	PIDs int `json:"pids"`

	OOM struct {
		Events int `json:"events"`
		Kills  int `json:"kills"`
	} `json:"oom"`
}

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
}
