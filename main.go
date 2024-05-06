package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	var ports string
	flag.StringVar(&ports, "p", "", "Ports to map in format '3000:3000'")
	flag.Parse()

	commandArgs := flag.Args()
	if len(commandArgs) < 1 {
		fmt.Println("Usage: ibox [-p port_mapping] command [arguments...]")
		os.Exit(1)
	}

	docImage, command := getDockerImageAndCommand(commandArgs)

	portMappings, err := getPortMappings(ports)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Incorrect port format:", err)
		os.Exit(1)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating Docker client:", err)
		os.Exit(1)
	}

	exposedPorts := nat.PortSet{}
	for port := range portMappings {
		exposedPorts[port] = struct{}{}
	}

	ctx := context.Background()

	err = pullImage(ctx, cli, docImage)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error pulling a Docker container:", err)
		os.Exit(1)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        docImage,
		Cmd:          command,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/app",
		OpenStdin:    true, // Ensure the stdin is open
		StdinOnce:    false,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		AutoRemove: true,

		PortBindings: portMappings,
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

func getDockerImageAndCommand(commandArgs []string) (string, []string) {
	baseCommand := commandArgs[0]
	additionalArgs := commandArgs[1:]
	strings.Split(baseCommand, ":")

	commandName, dockerTag := func() (string, string) {
		x := strings.Split(baseCommand, ":")
		if len(x) > 1 {
			return x[0], x[1]
		}
		return x[0], "latest"
	}()

	switch commandName {
	case "npm", "node", "yarn":
		return "node:" + dockerTag, append([]string{commandName}, additionalArgs...)
	case "forge":
		return "ghcr.io/foundry-rs/foundry:" + dockerTag, []string{strings.Join(commandArgs, " ")}
	case "anvil":
		fullArgs := append(commandArgs, []string{"--host", "0.0.0.0"}...)
		return "ghcr.io/foundry-rs/foundry:" + dockerTag, []string{strings.Join(fullArgs, " ")}
	case "python", "pip":
		return "python:" + dockerTag, append([]string{commandName}, additionalArgs...)
	case "ruby", "gem":
		return "ruby:" + dockerTag, append([]string{commandName}, additionalArgs...)
	default:
		return "ubuntu:" + dockerTag, commandArgs
	}
}

func getPortMappings(ports string) (nat.PortMap, error) {
	if ports == "" {
		return nil, nil
	}
	// Creating struct nat.PortMap
	portMap := nat.PortMap{}

	// Split the string into parts for an individual port and its mapping
	portBindings := strings.Split(ports, ":")
	if len(portBindings) != 2 {
		return nil, fmt.Errorf("invalid port mapping format")
	}

	// Define the port inside the container and the port on the host
	containerPort, err := nat.NewPort("tcp", portBindings[0])
	if err != nil {
		return nil, fmt.Errorf("invalid container port: %v", err)
	}

	// Create a list of port bindings
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: portBindings[1],
	}

	// Add binding to mapping
	portMap[containerPort] = []nat.PortBinding{hostBinding}

	return portMap, nil
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

	decoder := json.NewDecoder(pullRes)
	var message map[string]interface{}
	var prevMsg interface{}
	for decoder.More() {
		err := decoder.Decode(&message)
		if err != nil {
			// TODO: combine error and return
			fmt.Fprintln(os.Stderr, "Error decoding JSON stream:", err)
			os.Exit(1)
		}
		if status, ok := message["status"]; ok {
			if status == prevMsg {
				continue
			}
			prevMsg = status
			fmt.Println(status)
		}
	}

	return nil
}
