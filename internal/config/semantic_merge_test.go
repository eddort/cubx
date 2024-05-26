package config

import (
	"reflect"
	"sort"
	"testing"

	"github.com/go-playground/validator/v10"
)

func sortStrings(s []string) []string {
	sorted := make([]string, len(s))
	copy(sorted, s)
	sort.Strings(sorted)
	return sorted
}

func TestSemanticMerge(t *testing.T) {
	config := &ProgramConfig{
		Programs: []Program{
			{
				Name:     "program1",
				Image:    "image1",
				Command:  "command1",
				Settings: Settings{Net: "bridge", IgnorePaths: []string{"/program/path"}},
				Hooks: []Hook{
					{
						Command:  "hook1",
						Settings: Settings{IgnorePaths: []string{"/hook/path"}},
					},
				},
			},
		},
		Settings: Settings{
			Net:         "bridge",
			IgnorePaths: []string{"/global/path"},
		},
	}

	// Perform the semantic merge
	semanticMerge(config)

	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	// Expected settings after merge
	expectedProgramSettings := Settings{Net: "bridge", IgnorePaths: []string{"/program/path", "/global/path"}}
	expectedHookSettings := Settings{Net: "bridge", IgnorePaths: []string{"/hook/path", "/program/path", "/global/path"}}
	expectedGlobalSettings := Settings{Net: "bridge", IgnorePaths: []string{"/global/path"}}

	// Sort IgnorePaths for comparison
	configProgramSettings := config.Programs[0].Settings
	configProgramSettings.IgnorePaths = sortStrings(configProgramSettings.IgnorePaths)
	expectedProgramSettings.IgnorePaths = sortStrings(expectedProgramSettings.IgnorePaths)

	configHookSettings := config.Programs[0].Hooks[0].Settings
	configHookSettings.IgnorePaths = sortStrings(configHookSettings.IgnorePaths)
	expectedHookSettings.IgnorePaths = sortStrings(expectedHookSettings.IgnorePaths)

	// Sort IgnorePaths for comparison
	configGlobalSettings := config.Settings
	configGlobalSettings.IgnorePaths = sortStrings(configGlobalSettings.IgnorePaths)
	expectedGlobalSettings.IgnorePaths = sortStrings(expectedGlobalSettings.IgnorePaths)

	// Check program settings after merge
	if !reflect.DeepEqual(configProgramSettings, expectedProgramSettings) {
		t.Errorf("expected program settings to be %+v, but got %+v", expectedProgramSettings, configProgramSettings)
	}

	// Check hook settings after merge
	if !reflect.DeepEqual(configHookSettings, expectedHookSettings) {
		t.Errorf("expected hook settings to be %+v, but got %+v", expectedHookSettings, configHookSettings)
	}

	// Check global settings after merge
	if !reflect.DeepEqual(configGlobalSettings, expectedGlobalSettings) {
		t.Errorf("expected global settings to be %+v, but got %+v", expectedGlobalSettings, configGlobalSettings)
	}
}
