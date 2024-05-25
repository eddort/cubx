package config

type CLI struct {
	IsSelectMode bool
	FileIgnores  []string
}

type Hook struct {
	Command  string   `mapstructure:"command" validate:"required"`
	Settings Settings `mapstructure:"settings"`
}

type Program struct {
	Name        string   `mapstructure:"name" validate:"required"`
	Image       string   `mapstructure:"image" validate:"required"`
	Command     string   `mapstructure:"command"`
	Serializer  string   `mapstructure:"serializer" validate:"oneof='' default string testhandler"`
	Description string   `mapstructure:"description"`
	Category    string   `mapstructure:"category"`
	Hooks       []Hook   `mapstructure:"hooks" validate:"dive"`
	Settings    Settings `mapstructure:"settings"`
}

type Settings struct {
	Net         string   `mapstructure:"net" validate:"oneof='' none host bridge"`
	IgnorePaths []string `mapstructure:"ignore_paths"`
}

type ProgramConfig struct {
	Programs []Program `mapstructure:"programs" validate:"required,dive"`
	Settings Settings  `mapstructure:"settings"`
}
