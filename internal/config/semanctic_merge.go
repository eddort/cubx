package config

import (
	"encoding/json"

	"dario.cat/mergo"
)

// mergeSettings merges two Settings objects with the values from the override having priority
// and concatenates IgnorePaths slices without duplicates.
func mergeSettings(base, override Settings) (Settings, error) {
	// Perform deep cloning of the override settings
	merged := Settings{}
	data, err := json.Marshal(override)
	if err != nil {
		return base, err
	}
	err = json.Unmarshal(data, &merged)
	if err != nil {
		return base, err
	}
	err = mergo.Merge(&merged, base, mergo.WithOverride)
	if err != nil {
		return base, err
	}

	// Merge IgnorePaths without duplicates

	ignorePathMap := make(map[string]bool)
	for _, path := range merged.IgnorePaths {
		ignorePathMap[path] = true
	}
	for _, path := range base.IgnorePaths {
		if !ignorePathMap[path] {
			merged.IgnorePaths = append(merged.IgnorePaths, path)
			ignorePathMap[path] = true
		}
	}
	for _, path := range override.IgnorePaths {
		if !ignorePathMap[path] {
			merged.IgnorePaths = append(merged.IgnorePaths, path)
			ignorePathMap[path] = true
		}
	}
	return merged, nil
}

// semanticMerge updates the given config by applying inherited settings
func semanticMerge(config *ProgramConfig) error {
	// Merge global settings into each program's settings
	for i, program := range config.Programs {
		// Step 1: Merge global settings into program settings
		mergedSettings, err := mergeSettings(config.Settings, program.Settings)
		if err != nil {
			return err
		}
		program.Settings = mergedSettings

		// Merge program settings into each hook's settings
		for j, hook := range program.Hooks {
			// Step 2: Merge program settings into hook settings
			mergedHookSettings, err := mergeSettings(program.Settings, hook.Settings)
			if err != nil {
				return err
			}
			hook.Settings = mergedHookSettings
			program.Hooks[j] = hook
		}

		config.Programs[i] = program
	}
	return nil
}
