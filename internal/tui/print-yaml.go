package tui

import (
	"fmt"
	"log"
	"regexp"

	"gopkg.in/yaml.v3"
)

func colorizeYAML(yamlStr string) string {
	// Define regex patterns for different parts of the YAML
	keyPattern := regexp.MustCompile(`(?m)^[a-zA-Z_][a-zA-Z0-9_]*:`)
	strValuePattern := regexp.MustCompile(`(?m):\s*".*?"|\s*'.*?'`)
	numValuePattern := regexp.MustCompile(`(?m):\s*\d+`)

	// Apply Color formatting to the YAML output
	resYaml := keyPattern.ReplaceAllStringFunc(yamlStr, func(s string) string {
		return ColorCyan + s + ColorReset
	})
	resYaml = strValuePattern.ReplaceAllStringFunc(resYaml, func(s string) string {
		return ColorYellow + s + ColorReset
	})
	resYaml = numValuePattern.ReplaceAllStringFunc(resYaml, func(s string) string {
		return ColorGreen + s + ColorReset
	})

	return resYaml
}

func PrintColorizedYAML(program interface{}) {
	yamlData, err := yaml.Marshal(&program)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	yamlStr := string(yamlData)
	resYaml := colorizeYAML(yamlStr)

	fmt.Printf("---\n%s\n", resYaml)
}
