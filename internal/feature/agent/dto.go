package agent

import "time"

type ContainerMetric struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`

	CPUUsage      float64 `json:"cpu_usage"`
	CPUThrottling float64 `json:"cpu_throttling"`
	CPUPressure   float64 `json:"cpu_pressure"`

	MemoryUsage      float64 `json:"memory_usage"`
	MemoryWorkingSet float64 `json:"memory_working_set"`
	MemoryRSS        float64 `json:"memory_rss"`
	MemoryPressure   float64 `json:"memory_pressure"`

	OOMEvents float64 `json:"oom_events"`
	OOMKills  float64 `json:"oom_kills"`

	PIDs float64 `json:"pids"`

	ReadBPS    float64 `json:"read_bps"`
	WriteBPS   float64 `json:"write_bps"`
	IOPressure float64 `json:"io_pressure"`
}

type Payload struct {
	ServerID  int               `json:"server_id"`
	Timestamp string            `json:"timestamp"`
	Metrics   []ContainerMetric `json:"metrics"`
}

type MetricMessage struct {
	Timestamp     time.Time `json:"@timestamp"`
	ServerID      int       `json:"server_id"`
	ContainerName string    `json:"container_name"`

	Status bool `json:"status"`

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
