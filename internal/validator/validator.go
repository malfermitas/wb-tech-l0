package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(v interface{}) error
}

type DefaultValidator struct {
	validate *validator.Validate
}

func NewValidator() *DefaultValidator {
	return &DefaultValidator{
		validate: validator.New(),
	}
}

func (v *DefaultValidator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(validationErrors, "; "))
	}
	return nil
}
