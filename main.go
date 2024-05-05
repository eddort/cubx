package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var ports string
	flag.StringVar(&ports, "p", "", "Ports to map in format '3000:3000'")
	flag.Parse()

	commandArgs := flag.Args()
	if len(commandArgs) < 1 {
		fmt.Println("Usage: dock [-p port_mapping] command [arguments...]")
		os.Exit(1)
	}

	image, command := getDockerImageAndCommand(commandArgs)
	portMappings := getPortMappings(ports)

	args := []string{"run", "--platform", "linux/amd64", "--rm", "-it", "-v", fmt.Sprintf("%s:/app", getCurrentDir()), "-w", "/app"}
	args = append(args, portMappings...)
	args = append(args, image)
	args = append(args, command...)

	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running docker command:", err)
		os.Exit(1)
	}
}

func getDockerImageAndCommand(commandArgs []string) (string, []string) {
	baseCommand := commandArgs[0]
	additionalArgs := commandArgs[1:]

	switch baseCommand {
	case "npm", "node":
		return "node:latest", append([]string{baseCommand}, additionalArgs...)
	case "forge", "anvil":
		// Special handling for forge commands
		// Join all command parts into a single string as one argument
		return "ghcr.io/foundry-rs/foundry:latest", []string{strings.Join(commandArgs, " ")}
	case "python", "pip":
		return "python:latest", append([]string{baseCommand}, additionalArgs...)
	case "ruby", "gem":
		return "ruby:latest", append([]string{baseCommand}, additionalArgs...)
	default:
		return "ubuntu:latest", commandArgs
	}
}

func getPortMappings(ports string) []string {
	if ports == "" {
		return nil
	}
	return []string{"-p", ports}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get current directory:", err)
		os.Exit(1)
	}
	return dir
}
