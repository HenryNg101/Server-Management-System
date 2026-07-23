package runner

import (
	"net/http"
	"time"

	"github.com/HenryNg101/docker-monitoring-agent/internal/bootstrap"
)

const (
	metricsInterval  = 10 * time.Second
	rotationInterval = 24 * time.Hour // change for demo
)

type Runner struct {
	cfg    *bootstrap.Config
	sec    *bootstrap.Secret
	client *http.Client
}

func NewRunner(cfg *bootstrap.Config, sec *bootstrap.Secret, client *http.Client) *Runner {
	return &Runner{
		cfg:    cfg,
		sec:    sec,
		client: client,
	}
}

func (r *Runner) Start() {
	metricsTicker := time.NewTicker(metricsInterval)
	rotationTicker := time.NewTicker(rotationInterval)

	defer metricsTicker.Stop()
	defer rotationTicker.Stop()

	for {
		select {
		case <-metricsTicker.C:
			r.handleMetrics()

		case <-rotationTicker.C:
			r.handleRotation()
		}
	}
}
