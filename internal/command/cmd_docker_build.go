package command

import (
	"cubx/internal/config"
	"cubx/internal/docker"
	"path/filepath"
)

type DockerBuildCommand struct {
	Flags         config.CLI
	Configuration *config.ProgramConfig
	CommandArgs   []string
}

func (s *DockerBuildCommand) Execute() error {
	baseCommand := s.CommandArgs[0]
	// additionalArgs := s.CommandArgs[1:]

	commandName, _ := parseBaseCommand(baseCommand)
	for _, programConfig := range s.Configuration.Programs {
		if programConfig.Name == commandName && programConfig.Dockerfile != "" {
			imageWithTag := programConfig.Image + ":" + programConfig.DefaultTag
			return docker.BuildImage(programConfig.Dockerfile, imageWithTag, filepath.Dir(programConfig.Dockerfile))
		}
	}

	return nil
}
