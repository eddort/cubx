package config

type CLI struct {
	IsSelectMode bool
	FileIgnores  []string
}

type Program struct {
	Name string `mapstructure:"name" validate:"required"`
	// Aliases     []string `mapstructure:"aliases"`
	Image       string `mapstructure:"image" validate:"required"`
	Command     string `mapstructure:"command"`
	Serializer  string `yaml:"serializer" default:"default"`
	Description string `mapstructure:"description"`
}

type ProgramConfig struct {
	Programs []Program `mapstructure:"programs" validate:"required,dive"`
}
