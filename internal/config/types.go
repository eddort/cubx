package config

type CLI struct {
	IsSelectMode bool     `yaml:"is_select_mode"`
	FileIgnores  []string `yaml:"file_ignores"`
}

type Hook struct {
	Command  string   `yaml:"command" validate:"required"`
	Settings Settings `yaml:"settings"`
}

type Program struct {
	Name        string   `yaml:"name" validate:"required"`
	Image       string   `yaml:"image" validate:"required"`
	Command     string   `yaml:"command"`
	Serializer  string   `yaml:"serializer" validate:"oneof='' default string testhandler"`
	Description string   `yaml:"description"`
	DefaultTag  string   `yaml:"default_tag"`
	Category    string   `yaml:"category"`
	Hooks       []Hook   `yaml:"hooks" validate:"dive"`
	Settings    Settings `yaml:"settings"`
}

type Settings struct {
	Net         string   `yaml:"net" validate:"oneof='' none host bridge"`
	IgnorePaths []string `yaml:"ignore_paths"`
}

type ProgramConfig struct {
	Programs []Program `yaml:"programs" validate:"required,dive"`
	Settings Settings  `yaml:"settings"`
}
