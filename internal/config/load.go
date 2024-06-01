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
	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("getting current directory: %w", err)
	}
	// Load current directory config
	currentDirConfigPath := filepath.Join(pwd, ".cubx", configFileName)
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

	preparedConfig, err := configPreprocessing(finalConfig, loadedConfigs)
	if err != nil {
		return nil, nil, err
	}
	return preparedConfig, loadedConfigs, nil
}

func findDockerfileInDirectory(dir, dockerfileName string) bool {
	// Building a Dockerfile path
	dockerfilePath := filepath.Join(dir, dockerfileName)
	// Checking for file availability
	if _, err := os.Stat(dockerfilePath); err == nil {
		return true
	}
	return false
}

func configPreprocessing(finalConfig *ProgramConfig, loadedConfigs []string) (*ProgramConfig, error) {
	foundDockerfile := make(map[string]bool)

	// Reversing the paths to configuration files in reverse order
	for i := len(loadedConfigs) - 1; i >= 0; i-- {
		filePath := loadedConfigs[i]
		dir := filepath.Dir(filePath)

		// Searching all programs in the final configuration
		for index := range finalConfig.Programs {
			program := &finalConfig.Programs[index] // We get a link to the program
			if program.Dockerfile == "" {
				continue // Skip programs without a specified Dockerfile path
			}

			// Search for Dockerfile in the current directory
			if findDockerfileInDirectory(dir, program.Dockerfile) {
				fmt.Printf("Dockerfile '%s' found in directory: %s\n", program.Dockerfile, dir)
				program.Dockerfile = filepath.Join(dir, program.Dockerfile) // Changing the path in the original slice
				foundDockerfile[program.Dockerfile] = true
			}
		}
		fmt.Println("Processed file and directory:", filePath, dir)
	}

	// Checking for the presence of a Dockerfile for each program that requires it
	for _, program := range finalConfig.Programs {
		if program.Dockerfile != "" && !foundDockerfile[program.Dockerfile] {
			return nil, fmt.Errorf("dockerfile not found for program: %s", program.Dockerfile)
		}
	}

	return finalConfig, nil
}
