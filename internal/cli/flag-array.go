package cli

import (
	"flag"
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func FlagArray(flagName string, desc string) *arrayFlags {
	var arr arrayFlags
	flag.Var(&arr, flagName, desc)
	return &arr
}
