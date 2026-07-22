package main

import (
	"context"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// --------------------
// Discover only containers that need monitoring and are running
// --------------------
func listContainers() ([]container.Summary, error) {
	// Create Docker API client
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	f := make(client.Filters)
	f.Add("label", "monitor=true")
	containers, err := cli.ContainerList(context.Background(), client.ContainerListOptions{
		Filters: f,
		All:     false,
	})
	if err != nil {
		return nil, err
	}

	return containers.Items, nil
}
