package main

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// --------------------
// Discover only containers that needs monitoring and running
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
		All:     false,
	})
}
