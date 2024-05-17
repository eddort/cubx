package tui

import (
	"fmt"
	"os"
	"strings"

	bubbletea "github.com/charmbracelet/bubbletea"
)

type myModel struct {
	cursor   int
	choices  []string
	filtered []string
	input    string
	pageSize int
	selected string
}

func initialModel(choices []string) myModel {
	return myModel{
		choices:  choices,
		filtered: choices,
		pageSize: 5,
	}
}

func (m myModel) Init() bubbletea.Cmd {
	return nil
}

func (m myModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, bubbletea.Quit
		case "enter":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				m.selected = m.filtered[m.cursor]
				return m, bubbletea.Quit
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		default:
			if msg.Type == bubbletea.KeyRunes {
				m.input += msg.String()
			}
		}
		m.filtered = filterChoices(m.choices, m.input)
		if m.cursor >= len(m.filtered) {
			m.cursor = max(0, len(m.filtered)-1)
		}
	}
	return m, nil
}

func (m myModel) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Search: %s\n\n", m.input))
	start := max(0, m.cursor-m.pageSize/2)
	end := min(start+m.pageSize, len(m.filtered))
	for i := start; i < end; i++ {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, m.filtered[i]))
	}
	return s.String()
}

func filterChoices(choices []string, input string) []string {
	var filtered []string
	lowerInput := strings.ToLower(input)
	for _, choice := range choices {
		if strings.HasPrefix(strings.ToLower(choice), lowerInput) {
			filtered = append(filtered, choice)
		}
	}
	return filtered
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RunInteractivePrompt(variants []string, defaultVariant string) string {
	p := bubbletea.NewProgram(initialModel(variants))
	mdl, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
	// Type assertion to get a specific type
	if m, ok := mdl.(myModel); ok {
		return m.selected
	}
	return defaultVariant
}
