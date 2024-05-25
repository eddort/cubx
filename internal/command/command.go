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

func GetDockerImageAndCommand(commandArgs []string, flags config.CLI, configuration *config.ProgramConfig) (string, []string) {
	baseCommand := commandArgs[0]
	additionalArgs := commandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	for _, program := range configuration.Programs {
		if program.Name == commandName {

			image, tag, args := HandleProgram(dockerTag, commandName, additionalArgs, program)

			if flags.IsSelectMode {
				tags, _ := registry.FetchTags(image)
				tag = tui.RunInteractivePrompt(tags, "latest")
			}

			return image + ":" + tag, args
		}
	}

	return "ubuntu:" + dockerTag, commandArgs
}
