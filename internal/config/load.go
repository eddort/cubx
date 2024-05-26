package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

func cloneProgramConfig(config *ProgramConfig) (*ProgramConfig, error) {
	var clonedConfig ProgramConfig
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &clonedConfig); err != nil {
		return nil, err
	}
	return &clonedConfig, nil
}

func mergePrograms(baseConfig, overrideConfig *ProgramConfig) *[]Program {
	var allPrograms []Program
	programMap := make(map[string]Program)

	allPrograms = append(allPrograms, baseConfig.Programs...)
	allPrograms = append(allPrograms, overrideConfig.Programs...)

	for _, program := range allPrograms {
		programMap[program.Name] = program
	}

	var mergedPrograms []Program
	for _, program := range programMap {
		mergedPrograms = append(mergedPrograms, program)
	}

	sort.Slice(mergedPrograms, func(i, j int) bool {
		return mergedPrograms[i].Name < mergedPrograms[j].Name
	})

	return &mergedPrograms
}

func mergeConfigs(baseConfig, overrideConfig *ProgramConfig) (*ProgramConfig, error) {
	clonedConfig, err := cloneProgramConfig(baseConfig)
	if err != nil {
		return nil, err
	}

	clonedConfig.Programs = *mergePrograms(clonedConfig, overrideConfig)

	if err := mergo.Merge(&clonedConfig.Settings, &overrideConfig.Settings, mergo.WithOverride); err != nil {
		return nil, err
	}

	if err := semanticMerge(clonedConfig); err != nil {
		return nil, err
	}

	return clonedConfig, nil
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
		currentConfig, err = mergeConfigs(getProgramConfig(), currentConfig)
		if err != nil {
			return nil, nil, err
		}
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
	finalConfig, err := mergeConfigs(homeConfig, currentConfig)
	if err != nil {
		return nil, nil, err
	}
	return finalConfig, loadedConfigs, nil
}
