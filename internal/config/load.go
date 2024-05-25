package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func mergeConfigs(baseConfig, overrideConfig *ProgramConfig) *ProgramConfig {
	programSet := make(map[string]bool)
	var mergedPrograms []Program

	// Add programs from baseConfig
	for _, program := range baseConfig.Programs {
		if !programSet[program.Name] {
			mergedPrograms = append(mergedPrograms, program)
			programSet[program.Name] = true
		}
	}

	// Add/override programs from overrideConfig
	for _, program := range overrideConfig.Programs {
		if !programSet[program.Name] {
			mergedPrograms = append(mergedPrograms, program)
			programSet[program.Name] = true
		} else {
			// Replace existing program with the one from overrideConfig
			for i := range mergedPrograms {
				if mergedPrograms[i].Name == program.Name {
					mergedPrograms[i] = program
					break
				}
			}
		}
	}

	return &ProgramConfig{Programs: mergedPrograms}
}

func loadConfigFile(filePath string) (*ProgramConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ProgramConfig{}, nil
		}
		return nil, err
	}

	var config ProgramConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unable to decode file %s into struct: %w", filePath, err)
	}

	// Validate the configuration structure
	if err := validateProgramConfig(&config); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return &config, nil
}

func LoadConfig(withDefaults bool) (*ProgramConfig, []string, error) {
	configFileName := "config.yaml"
	var loadedConfigs []string

	// Load current directory config
	currentDirConfigPath := filepath.Join(".cubx", configFileName)
	currentConfig, err := loadConfigFile(currentDirConfigPath)
	if err != nil {
		return nil, nil, err
	}

	if withDefaults {
		currentConfig = mergeConfigs(currentConfig, getProgramConfig())
	}

	if _, err := os.Stat(currentDirConfigPath); err == nil {
		loadedConfigs = append(loadedConfigs, currentDirConfigPath)
	}

	// Load home directory config
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting home directory: %w", err)
	}
	homeConfigPath := filepath.Join(home, ".cubx", configFileName)
	homeConfig, err := loadConfigFile(homeConfigPath)
	if err != nil {
		return nil, nil, err
	}
	if _, err := os.Stat(homeConfigPath); err == nil && homeConfigPath != currentDirConfigPath {
		loadedConfigs = append(loadedConfigs, homeConfigPath)
	}

	// Merge the configurations
	finalConfig := mergeConfigs(currentConfig, homeConfig)

	return finalConfig, loadedConfigs, nil
}
