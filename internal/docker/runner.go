package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"

	"ibox/internal/streams"
)

func RunImageAndCommand(dockerImage string, command []string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating Docker client:", err)
		os.Exit(1)
	}

	ctx := context.Background()

	err = pullImage(ctx, cli, dockerImage)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error pulling a Docker container:", err)
		os.Exit(1)
	}

	// hostContainerId := EnsureHostContainer(ctx, cli)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        dockerImage,
		Cmd:          command,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/app",
		OpenStdin:    true, // Ensure the stdin is open
		StdinOnce:    false,
		// ExposedPorts: exposedPorts,
		// Labels: ["ibox-container"]
	}, &container.HostConfig{
		AutoRemove: true,
		// NetworkMode:  container.NetworkMode("container:" + hostContainerId),
		NetworkMode: "host",
		// PortBindings: portMappings,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: getCurrentDir(),
				Target: "/app",
			},
		},
	}, nil, nil, "")

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating a Docker container:", err)
		os.Exit(1)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting a Docker container:", err)
		os.Exit(1)
	}
	// TODO: check if we need this
	// defer cleanUpContainer(cli, ctx, resp.ID)

	out, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to the container:", err)
		os.Exit(1)
	}
	defer out.Close()

	go io.Copy(os.Stdout, out.Reader)
	go io.Copy(out.Conn, os.Stdin)

	// Processing of termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case <-sigCh:
		fmt.Println("Completion signal received, stop and delete the container...")
		cleanUpContainer(cli, ctx, resp.ID)

	case err := <-errCh:
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for container completion:", err)
			os.Exit(1)
		}

	case <-statusCh:
		// fmt.Printf("The container has completed its work with the status %d\n", status.StatusCode)
		cleanUpContainer(cli, ctx, resp.ID)
	}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve the current directory:", err)
		os.Exit(1)
	}
	return dir
}

func cleanUpContainer(cli *client.Client, ctx context.Context, containerID string) {
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stop the container: %v\n", err)
	}
	// if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to delete the container: %v\n", err)
	// }
}

func imageExists(ctx context.Context, cli *client.Client, imageName string) (bool, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", imageName)

	images, err := cli.ImageList(ctx, image.ListOptions{Filters: filterArgs})
	if err != nil {
		return false, err
	}

	return len(images) > 0, nil
}

func pullImage(ctx context.Context, cli *client.Client, docImage string) error {

	found, err := imageExists(ctx, cli, docImage)
	if err != nil {
		// TODO: combine error and return
		fmt.Fprintln(os.Stderr, "Error checking image existence:", err)
		os.Exit(1)
	}

	if found {
		return nil
	}

	pullRes, err := cli.ImagePull(ctx, docImage, image.PullOptions{Platform: "linux/amd64"})

	if err != nil {
		// TODO: combine error and return
		fmt.Fprintln(os.Stderr, "Error pulling a Docker container:", err)
		os.Exit(1)
	}
	defer pullRes.Close() // Ensure the response body is closed

	return jsonmessage.DisplayJSONMessagesToStream(pullRes, streams.NewOut(), nil)
}
