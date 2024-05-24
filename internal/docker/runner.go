package docker

import (
	"context"
	"cubx/internal/config"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func RunImageAndCommand(dockerImage string, command []string, config config.CLI) {
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
		// Labels: ["cubx-container"]
	}, &container.HostConfig{
		AutoRemove: true,
		// NetworkMode:  container.NetworkMode("container:" + hostContainerId),
		NetworkMode: "host",
		// PortBindings: portMappings,
		Mounts: generateMounts(config.FileIgnores),
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

// createTempDir creates a temporary empty directory
func createTempDir() string {
	tempDir, err := os.MkdirTemp("", "empty")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	return tempDir
}

func generateMounts(ignores []string) []mount.Mount {
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: getCurrentDir(),
			Target: "/app",
		},
	}

	for _, ignore := range ignores {
		// Convert the path to an absolute path
		absPath, err := filepath.Abs(ignore)
		if err != nil {
			fmt.Println("Error converting path to absolute:", err)
			os.Exit(1)
		}

		// Check if the path exists
		fileInfo, err := os.Stat(absPath)
		if os.IsNotExist(err) {
			fmt.Printf("Path does not exist: %s\n", absPath)
			continue
		}

		if err != nil {
			fmt.Println("Error stating path:", err)
			os.Exit(1)
		}

		var source string
		if fileInfo.IsDir() {
			// Create a temporary empty directory
			source = createTempDir()
		} else {
			// Use /dev/null for files
			source = "/dev/null"
		}

		mount := mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: "/app/" + filepath.Base(absPath),
		}

		mounts = append(mounts, mount)
	}

	return mounts
}
