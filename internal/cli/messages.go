package cli

import (
	"fmt"
	"ibox/internal/command"
	"strings"
)

const (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func getHelpMessage() string {
	header := "ibox - Your Applications Out of the Box"
	description := "A program to launch applications in isolated Docker containers."
	flags := []struct {
		Flag        string
		Description string
	}{
		{"--help", "Displays this help information"},
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s===== %s =====%s\n", colorBlue, header, colorReset))
	sb.WriteString(fmt.Sprintf("%s\n\n", description))
	sb.WriteString(fmt.Sprintf("%sFlags:%s\n", colorPurple, colorReset))
	for _, f := range flags {
		sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", colorGreen, f.Flag, colorReset, colorYellow, f.Description, colorReset))
	}
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("%sCommands:%s\n", colorPurple, colorReset))
	for _, cmd := range command.CommandHandlers {
		sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", colorGreen, cmd.Name, colorReset, colorYellow, cmd.Description, colorReset))
	}

	return sb.String()
}
