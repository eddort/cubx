package main

import (
	"fmt"
	"ibox/internal/cli"
	"ibox/internal/command"
	"ibox/internal/docker"
)

func main() {
	commandArgs, flags := cli.Parse()
	docImage, command := command.GetDockerImageAndCommand(commandArgs)
	fmt.Println(flags.IsSelectMode)
	docker.RunImageAndCommand(docImage, command)
}
