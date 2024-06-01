package tui

import (
	"errors"
	"fmt"
)

func PrintError(err error) {
	fmt.Println(ColorRed + "error" + ColorReset + ":")
	for err != nil {
		fmt.Printf("  %v\n", err)
		err = errors.Unwrap(err)
	}
}
