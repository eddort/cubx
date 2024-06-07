package command

import (
	"fmt"
	"github.com/eddort/cubx/internal/config"
	"github.com/eddort/cubx/internal/docker"
	"github.com/eddort/cubx/internal/registry"
	"github.com/eddort/cubx/internal/tui"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
)

type DockerRunCommand struct {
	Flags         config.CLI
	Configuration *config.ProgramConfig
	CommandArgs   []string
}

func (s *DockerRunCommand) GetDockerMeta() (string, []string, *config.Settings, error) {
	baseCommand := s.CommandArgs[0]
	additionalArgs := s.CommandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	for _, programConfig := range s.Configuration.Programs {
		if programConfig.Name == commandName {
			if programConfig.Dockerfile != "" {
				imageWithTag := programConfig.Image + ":" + programConfig.Tag
				err := docker.BuildImage(programConfig.Dockerfile, imageWithTag, filepath.Dir(programConfig.Dockerfile))
				if err != nil {
					return "", nil, nil, fmt.Errorf("error while building docker image: %w", err)
				}
			}
			// merge setting with flags
			image, tag, args, err := handleProgram(dockerTag, commandName, additionalArgs, programConfig)
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
				// TODO: move to the validation part
				if programConfig.Dockerfile != "" {
					return "", nil, nil, fmt.Errorf("use of the select flag is not allowed in local builds")
				}
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

func (s *DockerRunCommand) Execute() error {
	docImage, command, settings, err := s.GetDockerMeta()
	if err != nil {
		return err
	}
	return docker.RunImageAndCommand(docImage, command, s.Flags, settings)
}

func handleProgram(tag string, _ string, args []string, programConfig config.Program) (string, string, []string, error) {
	arguments := args

	if programConfig.Command != "" {
		parts, err := shlex.Split(programConfig.Command)
		if err != nil {
			return "", "", nil, fmt.Errorf("error parsing command: %w", err)
		}
		arguments = append(parts, args...)
	}

	if programConfig.Serializer == "string" {
		escArgs := escapeArgs(arguments)
		arguments = []string{strings.Join(escArgs, " ")}
	}

	if programConfig.Tag != "" {
		tag = programConfig.Tag
	}

	return programConfig.Image, tag, arguments, nil
}

// resolveProgramSettings returns the program settings if they exist, otherwise returns the global settings
func resolveProgramSettings(globalSettings, programSettings *config.Settings, hooks *[]config.Hook, additionalArgs []string) (*config.Settings, error) {
	for _, hook := range *hooks {
		escArgs, err := shlex.Split(strings.Join(additionalArgs, " "))
		if err != nil {
			return nil, fmt.Errorf("error parsing hook command: %w", err)
		}

		hookParts, err := shlex.Split(hook.Command)
		if err != nil {
			return nil, fmt.Errorf("error parsing hook command: %w", err)
		}

		argsString := strings.Join(escArgs, " ")
		hookString := strings.Join(hookParts, " ")

		if strings.HasPrefix(argsString, hookString) {
			return &hook.Settings, nil
		}
	}

	if !programSettings.IsEmpty() {
		return programSettings, nil
	}

	return globalSettings, nil
}

func mergeFlagsWithSettings(programSettings *config.Settings, flags config.CLI) (*config.Settings, error) {
	flagsSetting := config.Settings{
		IgnorePaths: flags.FileIgnores,
	}

	merged, err := config.MergeSettings(*programSettings, flagsSetting)
	if err != nil {
		return nil, fmt.Errorf("error merging command config: %w", err)
	}

	return &merged, nil
}
