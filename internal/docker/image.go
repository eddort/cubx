package docker

import (
	"context"
	"fmt"
	"github.com/eddort/cubx/internal/config"
	"github.com/eddort/cubx/internal/platform"
	"github.com/eddort/cubx/internal/streams"

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

func pullImage(ctx context.Context, cli *client.Client, dockerImage string, settings *config.Settings) error {
	found, err := imageExists(ctx, cli, dockerImage)
	if err != nil {
		return fmt.Errorf("error checking image existence: %w", err)
	}

	if found {
		return nil
	}

	platforms, err := platform.GetPlatforms(dockerImage)
	if err != nil {
		return err
	}
	platformKey := platforms.GetString()

	fmt.Printf("Download image with platform: %s\n", platformKey)

	pullRes, err := cli.ImagePull(ctx, dockerImage, image.PullOptions{Platform: settings.Platform})

	if err != nil {
		return fmt.Errorf("error pulling a Docker container: %w", err)
	}
	defer pullRes.Close() // Ensure the response body is closed

	return jsonmessage.DisplayJSONMessagesToStream(pullRes, streams.NewOut(), nil)
}
