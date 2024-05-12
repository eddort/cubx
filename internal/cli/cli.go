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
	flag.Bool("kek", false, "Activate select mode")
	flag.Bool("lol", false, "Activate select mode")
	flag.Bool("a", false, "Activate select mode")
	flag.Bool("b", false, "Activate select mode")
	flag.Bool("c", false, "Activate select mode")
	flag.Parse()
	commandArgs := flag.Args()
	fmt.Println(commandArgs)
	if len(commandArgs) < 1 {
		fmt.Println(getHelpMessage())
		os.Exit(1)
	}

	return commandArgs
}
