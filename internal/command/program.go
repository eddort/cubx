package command

import (
	"cubx/internal/config"
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

	if programConfig.DefaultTag != "" {
		tag = programConfig.DefaultTag
	}

	return programConfig.Image, tag, arguments
}

// resolveProgramSettings returns the program settings if they exist, otherwise returns the global settings
func resolveProgramSettings(globalSettings, programSettings *config.Settings, hooks *[]config.Hook, additionalArgs []string) *config.Settings {
	for _, hook := range *hooks {
		escArgs, err := shlex.Split(strings.Join(additionalArgs, " "))
		if err != nil {
			log.Fatalf("Error parsing hook command: %v", err)
		}

		hookParts, err := shlex.Split(hook.Command)
		if err != nil {
			log.Fatalf("Error parsing hook command: %v", err)
		}

		argsString := strings.Join(escArgs, " ")
		hookString := strings.Join(hookParts, " ")

		if strings.HasPrefix(argsString, hookString) {

			return &hook.Settings
		}
	}

	if !programSettings.IsEmpty() {
		return programSettings
	}

	return globalSettings
}

func mergeFlagsWithSettings(programSettings *config.Settings, flags config.CLI) *config.Settings {
	flagsSetting := config.Settings{
		IgnorePaths: flags.FileIgnores,
	}

	merged, err := config.MergeSettings(*programSettings, flagsSetting)

	if err != nil {
		log.Fatalf("Error merging command config: %v", err)
	}

	return &merged
}
