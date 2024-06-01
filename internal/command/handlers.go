package command

import (
	"cubx/internal/config"
	"cubx/internal/docker"
	"cubx/internal/registry"
	"cubx/internal/session"
	"cubx/internal/tui"
	"fmt"
)

type ShowConfigCommand struct {
	Flags         config.CLI
	Configuration *config.ProgramConfig
}

func (c *ShowConfigCommand) Execute() error {
	for _, programConfig := range c.Configuration.Programs {
		if c.Flags.ShowConfig == programConfig.Name {
			tui.PrintColorizedYAML(programConfig)
			return nil
		}
	}
	return fmt.Errorf("not found command: %v", c.Flags.ShowConfig)
}

type SessionCommand struct {
	Flags         config.CLI
	Configuration *config.ProgramConfig
}

func (s *SessionCommand) Execute() error {
	aliases := []string{}
	for _, programConfig := range s.Configuration.Programs {
		aliases = append(aliases, programConfig.Name+"='cubx "+programConfig.Name+"'")
	}
	conf := session.Settings{
		Prompt:  "[cubx]$PS1",
		Aliases: aliases,
	}
	return session.Run(conf)
}

type DockerCommand struct {
	Flags         config.CLI
	Configuration *config.ProgramConfig
	CommandArgs   []string
}

func (s *DockerCommand) GetDockerMeta() (string, []string, *config.Settings) {
	baseCommand := s.CommandArgs[0]
	additionalArgs := s.CommandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	for _, programConfig := range s.Configuration.Programs {
		if programConfig.Name == commandName {
			// merge setting with flags
			image, tag, args := HandleProgram(dockerTag, commandName, additionalArgs, programConfig)

			settings := resolveProgramSettings(&s.Configuration.Settings, &programConfig.Settings, &programConfig.Hooks, additionalArgs)
			settingsWithFlags := mergeFlagsWithSettings(settings, s.Flags)

			if s.Flags.IsSelectMode {
				// TODO: add loader
				tags, _ := registry.FetchTags(image)
				tag = tui.RunInteractivePrompt(tags, "latest")
			}

			return image + ":" + tag, args, settingsWithFlags
		}
	}

	return "ubuntu:" + dockerTag, s.CommandArgs, &s.Configuration.Settings
}

func (s *DockerCommand) Execute() error {
	docImage, command, settings := s.GetDockerMeta()
	docker.RunImageAndCommand(docImage, command, s.Flags, settings)
	// TODO: refactor RunImageAndCommand
	return nil
}
