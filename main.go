package main

import (
	"cubx/internal/cli"
	"cubx/internal/command"
	"cubx/internal/config"
	"errors"
	"fmt"
	"os"
)

func main() {

	configuration, _, err := config.LoadConfig(true)

	if err != nil {
		fmt.Println("An unexpected error occurred:", err)
		os.Exit(1)
	}

	commandArgs, flags := cli.Parse(*configuration)

	err = command.Execute(commandArgs, flags, configuration)
	if err != nil {
		if errors.Is(err, command.ErrCommandNotFound) {
			cli.ShowHelpMessage(*configuration)
			os.Exit(0)
		} else {
			fmt.Println("An unexpected error occurred:", err)
			os.Exit(1)
		}
	}
}
