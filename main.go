package main

import (
	"ibox/internal/cli"
	"ibox/internal/command"
	"ibox/internal/docker"
)

func main() {
	commandArgs, flags := cli.Parse()
	docImage, command := command.GetDockerImageAndCommand(commandArgs, flags)
	docker.RunImageAndCommand(docImage, command, flags)
}
