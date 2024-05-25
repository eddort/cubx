package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func validateSerializer(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, v := range validSerializers {
		if value == v {
			return true
		}
	}
	return false
}

func getValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("serializer", validateSerializer)

	return validate
}

func validateProgramConfig(config *ProgramConfig) error {
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
