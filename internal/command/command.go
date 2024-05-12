package command

type CommandHandler func(string, string, []string) (string, []string)

type Command struct {
	Name           string
	CommandHandler CommandHandler
	Description    string
}

var CommandHandlers = []Command{
	{"npm", nodeCommandHandler, "Handle Node package manager operations"},
	{"node", nodeCommandHandler, "Execute Node.js programs"},
	{"yarn", nodeCommandHandler, "Manage Node.js packages with Yarn"},
	{"npx", nodeCommandHandler, "Execute Node package binaries"},
	{"forge", foundryCommandHandler, "Interact with smart contracts via Forge"},
	{"cast", foundryCommandHandler, "Send transactions or query blockchain state with Cast"},
	{"anvil", foundryCommandHandler, "Run a local Ethereum node using Anvil"},
	{"python", pythonCommandHandler, "Execute Python scripts"},
	{"pip", pythonCommandHandler, "Manage Python packages with pip"},
	{"ruby", rubyCommandHandler, "Execute Ruby scripts"},
	{"gem", rubyCommandHandler, "Manage Ruby gems"},
}

func GetDockerImageAndCommand(commandArgs []string) (string, []string) {
	baseCommand := commandArgs[0]
	additionalArgs := commandArgs[1:]

	commandName, dockerTag := parseBaseCommand(baseCommand)

	for _, command := range CommandHandlers {
		if command.Name == commandName {
			return command.CommandHandler(dockerTag, commandName, additionalArgs)
		}
	}

	return "ubuntu:" + dockerTag, commandArgs
}
