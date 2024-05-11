package command

func nodeCommandHandler(tag string, commandName string, additionalArgs []string) (string, []string) {
	return "node:" + tag, append([]string{commandName}, additionalArgs...)
}

func foundryCommandHandler(tag string, commandName string, additionalArgs []string) (string, []string) {
	// args := escapeArgs(additionalArgs)
	return "ghcr.io/foundry-rs/foundry:" + tag, append([]string{commandName}, additionalArgs...)
}

func pythonCommandHandler(tag string, commandName string, additionalArgs []string) (string, []string) {
	return "python:" + tag, append([]string{commandName}, additionalArgs...)
}

func rubyCommandHandler(tag string, commandName string, additionalArgs []string) (string, []string) {
	return "ruby:" + tag, append([]string{commandName}, additionalArgs...)
}
