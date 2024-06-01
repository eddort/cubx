package config

var defaultPrograms = []Program{
	{Name: "npm", Image: "node", Command: "npm", Description: "Handle Node package manager operations", Category: "Node.js", Tag: "latest"},
	{Name: "node", Image: "node", Command: "node", Description: "Execute Node.js programs", Category: "Node.js", Tag: "latest"},
	{Name: "yarn", Image: "node", Command: "yarn", Description: "Manage Node.js packages with Yarn", Category: "Node.js", Tag: "latest"},
	{Name: "npx", Image: "node", Command: "npx", Description: "Execute Node package binaries", Category: "Node.js", Tag: "latest"},
	{Name: "python", Image: "python", Command: "python", Description: "Execute Python scripts", Category: "Python", Tag: "latest"},
	{Name: "ruff", Image: "ghcr.io/astral-sh/ruff", Description: "Python linter and code formatter, written in Rust.", Category: "Python", Tag: "latest"},
	{Name: "pip", Image: "python", Command: "pip", Description: "Manage Python packages with pip", Category: "Python", Tag: "latest"},
	{Name: "ruby", Image: "ruby", Command: "ruby", Description: "Execute Ruby scripts", Category: "Ruby", Tag: "latest"},
	{Name: "gem", Image: "ruby", Command: "gem", Description: "Manage Ruby gems", Category: "Ruby", Tag: "latest"},
}

func getProgramConfig() *ProgramConfig {
	return &ProgramConfig{Programs: defaultPrograms}
}
