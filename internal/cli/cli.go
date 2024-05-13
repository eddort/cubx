package cli

import (
	"flag"
	"fmt"
	"os"
)

func myUsage() {
	fmt.Println(getHelpMessage())
}

type Flags struct {
	IsSelectMode bool
}

func Parse() ([]string, Flags) {
	flag.Usage = myUsage
	IsSelectMode := flag.Bool("select", false, "Activate select mode")
	// flag.Bool("lol", false, "Activate select mode")
	// flag.Bool("a", false, "Activate select mode")
	// flag.Bool("b", false, "Activate select mode")
	// kek := flag.String("select", "latest", "Select")
	flag.Parse()
	commandArgs := flag.Args()
	fmt.Println(commandArgs, "FLAGS", *IsSelectMode)
	if len(commandArgs) < 1 {
		fmt.Println(getHelpMessage())
		os.Exit(1)
	}

	return commandArgs, Flags{IsSelectMode: *IsSelectMode}
}
