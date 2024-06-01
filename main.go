package main

import (
	"cubx/internal/cli"
	"cubx/internal/command"
	"cubx/internal/config"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {

	configuration, _, err := config.LoadConfig(true)

	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	commandArgs, flags := cli.Parse(*configuration)

	err = command.Execute(commandArgs, flags, configuration)
	if err != nil {
		if errors.Is(err, command.ErrCommandNotFound) {
			cli.ShowHelpMessage(*configuration)
			os.Exit(0)
		} else {
			fmt.Println("An unexpected error occurred:", err)
		}
	}
}
