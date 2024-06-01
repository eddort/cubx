package command

import (
	"cubx/internal/config"
	"cubx/internal/tui"
	"fmt"
)

type ShowConfigCommand struct {
	Flags         config.CLI
	Configuration *config.ProgramConfig
}

func (c *ShowConfigCommand) Execute() error {
	for _, programConfig := range c.Configuration.Programs {
		if c.Flags.ShowConfig == programConfig.Name {
			tui.PrintColorizedYAML(programConfig)
			return nil
		}
	}
	return fmt.Errorf("not found command: %v", c.Flags.ShowConfig)
}
