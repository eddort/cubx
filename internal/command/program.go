package command

import (
	"cubx/internal/config"
	"fmt"
	"strings"

	"github.com/google/shlex"
)

func HandleProgram(tag string, commandName string, args []string, programConfig config.Program) (string, string, []string, error) {
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

	if programConfig.DefaultTag != "" {
		tag = programConfig.DefaultTag
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
