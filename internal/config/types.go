package config

type CLI struct {
	IsSelectMode bool
	FileIgnores  []string
}

type Command struct {
	Name        string   `mapstructure:"name" validate:"required"`
	Aliases     []string `mapstructure:"aliases"`
	Image       string   `mapstructure:"image" validate:"required"`
	Handler     string   `mapstructure:"handler"`
	Description string   `mapstructure:"description"`
}

type CommandConfig struct {
	Commands []Command `mapstructure:"commands" validate:"dive"`
}
