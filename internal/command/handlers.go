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

func (s *DockerCommand) GetDockerMeta() (string, []string, *config.Settings, error) {
	baseCommand := s.CommandArgs[0]
	additionalArgs := s.CommandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	for _, programConfig := range s.Configuration.Programs {
		if programConfig.Name == commandName {
			// merge setting with flags
			image, tag, args, err := HandleProgram(dockerTag, commandName, additionalArgs, programConfig)
			if err != nil {
				return "", nil, nil, fmt.Errorf("error handling program: %w", err)
			}

			settings, err := resolveProgramSettings(&s.Configuration.Settings, &programConfig.Settings, &programConfig.Hooks, additionalArgs)
			if err != nil {
				return "", nil, nil, fmt.Errorf("error resolving program settings: %w", err)
			}
			settingsWithFlags, err := mergeFlagsWithSettings(settings, s.Flags)
			if err != nil {
				return "", nil, nil, fmt.Errorf("error merging flags with settings: %w", err)
			}

			if s.Flags.IsSelectMode {
				// TODO: add loader
				tags, err := registry.FetchTags(image)
				if err != nil {
					return "", nil, nil, fmt.Errorf("error fetching tags: %w", err)
				}
				tag, err = tui.RunInteractivePrompt(tags, "latest")
				if err != nil {
					return "", nil, nil, fmt.Errorf("error processing tags: %w", err)
				}
			}

			return image + ":" + tag, args, settingsWithFlags, nil
		}
	}

	return "ubuntu:" + dockerTag, s.CommandArgs, &s.Configuration.Settings, nil
}

func (s *DockerCommand) Execute() error {
	docImage, command, settings, err := s.GetDockerMeta()
	if err != nil {
		return err
	}
	return docker.RunImageAndCommand(docImage, command, s.Flags, settings)
}
