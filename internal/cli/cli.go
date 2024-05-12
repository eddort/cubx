package cli

import (
	"flag"
	"fmt"
	"os"
)

func myUsage() {
	fmt.Println(getHelpMessage())
}

func Parse() []string {
	flag.Usage = myUsage
	flag.Parse()
	commandArgs := flag.Args()
	if len(commandArgs) < 1 {
		fmt.Println(getHelpMessage())
		os.Exit(1)
	}

	return commandArgs
}
