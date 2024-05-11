package command

import "strings"

func nodeCommandHandler(tag string, platform string, additionalArgs []string) (string, string, []string) {
	return "node:" + tag, platform, append([]string{"node"}, additionalArgs...)
}

func foundryCommandHandler(tag string, platform string, additionalArgs []string) (string, string, []string) {
	args := escapeArgs(additionalArgs)
	return "ghcr.io/foundry-rs/foundry:" + tag, platform, []string{strings.Join(args, " ")}
}

func pythonCommandHandler(tag string, platform string, additionalArgs []string) (string, string, []string) {
	return "python:" + tag, platform, append([]string{"python"}, additionalArgs...)
}

func rubyCommandHandler(tag string, platform string, additionalArgs []string) (string, string, []string) {
	return "ruby:" + tag, platform, append([]string{"ruby"}, additionalArgs...)
}
