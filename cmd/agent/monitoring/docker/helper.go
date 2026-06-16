package main

import (
	"context"
	"net"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// --------------------
// HEALTH CHECK (simple TCP)
// --------------------
func isPortOpen(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// --------------------
// Discover only containers that needs monitoring
// --------------------
func listContainers() ([]container.Summary, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	f := filters.NewArgs()
	f.Add("label", "monitor=true")

	return cli.ContainerList(context.Background(), container.ListOptions{
		Filters: f,
	})
}
