package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func mergeConfigs(baseConfig, overrideConfig *CommandConfig) *CommandConfig {
	commandMap := make(map[string]Command)
	for _, cmd := range baseConfig.Commands {
		commandMap[cmd.Name] = cmd
	}
	for _, cmd := range overrideConfig.Commands {
		commandMap[cmd.Name] = cmd
	}
	mergedCommands := make([]Command, 0, len(commandMap))
	for _, cmd := range commandMap {
		mergedCommands = append(mergedCommands, cmd)
	}
	return &CommandConfig{Commands: mergedCommands}
}

func LoadConfig() (*CommandConfig, []string, error) {
	configFile := "config"
	var loadedConfigs []string

	// Create a new viper instance for the current directory config
	viperCurrent := viper.New()
	viperCurrent.SetConfigName(configFile)
	viperCurrent.SetConfigType("yaml")
	viperCurrent.AddConfigPath(".cubx")

	if err := viperCurrent.ReadInConfig(); err == nil {
		loadedConfigs = append(loadedConfigs, viperCurrent.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return nil, nil, err
	}

	// Create a new viper instance for the home directory config
	viperHome := viper.New()
	viperHome.SetConfigName(configFile)
	viperHome.SetConfigType("yaml")
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting home directory: %w", err)
	}
	viperHome.AddConfigPath(filepath.Join(home, ".cubx"))

	if err := viperHome.ReadInConfig(); err == nil {
		if viperHome.ConfigFileUsed() != "" && viperHome.ConfigFileUsed() != viperCurrent.ConfigFileUsed() {
			loadedConfigs = append(loadedConfigs, viperHome.ConfigFileUsed())
		}
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return nil, nil, err
	}

	var currentConfig CommandConfig
	if err := viperCurrent.Unmarshal(&currentConfig); err != nil {
		return nil, nil, fmt.Errorf("unable to decode current directory config into struct: %w", err)
	}

	var homeConfig CommandConfig
	if err := viperHome.Unmarshal(&homeConfig); err != nil {
		return nil, nil, fmt.Errorf("unable to decode home directory config into struct: %w", err)
	}

	// Merge the configurations
	finalConfig := mergeConfigs(&currentConfig, &homeConfig)

	// Set default value for Handler if not specified
	for i, cmd := range finalConfig.Commands {
		if cmd.Handler == "" {
			finalConfig.Commands[i].Handler = "default"
		}
	}

	// Initialize the validator
	validate := validator.New()

	// Validate the configuration structure
	if err := validate.Struct(finalConfig); err != nil {
		return nil, nil, fmt.Errorf("validation error: %w", err)
	}

	return finalConfig, loadedConfigs, nil
}
