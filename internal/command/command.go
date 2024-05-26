package command

import (
	"cubx/internal/config"
	"cubx/internal/registry"
	"cubx/internal/tui"
	"log"
	"os"
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

func HandleShowConfig(flags config.CLI, configuration *config.ProgramConfig) {

	if flags.ShowConfig == "" {
		return
	}

	for _, programConfig := range configuration.Programs {
		if flags.ShowConfig == programConfig.Name {
			tui.PrintColorizedYAML(programConfig)
			os.Exit(0)
		}
	}

	log.Fatalf("not found command: %v", flags.ShowConfig)
}

func GetDockerMeta(commandArgs []string, flags config.CLI, configuration *config.ProgramConfig) (string, []string, *config.Settings) {
	baseCommand := commandArgs[0]
	additionalArgs := commandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	for _, programConfig := range configuration.Programs {
		if programConfig.Name == commandName {
			// merge setting with flags
			image, tag, args := HandleProgram(dockerTag, commandName, additionalArgs, programConfig)

			settings := resolveProgramSettings(&configuration.Settings, &programConfig.Settings, &programConfig.Hooks, additionalArgs)
			settingsWithFlags := mergeFlagsWithSettings(settings, flags)

			if flags.IsSelectMode {
				tags, _ := registry.FetchTags(image)
				tag = tui.RunInteractivePrompt(tags, "latest")
			}

			return image + ":" + tag, args, settingsWithFlags
		}
	}

	return "ubuntu:" + dockerTag, commandArgs, &configuration.Settings
}
