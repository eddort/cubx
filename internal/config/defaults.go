package config

var defaultPrograms = []Program{
	{Name: "npm", Image: "node", Command: "npm", Description: "Handle Node package manager operations", Category: "Node.js"},
	{Name: "node", Image: "node", Command: "node", Description: "Execute Node.js programs", Category: "Node.js"},
	{Name: "yarn", Image: "node", Command: "yarn", Description: "Manage Node.js packages with Yarn", Category: "Node.js"},
	{Name: "npx", Image: "node", Command: "npx", Description: "Execute Node package binaries", Category: "Node.js"},
	{Name: "forge", Image: "ghcr.io/foundry-rs/foundry", Command: "forge", Description: "Interact with smart contracts via Forge", Category: "Ethereum"},
	{Name: "cast", Image: "ghcr.io/foundry-rs/foundry", Command: "cast", Description: "Send transactions or query blockchain state with Cast", Category: "Ethereum"},
	{Name: "anvil", Image: "ghcr.io/foundry-rs/foundry", Command: "anvil", Description: "Run a local Ethereum node using Anvil", Category: "Ethereum"},
	{Name: "python", Image: "python", Command: "python", Description: "Execute Python scripts", Category: "Python"},
	{Name: "ruff", Image: "ghcr.io/astral-sh/ruff", Command: "ruff", Description: "Python linter and code formatter, written in Rust.", Category: "Python"},
	{Name: "pip", Image: "python", Command: "pip", Description: "Manage Python packages with pip", Category: "Python"},
	{Name: "ruby", Image: "ruby", Command: "ruby", Description: "Execute Ruby scripts", Category: "Ruby"},
	{Name: "gem", Image: "ruby", Command: "gem", Description: "Manage Ruby gems", Category: "Ruby"},
}

func getProgramConfig() *ProgramConfig {
	return &ProgramConfig{Programs: defaultPrograms}
}
