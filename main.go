package main

import (
	"cubx/internal/cli"
	"cubx/internal/command"
	"cubx/internal/docker"
)

func main() {
	commandArgs, flags := cli.Parse()
	docImage, command := command.GetDockerImageAndCommand(commandArgs, flags)
	docker.RunImageAndCommand(docImage, command, flags)
}
