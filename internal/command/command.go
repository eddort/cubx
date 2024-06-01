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
	var commands []Command

	if flags.ShowConfig != "" {
		commands = append(commands, &ShowConfigCommand{Flags: flags, Configuration: configuration})
	} else if flags.Session {
		commands = append(commands, &SessionCommand{Flags: flags, Configuration: configuration})
	} else if len(commandArgs) > 0 {
		commands = append(commands, &DockerBuildCommand{Flags: flags, Configuration: configuration, CommandArgs: commandArgs})
		commands = append(commands, &DockerRunCommand{Flags: flags, Configuration: configuration, CommandArgs: commandArgs})
	}

	if len(commands) == 0 {
		return ErrCommandNotFound
	}

	for _, command := range commands {
		if err := command.Execute(); err != nil {
			return err
		}
	}

	return nil
}
