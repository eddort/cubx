package cli

import (
	"flag"
	"fmt"
	"ibox/internal/config"
	"os"
)

func myUsage() {
	fmt.Println(getHelpMessage())
}

func Parse() ([]string, config.CLI) {
	flag.Usage = myUsage

	IsSelectMode := flag.Bool("select", false, "Interactive selection of the required application version")

	FileIgnores := FlagArray("ignore-file", "Files to ignore (can be specified multiple times)")

	flag.Parse()
	commandArgs := flag.Args()

	if len(commandArgs) < 1 {
		fmt.Println(getHelpMessage())
		os.Exit(1)
	}

	return commandArgs, config.CLI{IsSelectMode: *IsSelectMode, FileIgnores: *FileIgnores}
}
