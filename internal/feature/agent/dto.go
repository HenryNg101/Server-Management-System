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
