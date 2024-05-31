package session

import "fmt"

type Settings struct {
	Prompt  string
	Aliases []string
}

func (s *Settings) SettingsToStrings() []string {
	var settings []string
	if s.Prompt != "" {
		settings = append(settings, fmt.Sprintf("export PS1=\"%s\"\n", s.Prompt))
	}
	for _, alias := range s.Aliases {
		settings = append(settings, fmt.Sprintf("alias %s\n", alias))
	}
	return settings
}
