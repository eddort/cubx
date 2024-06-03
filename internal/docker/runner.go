package docker

import (
	"context"
	"cubx/internal/config"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func RunImageAndCommand(dockerImage string, command []string, config config.CLI, settings *config.Settings) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating Docker client: %w", err)
	}

	ctx := context.Background()

	err = pullImage(ctx, cli, dockerImage)

	if err != nil {
		return fmt.Errorf("error pulling a Docker container: %w", err)
	}

	// hostContainerId := EnsureHostContainer(ctx, cli)

	cwd, err := getCurrentDir()
	if err != nil {
		return err
	}
	ogCWD := ""
	containerENVS := []string{}
	hostCwd := os.Getenv("CUBX_HOST_CWD")
	// isPrivileged := false

	if hostCwd == "" {
		ogCWD = cwd
	} else {
		// isPrivileged = true
		ogCWD = hostCwd

	}
	containerENVS = append(containerENVS, fmt.Sprintf("CUBX_HOST_CWD=%s", ogCWD))
	// fmt.Println("isPrivileged", isPrivileged)
	dockerContainerConfig := &container.Config{
		Image:        dockerImage,
		Cmd:          command,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/app",
		OpenStdin:    true, // Ensure the stdin is open
		StdinOnce:    false,
		Env:          containerENVS,
		// ExposedPorts: exposedPorts,
		// Labels: ["cubx-container"]
	}

	mounts, err := generateMounts(ogCWD, settings.IgnorePaths, settings.Mounts)
	if err != nil {
		return fmt.Errorf("generate mounts error: %w", err)
	}

	dockerHostConfig := &container.HostConfig{
		// AutoRemove: true,
		// NetworkMode:  container.NetworkMode("container:" + hostContainerId),
		NetworkMode: "host",
		// PortBindings: portMappings,
		Mounts: mounts,
		// Privileged: isPrivileged,
	}

	if settings.Net != "" {
		dockerHostConfig.NetworkMode = container.NetworkMode(settings.Net)
	}

	resp, err := cli.ContainerCreate(ctx, dockerContainerConfig, dockerHostConfig, nil, nil, "")

	if err != nil {
		return fmt.Errorf("error creating a Docker container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("error starting a Docker container: %w", err)
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
		return fmt.Errorf("error connecting to the container: %w", err)
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
		// fmt.Println("Completion signal received, stop and delete the container...")
		return cleanUpContainer(cli, ctx, resp.ID)

	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container completion: %w", err)
		}

	case <-statusCh:
		// fmt.Printf("The container has completed its work with the status %d\n", status.StatusCode)
		return cleanUpContainer(cli, ctx, resp.ID)
	}

	return nil
}

func getCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve the current directory: %w", err)
	}
	return dir, nil
}

func cleanUpContainer(cli *client.Client, ctx context.Context, containerID string) error {
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop the container: %w", err)
	}
	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to delete the container: %w", err)
	}
	return nil
}

// createTempDir creates a temporary empty directory
func createTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "empty")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %w", err)
	}
	return tempDir, nil
}

func parseMountsString(input string) mount.Mount {
	parts := strings.Split(input, ":")
	if len(parts) == 1 {
		return mount.Mount{
			Type:   mount.TypeBind,
			Source: parts[0],
			Target: parts[0],
		}
	}
	return mount.Mount{
		Type:   mount.TypeBind,
		Source: parts[0],
		Target: parts[1],
	}
}

func generateMounts(cwd string, ignores []string, user_mounts []string) ([]mount.Mount, error) {
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: cwd,
			Target: "/app",
		},
		// TODO: add volume to config and flags
		// {
		// 	Type:   mount.TypeBind,
		// 	Source: "/var/run/docker.sock",
		// 	Target: "/var/run/docker.sock",
		// },
	}

	for _, s := range user_mounts {
		mount := parseMountsString(s)
		mounts = append(mounts, mount)
	}

	for _, ignore := range ignores {
		// Convert the path to an absolute path
		absPath, err := filepath.Abs(ignore)
		if err != nil {
			return nil, fmt.Errorf("error converting path to absolute: %w", err)
		}

		// Check if the path exists
		fileInfo, err := os.Stat(absPath)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", absPath)
		}

		if err != nil {
			return nil, fmt.Errorf("error stating path: %w", err)
		}

		var source string
		if fileInfo.IsDir() {
			// Create a temporary empty directory
			source, err = createTempDir()
			if err != nil {
				return nil, err
			}
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

	return mounts, nil
}
