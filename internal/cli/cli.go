package cli

import (
	"flag"
	"fmt"
	"os"
)

func Parse() []string {
	flag.Parse()
	commandArgs := flag.Args()
	if len(commandArgs) < 1 {
		fmt.Println("Usage: ibox command [arguments...]")
		os.Exit(1)
	}

	return commandArgs
}
