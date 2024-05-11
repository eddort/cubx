package main

import (
	"ibox/internal/cli"
	"ibox/internal/command"
	"ibox/internal/docker"
)

func main() {

	commandArgs := cli.Parse()

	docImage, command := command.GetDockerImageAndCommand(commandArgs)

	docker.RunImageAndCommand(docImage, command)
}
