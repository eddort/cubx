package cli

import (
	"cubx/internal/config"
	"flag"
	"fmt"
	"os"
)

func Parse(configuration config.ProgramConfig) ([]string, config.CLI) {
	flag.Usage = func() {
		fmt.Println(getHelpMessage(configuration))
	}

	IsSelectMode := flag.Bool("select", false, "Interactive selection of the required application version")
	ShowConfig := flag.String("show-config", "", "Show the configuration for the specified command")
	FileIgnores := FlagArray("ignore-path", "Files or dirs to ignore (can be specified multiple times)")

	flag.Parse()
	commandArgs := flag.Args()

	if len(commandArgs) < 1 && *ShowConfig == "" {
		fmt.Println(getHelpMessage(configuration))
		os.Exit(1)
	}

	return commandArgs, config.CLI{IsSelectMode: *IsSelectMode, FileIgnores: *FileIgnores, ShowConfig: *ShowConfig}
}
