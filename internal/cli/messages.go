package cli

import (
	"flag"
	"fmt"
	"github.com/eddort/cubx/internal/config"
	"github.com/eddort/cubx/internal/tui"
	"sort"
	"strings"
)

func getHelpMessage(configuration config.ProgramConfig) string {
	header := "cubx - Isolated App Launch Made Easy"
	description := "A program to launch applications in isolated Docker containers"
	flags := []struct {
		Flag        string
		Description string
	}{
		{"help", "Displays this help information"},
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s===== %s =====%s\n", tui.ColorBlue, header, tui.ColorReset))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", description))
	sb.WriteString(fmt.Sprintf("%sFlags%s\n", tui.ColorPurple, tui.ColorReset))
	sb.WriteString("\n")

	for _, f := range flags {
		sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", tui.ColorGreen, "--"+f.Flag, tui.ColorReset, tui.ColorYellow, f.Description, tui.ColorReset))
	}
	flag.VisitAll(func(f *flag.Flag) {
		sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", tui.ColorGreen, "--"+f.Name, tui.ColorReset, tui.ColorYellow, f.Usage, tui.ColorReset))
	})

	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("%sCommands%s\n", tui.ColorPurple, tui.ColorReset))

	categories := make(map[string][]config.Program)
	for _, cmd := range configuration.Programs {
		category := cmd.Category
		if category == "" {
			category = "Default"
		}
		categories[category] = append(categories[category], cmd)
	}

	sortedCategories := make([]string, 0, len(categories))
	for category := range categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Strings(sortedCategories)
	for _, category := range sortedCategories {
		if category == "Hidden" {
			continue
		}
		commands := categories[category]
		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Name < commands[j].Name
		})
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%s%s:%s\n", tui.ColorCyan, category, tui.ColorReset))
		sb.WriteString("\n")
		for _, cmd := range commands {
			sb.WriteString(fmt.Sprintf("%s%-15s%s - %s%s%s\n", tui.ColorGreen, cmd.Name, tui.ColorReset, tui.ColorYellow, cmd.Description, tui.ColorReset))
		}
	}
	return sb.String()

}
