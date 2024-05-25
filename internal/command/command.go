package command

import (
	"cubx/internal/config"
	"cubx/internal/registry"
	"cubx/internal/tui"
	"log"
	"strings"

	"github.com/google/shlex"
)

func HandleProgram(tag string, commandName string, args []string, programConfig config.Program) (string, string, []string) {
	arguments := args

	if programConfig.Command != "" {
		parts, err := shlex.Split(programConfig.Command)
		if err != nil {
			log.Fatalf("Error parsing command: %v", err)
		}

		arguments = append(parts, args...)
	}

	if programConfig.Serializer == "string" {
		escArgs := escapeArgs(arguments)
		arguments = []string{strings.Join(escArgs, " ")}
	}

	return programConfig.Image, tag, arguments
}

func isEmptySettings(s *config.Settings) bool {
	return s.Net == "" && len(s.IgnorePaths) == 0
}

// getProgramSettings returns the program settings if they exist, otherwise returns the global settings
func getProgramSettings(globalSettings, programSettings *config.Settings) *config.Settings {
	// fmt.Println(programSettings)
	if !isEmptySettings(programSettings) {
		return programSettings
	}
	return globalSettings
}

func GetDockerImageAndCommand(commandArgs []string, flags config.CLI, configuration *config.ProgramConfig) (string, []string, *config.Settings) {
	baseCommand := commandArgs[0]
	additionalArgs := commandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	for _, programConfig := range configuration.Programs {
		if programConfig.Name == commandName {

			image, tag, args := HandleProgram(dockerTag, commandName, additionalArgs, programConfig)

			settings := getProgramSettings(&configuration.Settings, &programConfig.Settings)

			if flags.IsSelectMode {
				tags, _ := registry.FetchTags(image)
				tag = tui.RunInteractivePrompt(tags, "latest")
			}

			return image + ":" + tag, args, settings
		}
	}

	return "ubuntu:" + dockerTag, commandArgs, &configuration.Settings
}
