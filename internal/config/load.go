package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func LoadConfig() (*CommandConfig, []string, error) {
	viper.SetConfigType("yaml")
	configFile := "config"

	// Add paths for configuration search
	viper.AddConfigPath(".cubx") // current directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting home directory: %w", err)
	}
	viper.AddConfigPath(filepath.Join(home, ".cubx")) // user's home directory

	var loadedConfigs []string

	// Set the configuration file name
	viper.SetConfigName(configFile)

	// Attempt to read configuration from current directory
	if err := viper.ReadInConfig(); err == nil {
		loadedConfigs = append(loadedConfigs, viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return nil, nil, err
	}

	// Attempt to read configuration from home directory
	if err := viper.MergeInConfig(); err == nil {
		loadedConfigs = append(loadedConfigs, viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return nil, nil, err
	}

	var config CommandConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Set default value for Handler if not specified
	for i, cmd := range config.Commands {
		if cmd.Handler == "" {
			config.Commands[i].Handler = "default"
		}
	}

	// Initialize the validator
	validate := validator.New()

	// Validate the configuration structure
	if err := validate.Struct(&config); err != nil {
		return nil, nil, fmt.Errorf("validation error: %w", err)
	}

	return &config, loadedConfigs, nil
}
