package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func ContainsNoHost(args []string) bool {
	for _, v := range args {
		fmt.Println(v)
		if v == "\"--host\"" {
			return false
		}
	}
	return true
}

func EnsureHostContainer(ctx context.Context, cli *client.Client) string {
	filter := filters.NewArgs()
	filter.Add("label", "ibox-host-container=true")
	containers, err := cli.ContainerList(ctx, container.ListOptions{Filters: filter})
	if err != nil {
		panic(err)
	}

	if len(containers) > 0 {
		hostContainerID := containers[0].ID
		// fmt.Printf("Found existing host container: %s\n", hostContainerID)
		return hostContainerID
	}

	containerConfig := &container.Config{
		Image: "alpine",
		Cmd:   []string{"crond", "-f", "-d", "8"},
		Labels: map[string]string{
			"ibox-host-container": "true",
		},
	}
	resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}
	hostContainerID := resp.ID
	// fmt.Printf("Created new host container: %s\n", hostContainerID)

	if err := cli.ContainerStart(ctx, hostContainerID, container.StartOptions{}); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting a ibox-host container:", err)
		os.Exit(1)
	}
	return hostContainerID
}
