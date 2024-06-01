package cli

import (
	"cubx/internal/config"
	"flag"
	"fmt"
)

func ShowHelpMessage(configuration config.ProgramConfig) {
	fmt.Println(getHelpMessage(configuration))
}

func Parse(configuration config.ProgramConfig) ([]string, config.CLI) {
	flag.Usage = func() {
		fmt.Println(getHelpMessage(configuration))
	}

	IsSelectMode := flag.Bool("select", false, "Interactive selection of the required application version")
	ShowConfig := flag.String("show-config", "", "Show the configuration for the specified command")
	FileIgnores := FlagArray("ignore-path", "Files or dirs to ignore (can be specified multiple times)")
	Session := flag.Bool("session", false, "Start a session in which all programs are available directly")

	flag.Parse()
	commandArgs := flag.Args()

	return commandArgs, config.CLI{IsSelectMode: *IsSelectMode, FileIgnores: *FileIgnores, ShowConfig: *ShowConfig, Session: *Session}
}
