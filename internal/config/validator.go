package config

import (
	"cubx/internal/platform"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validatePlatform(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// allow empty value
	if value == "" {
		return true
	}
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return false
	}
	return platform.IsValidOsArch(parts[0], parts[1])
}

// setDefaults sets default values for fields that are not set.
func setDefaults(config *ProgramConfig) {
	for i := range config.Programs {
		if config.Programs[i].Serializer == "" {
			config.Programs[i].Serializer = "default"
		}
		if config.Programs[i].Tag == "" {
			config.Programs[i].Tag = "latest"
		}
	}
}

func getValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("platform", validatePlatform)
	return validate
}

func validateProgramConfig(config *ProgramConfig) error {
	// Set default values
	setDefaults(config)

	validate := getValidator()

	err := validate.Struct(config)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, err := range validationErrors {
				errorMessages = append(errorMessages, fmt.Sprintf("Field validation error on '%s': '%v' is not a valid value", err.Field(), err.Value()))
			}
			return fmt.Errorf("validation errors: %s", errorMessages)
		}
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}
