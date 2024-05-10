package command

import "strings"

func nodeCommandHandler(tag string, additionalArgs []string) (string, []string) {
	return "node:" + tag, append([]string{"node"}, additionalArgs...)
}

func foundryCommandHandler(tag string, additionalArgs []string) (string, []string) {
	args := escapeArgs(additionalArgs)
	return "ghcr.io/foundry-rs/foundry:" + tag, []string{strings.Join(args, " ")}
}

func pythonCommandHandler(tag string, additionalArgs []string) (string, []string) {
	return "python:" + tag, append([]string{"python"}, additionalArgs...)
}

func rubyCommandHandler(tag string, additionalArgs []string) (string, []string) {
	return "ruby:" + tag, append([]string{"ruby"}, additionalArgs...)
}
