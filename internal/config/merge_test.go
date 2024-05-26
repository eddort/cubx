package config

import (
	"testing"
)

// Test function for merging configurations with override
func TestMergeConfigs_Override(t *testing.T) {
	baseConfig := &ProgramConfig{
		Programs: []Program{
			{
				Name:        "program1",
				Image:       "baseImage",
				Command:     "baseCommand",
				Serializer:  "baseSerializer",
				Description: "baseDescription",
				Category:    "baseCategory",
				Hooks: []Hook{
					{
						Command: "baseHookCommand",
						Settings: Settings{
							Net:         "host",
							IgnorePaths: []string{"/base/path"},
						},
					},
				},
				Settings: Settings{
					Net:         "host",
					IgnorePaths: []string{"/base/path"},
				},
			},
		},
		Settings: Settings{
			Net:         "host",
			IgnorePaths: []string{"/base/path"},
		},
	}

	overrideConfig := &ProgramConfig{
		Programs: []Program{
			{
				Name:        "program1",
				Image:       "overrideImage",
				Command:     "overrideCommand",
				Serializer:  "string",
				Description: "overrideDescription",
				Category:    "overrideCategory",
				Hooks: []Hook{
					{
						Command: "overrideHookCommand",
						Settings: Settings{
							Net:         "none",
							IgnorePaths: []string{"/override/path"},
						},
					},
				},
				Settings: Settings{
					Net:         "none",
					IgnorePaths: []string{"/override/path"},
				},
			},
		},
		Settings: Settings{
			Net:         "none",
			IgnorePaths: []string{"/override/path"},
		},
	}

	mergedConfig, err := mergeConfigs(baseConfig, overrideConfig)
	if err != nil {
		t.Fatalf("mergeConfigs failed: %v", err)
	}

	// Initialize the validator
	if err := validateProgramConfig(mergedConfig); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	if err != nil {
		t.Fatalf("mergeConfigs failed: %v", err)
	}

	// Validate the merged config
	if len(mergedConfig.Programs) != 1 {
		t.Fatalf("expected 1 program, got %d", len(mergedConfig.Programs))
	}

	program := mergedConfig.Programs[0]
	if program.Image != "overrideImage" {
		t.Errorf("expected Image to be overrideImage, got %s", program.Image)
	}
	if program.Command != "overrideCommand" {
		t.Errorf("expected Command to be overrideCommand, got %s", program.Command)
	}
	if program.Serializer != "string" {
		t.Errorf("expected Serializer to be string, got %s", program.Serializer)
	}
	if program.Description != "overrideDescription" {
		t.Errorf("expected Description to be overrideDescription, got %s", program.Description)
	}
	if program.Category != "overrideCategory" {
		t.Errorf("expected Category to be overrideCategory, got %s", program.Category)
	}
	if len(program.Hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(program.Hooks))
	}
	hook := program.Hooks[0]
	if hook.Command != "overrideHookCommand" {
		t.Errorf("expected Hook Command to be overrideHookCommand, got %s", hook.Command)
	}
	if hook.Settings.Net != "none" {
		t.Errorf("expected Hook Settings Net to be none, got %s", hook.Settings.Net)
	}
	if len(hook.Settings.IgnorePaths) != 1 || hook.Settings.IgnorePaths[0] != "/override/path" {
		t.Errorf("expected Hook Settings IgnorePaths to be [/override/path], got %v", hook.Settings.IgnorePaths)
	}
	if mergedConfig.Settings.Net != "none" {
		t.Errorf("Settings to be none, got %s", mergedConfig.Settings.Net)
	}
	if len(mergedConfig.Settings.IgnorePaths) != 1 || mergedConfig.Settings.IgnorePaths[0] != "/override/path" {
		t.Errorf("expected Settings IgnorePaths to be [/override/path], got %v", mergedConfig.Settings.IgnorePaths)
	}
}

func TestMergeConfigs_AddNewProgram(t *testing.T) {
	baseConfig := &ProgramConfig{
		Programs: []Program{
			{
				Name:        "baseProgram",
				Image:       "baseImage",
				Command:     "baseCommand",
				Serializer:  "string",
				Description: "baseDescription",
				Category:    "baseCategory",
				Hooks: []Hook{
					{
						Command: "baseHookCommand",
						Settings: Settings{
							Net:         "host",
							IgnorePaths: []string{"/base/path"},
						},
					},
				},
				Settings: Settings{
					Net:         "host",
					IgnorePaths: []string{"/base/path"},
				},
			},
		},
		Settings: Settings{
			Net:         "host",
			IgnorePaths: []string{"/base/path"},
		},
	}

	overrideConfig := &ProgramConfig{
		Programs: []Program{
			{
				Name:        "newProgram",
				Image:       "newImage",
				Command:     "newCommand",
				Serializer:  "default",
				Description: "newDescription",
				Category:    "newCategory",
				Hooks: []Hook{
					{
						Command: "newHookCommand",
						Settings: Settings{
							Net:         "none",
							IgnorePaths: []string{"/new/path"},
						},
					},
				},
				Settings: Settings{
					Net:         "none",
					IgnorePaths: []string{"/new/path"},
				},
			},
		},
		Settings: Settings{
			Net:         "",
			IgnorePaths: []string{"/override/path"},
		},
	}

	mergedConfig, err := mergeConfigs(baseConfig, overrideConfig)
	if err != nil {
		t.Fatalf("mergeConfigs failed: %v", err)
	}

	if mergedConfig.Settings.Net != "host" {
		t.Errorf("expected Hook Settings Net to be host, got %s", mergedConfig.Settings.Net)
	}

	// Initialize the validator
	if err := validateProgramConfig(mergedConfig); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	// Validate the merged config
	if len(mergedConfig.Programs) != 2 {
		t.Fatalf("expected 2 programs, got %d", len(mergedConfig.Programs))
	}

	program1 := mergedConfig.Programs[0]
	if program1.Name != "baseProgram" {
		t.Errorf("expected Name to be baseProgram, got %s", program1.Name)
	}

	program2 := mergedConfig.Programs[1]
	if program2.Name != "newProgram" {
		t.Errorf("expected Name to be newProgram, got %s", program2.Name)
	}
	if program2.Image != "newImage" {
		t.Errorf("expected Image to be newImage, got %s", program2.Image)
	}
	if program2.Command != "newCommand" {
		t.Errorf("expected Command to be newCommand, got %s", program2.Command)
	}
	if program2.Serializer != "default" {
		t.Errorf("expected Serializer to be default, got %s", program2.Serializer)
	}
	if program2.Description != "newDescription" {
		t.Errorf("expected Description to be newDescription, got %s", program2.Description)
	}
	if program2.Category != "newCategory" {
		t.Errorf("expected Category to be newCategory, got %s", program2.Category)
	}
	if len(program2.Hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(program2.Hooks))
	}
	hook := program2.Hooks[0]
	if hook.Command != "newHookCommand" {
		t.Errorf("expected Hook Command to be newHookCommand, got %s", hook.Command)
	}
	if hook.Settings.Net != "host" {
		t.Errorf("expected Hook Settings Net to be host, got %s", hook.Settings.Net)
	}
	if len(hook.Settings.IgnorePaths) != 2 || hook.Settings.IgnorePaths[0] != "/override/path" || hook.Settings.IgnorePaths[1] != "/new/path" {
		t.Errorf("expected Hook Settings IgnorePaths to be [/override/path /new/path], got %v", hook.Settings.IgnorePaths)
	}
	if mergedConfig.Settings.Net != "host" {
		t.Errorf("Settings to be host, got %s", mergedConfig.Settings.Net)
	}
	if len(mergedConfig.Settings.IgnorePaths) != 1 || mergedConfig.Settings.IgnorePaths[0] != "/override/path" {
		t.Errorf("expected Settings IgnorePaths to be [/override/path], got %v", mergedConfig.Settings.IgnorePaths)
	}
}
