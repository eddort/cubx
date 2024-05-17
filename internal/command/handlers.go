package command

import "strings"

func nodeCommandHandler(tag string, commandName string, additionalArgs []string) (string, string, []string) {
	return "node", tag, append([]string{commandName}, additionalArgs...)
}

func foundryCommandHandler(tag string, commandName string, additionalArgs []string) (string, string, []string) {
	args := escapeArgs(additionalArgs)
	fullCommand := append([]string{commandName}, args...)
	return "ghcr.io/foundry-rs/foundry", tag, []string{strings.Join(fullCommand, " ")}
}

func pythonCommandHandler(tag string, commandName string, additionalArgs []string) (string, string, []string) {
	return "python", tag, append([]string{commandName}, additionalArgs...)
}

func rubyCommandHandler(tag string, commandName string, additionalArgs []string) (string, string, []string) {
	return "ruby", tag, append([]string{commandName}, additionalArgs...)
}
