package command

type CommandHandler func(string, string, []string) (string, []string)

var commandHandlers = map[string]CommandHandler{
	"npm":    nodeCommandHandler,
	"node":   nodeCommandHandler,
	"yarn":   nodeCommandHandler,
	"npx":    nodeCommandHandler,
	"forge":  foundryCommandHandler,
	"cast":   foundryCommandHandler,
	"anvil":  foundryCommandHandler,
	"python": pythonCommandHandler,
	"pip":    pythonCommandHandler,
	"ruby":   rubyCommandHandler,
	"gem":    rubyCommandHandler,
}

func GetDockerImageAndCommand(commandArgs []string) (string, []string) {
	baseCommand := commandArgs[0]
	additionalArgs := commandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	if handler, ok := commandHandlers[commandName]; ok {
		return handler(dockerTag, commandName, additionalArgs)
	}

	return "ubuntu:" + dockerTag, commandArgs
}
