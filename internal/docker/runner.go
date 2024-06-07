package docker

import (
	"context"
	"cubx/internal/config"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunImageAndCommand(dockerImage string, command []string, config config.CLI, settings *config.Settings) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating Docker client: %w", err)
	}

	ctx := context.Background()

	err = pullImage(ctx, cli, dockerImage, settings)

	if err != nil {
		return fmt.Errorf("error pulling a Docker container: %w", err)
	}

	currentCWD, err := getCWD()
	if err != nil {
		return err
	}

	containerENV := getENV(currentCWD)

	dockerContainerConfig := &container.Config{
		Image:        dockerImage,
		Cmd:          command,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/app",
		OpenStdin:    true,
		StdinOnce:    false,
		Env:          containerENV,
		// ExposedPorts: exposedPorts,
		// Labels: ["cubx-container"]
	}

	mounts, err := generateMounts(currentCWD, settings.IgnorePaths, settings.Mounts)
	if err != nil {
		return fmt.Errorf("generate mounts error: %w", err)
	}

	dockerHostConfig := &container.HostConfig{
		// NetworkMode:  container.NetworkMode("container:" + hostContainerId),
		NetworkMode: "host",
		// TODO: portMappings from config
		// PortBindings: portMappings,
		Mounts: mounts,
	}

	if settings.Net != "" {
		dockerHostConfig.NetworkMode = container.NetworkMode(settings.Net)
	}

	resp, err := cli.ContainerCreate(ctx, dockerContainerConfig, dockerHostConfig, nil, nil, "")

	if err != nil {
		return fmt.Errorf("error creating a Docker container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return cleanUpContainer(cli, ctx, resp.ID, fmt.Errorf("error starting a Docker container: %w", err))
	}

	out, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		return cleanUpContainer(cli, ctx, resp.ID, fmt.Errorf("error connecting to the container: %w", err))
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
		return cleanUpContainer(cli, ctx, resp.ID, nil)

	case err := <-errCh:
		return cleanUpContainer(cli, ctx, resp.ID, fmt.Errorf("error waiting for container completion: %w", err))

	case <-statusCh:
		// fmt.Printf("The container has completed its work with the status %d\n", status.StatusCode)
		return cleanUpContainer(cli, ctx, resp.ID, nil)
	}
}

func cleanUpContainer(cli *client.Client, ctx context.Context, containerID string, customError error) error {
	wrapError := func(err error, msg string) error {
		if customError != nil {
			return fmt.Errorf("%s: %v: %w", customError, msg, err)
		}
		return fmt.Errorf("%s: %w", msg, err)
	}

	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return wrapError(err, "failed to stop the container")
	}

	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return wrapError(err, "failed to delete the container")
	}

	return nil
}

func getCWD() (string, error) {
	cwd, err := getCurrentDir()
	if err != nil {
		return "", err
	}
	ogCWD := ""
	hostCwd := os.Getenv("CUBX_HOST_CWD")

	if hostCwd == "" {
		ogCWD = cwd
	} else {
		ogCWD = hostCwd
	}
	return ogCWD, nil
}

func getENV(currentCWD string) []string {
	containerENVS := []string{}
	termEnv := os.Getenv("TERM")
	if termEnv != "" {
		containerENVS = append(containerENVS, fmt.Sprintf("TERM=%s", termEnv))
	}

	containerENVS = append(containerENVS, fmt.Sprintf("CUBX_HOST_CWD=%s", currentCWD))
	// TODO: pass env from .cubx/config
	return containerENVS
}
