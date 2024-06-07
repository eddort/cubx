package main

import (
	"errors"
	"github.com/eddort/cubx/internal/cli"
	"github.com/eddort/cubx/internal/command"
	"github.com/eddort/cubx/internal/config"
	"github.com/eddort/cubx/internal/tui"
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
