package config

import (
	"encoding/json"

	"dario.cat/mergo"
)

// mergeSettings merges two Settings objects with the values from the override having priority
// and concatenates IgnorePaths slices without duplicates.
func mergeSettings(base, override Settings) Settings {
	// Perform deep cloning of the override settings
	merged := Settings{}
	data, _ := json.Marshal(override)
	_ = json.Unmarshal(data, &merged)
	_ = mergo.Merge(&merged, base, mergo.WithOverride)

	// Merge IgnorePaths without duplicates
	ignorePathMap := make(map[string]bool)
	for _, path := range base.IgnorePaths {
		ignorePathMap[path] = true
	}
	for _, path := range override.IgnorePaths {
		if !ignorePathMap[path] {
			merged.IgnorePaths = append(merged.IgnorePaths, path)
			ignorePathMap[path] = true
		}
	}
	return merged
}

// semanticMerge updates the given config by applying inherited settings
func semanticMerge(config *ProgramConfig) {
	// Merge global settings into each program's settings
	for i, program := range config.Programs {
		// Step 1: Merge global settings into program settings
		program.Settings = mergeSettings(config.Settings, program.Settings)

		// Merge program settings into each hook's settings
		for j, hook := range program.Hooks {
			// Step 2: Merge program settings into hook settings
			hook.Settings = mergeSettings(program.Settings, hook.Settings)
			program.Hooks[j] = hook
		}

		config.Programs[i] = program
	}
}
