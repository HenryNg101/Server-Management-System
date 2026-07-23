package runner

import (
	"log"

	"github.com/HenryNg101/docker-monitoring-agent/internal/bootstrap"
)

func (r *Runner) handleRotation() {
	log.Println("[Agent] rotating API key...")

	if err := bootstrap.RotateKey(r.cfg, r.sec); err != nil {
		log.Println("[Agent] rotate failed:", err)
		return
	}

	log.Println("[Agent] API key rotated")
}
