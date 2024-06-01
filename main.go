package main

import (
	"cubx/internal/cli"
	"cubx/internal/command"
	"cubx/internal/config"
	"cubx/internal/tui"
	"errors"
	"os"
)

func main() {

	configuration, _, err := config.LoadConfig(true)

	if err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}

	commandArgs, flags := cli.Parse(*configuration)

	err = command.Execute(commandArgs, flags, configuration)
	if err != nil {
		if errors.Is(err, command.ErrCommandNotFound) {
			cli.ShowHelpMessage(*configuration)
			os.Exit(0)
		} else {
			tui.PrintError(err)
			os.Exit(1)
		}
	}
}
