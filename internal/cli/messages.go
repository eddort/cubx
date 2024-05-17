package cli

import (
	"flag"
	"fmt"
	"ibox/internal/command"
	"sort"
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
	header := "ibox - Isolated App Launch Made Easy and Out of the Box"
	description := "A program to launch applications in isolated Docker containers"
	flags := []struct {
		Flag        string
		Description string
	}{
		{"help", "Displays this help information"},
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s===== %s =====%s\n", colorBlue, header, colorReset))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", description))
	sb.WriteString(fmt.Sprintf("%sFlags%s\n", colorPurple, colorReset))
	sb.WriteString("\n")

	for _, f := range flags {
		sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", colorGreen, "--"+f.Flag, colorReset, colorYellow, f.Description, colorReset))
	}
	flag.VisitAll(func(f *flag.Flag) {
		sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", colorGreen, "--"+f.Name, colorReset, colorYellow, f.Usage, colorReset))
	})

	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("%sCommands%s\n", colorPurple, colorReset))

	categories := make(map[string][]command.Command)
	for _, cmd := range command.CommandHandlers {
		categories[cmd.Category] = append(categories[cmd.Category], cmd)
	}

	sortedCategories := make([]string, 0, len(categories))
	for category := range categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Strings(sortedCategories)
	for _, category := range sortedCategories {
		commands := categories[category]
		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Name < commands[j].Name
		})
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%s%s:%s\n", colorCyan, category, colorReset))
		sb.WriteString("\n")
		for _, cmd := range commands {
			sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", colorGreen, cmd.Name, colorReset, colorYellow, cmd.Description, colorReset))
		}
	}
	return sb.String()

}
