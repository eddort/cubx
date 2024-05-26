package main

import (
	"cubx/internal/cli"
	"cubx/internal/command"
	"cubx/internal/config"
	"cubx/internal/docker"
	"log"
)

func main() {
	configuration, _, err := config.LoadConfig(true)

	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	commandArgs, flags := cli.Parse(*configuration)

	command.HandleShowConfig(flags, configuration)

	docImage, command, settings := command.GetDockerMeta(commandArgs, flags, configuration)

	docker.RunImageAndCommand(docImage, command, flags, settings)
}
