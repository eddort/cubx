package command

import (
	"cubx/internal/config"
	"errors"
)

var ErrCommandNotFound = errors.New("command is nil")

type Command interface {
	Execute() error
}

func Execute(commandArgs []string, flags config.CLI, configuration *config.ProgramConfig) error {
	var command Command
	if flags.ShowConfig != "" {
		command = &ShowConfigCommand{Flags: flags, Configuration: configuration}
	} else if flags.Session {
		command = &SessionCommand{Flags: flags, Configuration: configuration}
	} else if len(commandArgs) > 0 {
		command = &DockerRunCommand{Flags: flags, Configuration: configuration, CommandArgs: commandArgs}
	}

	if command == nil {
		return ErrCommandNotFound
	}

	return command.Execute()
}
