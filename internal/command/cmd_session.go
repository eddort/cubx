package command

import (
	"github.com/eddort/cubx/internal/config"
	"github.com/eddort/cubx/internal/session"
)

type SessionCommand struct {
	Flags         config.CLI
	Configuration *config.ProgramConfig
}

func (s *SessionCommand) Execute() error {
	aliases := []string{}
	for _, programConfig := range s.Configuration.Programs {
		aliases = append(aliases, programConfig.Name+"='cubx "+programConfig.Name+"'")
	}
	conf := session.Settings{
		Prompt:  "[cubx]$PS1",
		Aliases: aliases,
	}
	return session.Run(conf)
}
