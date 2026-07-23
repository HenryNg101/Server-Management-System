package main

import (
	"log"
	"net/http"
	"time"

	"github.com/HenryNg101/docker-monitoring-agent/internal/bootstrap"
	"github.com/HenryNg101/docker-monitoring-agent/internal/runner"
)

func main() {
	cfg, sec := bootstrap.InitiateAgent()

	client := &http.Client{Timeout: 5 * time.Second}

	r := runner.NewRunner(cfg, sec, client)

	log.Println("[Agent] started")

	r.Start()
}
