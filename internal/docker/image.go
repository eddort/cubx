package docker

import (
	"context"
	"fmt"
	"ibox/internal/platform"
	"ibox/internal/streams"
	"os"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

func imageExists(ctx context.Context, cli *client.Client, imageName string) (bool, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", imageName)

	images, err := cli.ImageList(ctx, image.ListOptions{Filters: filterArgs})
	if err != nil {
		return false, err
	}

	return len(images) > 0, nil
}

func pullImage(ctx context.Context, cli *client.Client, dockerImage string) error {

	found, err := imageExists(ctx, cli, dockerImage)
	if err != nil {
		// TODO: combine error and return
		fmt.Fprintln(os.Stderr, "Error checking image existence:", err)
		os.Exit(1)
	}

	if found {
		return nil
	}

	platforms := platform.GetPlatforms(dockerImage)
	platformKey := platforms.GetString()

	fmt.Printf("Download image with platform: %s", platformKey)

	pullRes, err := cli.ImagePull(ctx, dockerImage, image.PullOptions{Platform: platformKey})

	if err != nil {
		// TODO: combine error and return
		fmt.Fprintln(os.Stderr, "Error pulling a Docker container:", err)
		os.Exit(1)
	}
	defer pullRes.Close() // Ensure the response body is closed

	return jsonmessage.DisplayJSONMessagesToStream(pullRes, streams.NewOut(), nil)
}
