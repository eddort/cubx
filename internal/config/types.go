package config

import "reflect"

type CLI struct {
	IsSelectMode bool     `yaml:"is_select_mode"`
	FileIgnores  []string `yaml:"file_ignores"`
	ShowConfig   string   `yaml:"show_config"`
	Session      bool     `yaml:"session"`
}

type Hook struct {
	Command  string   `yaml:"command"`
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

func (s *Settings) IsEmpty() bool {
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !isEmptyValue(field) {
			return false
		}
	}
	return true
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Chan, reflect.Ptr, reflect.Interface, reflect.Func:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
