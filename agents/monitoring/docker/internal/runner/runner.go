package runner

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/HenryNg101/docker-monitoring-agent/internal/bootstrap"
)

const (
	metricsInterval  = 30 * time.Second
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

// Add a bit of tiny random interval -> Prevent race condition where two actions clashing
func jitter(base time.Duration) time.Duration {
	return base + time.Duration(rand.Int63n(int64(base/5))) // +0~20%
}

func (r *Runner) Start() {
	metricsTicker := time.NewTicker(jitter(metricsInterval))
	rotationTicker := time.NewTicker(jitter(rotationInterval))

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
