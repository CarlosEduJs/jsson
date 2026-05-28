package validator

import (
	"fmt"
	"regexp"
)

// validateString validates a string value.
func (v *Validator) validateString(value string, schema *Schema, path string, result *ValidationResult) {
	if schema.MinLength != nil && len(value) < *schema.MinLength {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("String length must be at least %d", *schema.MinLength),
			Value:   value,
		})
	}

	if schema.MaxLength != nil && len(value) > *schema.MaxLength {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("String length must be at most %d", *schema.MaxLength),
			Value:   value,
		})
	}

	if schema.Pattern != "" {
		matched, err := regexp.MatchString(schema.Pattern, value)
		if err != nil || !matched {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("String does not match pattern '%s'", schema.Pattern),
				Value:   value,
			})
		}
	}

	if schema.Format != "" {
		v.validateJSSonFormat(value, schema.Format, path, result)
	}

	if schema.JSSonFormat != "" {
		v.validateJSSonFormat(value, schema.JSSonFormat, path, result)
	}
}

// validateNumber validates a numeric value.
func (v *Validator) validateNumber(value any, schema *Schema, path string, result *ValidationResult) {
	num := toFloat64(value)

	if schema.Minimum != nil && num < *schema.Minimum {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Number must be at least %v", *schema.Minimum),
			Value:   value,
		})
	}

	if schema.Maximum != nil && num > *schema.Maximum {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Number must be at most %v", *schema.Maximum),
			Value:   value,
		})
	}
}
