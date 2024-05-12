package command

type CommandHandler func(string, string, []string) (string, []string)

type Command struct {
	Name           string
	CommandHandler CommandHandler
	Description    string
	Category       string
}

var CommandHandlers = []Command{
	{Name: "npm", CommandHandler: nodeCommandHandler, Description: "Handle Node package manager operations", Category: "Node.js"},
	{Name: "node", CommandHandler: nodeCommandHandler, Description: "Execute Node.js programs", Category: "Node.js"},
	{Name: "yarn", CommandHandler: nodeCommandHandler, Description: "Manage Node.js packages with Yarn", Category: "Node.js"},
	{Name: "npx", CommandHandler: nodeCommandHandler, Description: "Execute Node package binaries", Category: "Node.js"},
	{Name: "forge", CommandHandler: foundryCommandHandler, Description: "Interact with smart contracts via Forge", Category: "Ethereum"},
	{Name: "cast", CommandHandler: foundryCommandHandler, Description: "Send transactions or query blockchain state with Cast", Category: "Ethereum"},
	{Name: "anvil", CommandHandler: foundryCommandHandler, Description: "Run a local Ethereum node using Anvil", Category: "Ethereum"},
	{Name: "python", CommandHandler: pythonCommandHandler, Description: "Execute Python scripts", Category: "Python"},
	{Name: "pip", CommandHandler: pythonCommandHandler, Description: "Manage Python packages with pip", Category: "Python"},
	{Name: "ruby", CommandHandler: rubyCommandHandler, Description: "Execute Ruby scripts", Category: "Ruby"},
	{Name: "gem", CommandHandler: rubyCommandHandler, Description: "Manage Ruby gems", Category: "Ruby"},
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
