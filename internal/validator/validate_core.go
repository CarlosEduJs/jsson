package validator

import "fmt"

// validateValue recursively validates a value against a schema.
func (v *Validator) validateValue(value any, schema *Schema, path string, result *ValidationResult) {
	if schema == nil {
		return
	}

	// Check const
	if schema.Const != nil {
		if !deepEqual(value, schema.Const) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Value must be equal to constant %v", schema.Const),
				Value:   value,
			})

			return
		}
	}

	// Check enum
	if len(schema.Enum) > 0 {
		found := false

		for _, enumVal := range schema.Enum {
			if deepEqual(value, enumVal) {
				found = true

				break
			}
		}

		if !found {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Value must be one of: %v", schema.Enum),
				Value:   value,
			})

			return
		}
	}

	// Check oneOf
	if len(schema.OneOf) > 0 {
		validCount := 0

		for _, subSchema := range schema.OneOf {
			subResult := &ValidationResult{Valid: true, Errors: []ValidationError{}}
			v.validateValue(value, subSchema, path, subResult)

			if subResult.Valid {
				validCount++
			}
		}

		if validCount != 1 {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: "Value must match exactly one of the oneOf schemas",
				Value:   value,
			})

			return
		}
	}

	// Check anyOf
	if len(schema.AnyOf) > 0 {
		validAny := false

		for _, subSchema := range schema.AnyOf {
			subResult := &ValidationResult{Valid: true, Errors: []ValidationError{}}
			v.validateValue(value, subSchema, path, subResult)

			if subResult.Valid {
				validAny = true

				break
			}
		}

		if !validAny {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: "Value must match at least one of the anyOf schemas",
				Value:   value,
			})

			return
		}
	}

	// Check allOf
	if len(schema.AllOf) > 0 {
		for _, subSchema := range schema.AllOf {
			v.validateValue(value, subSchema, path, result)
		}
	}

	// Check not
	if schema.Not != nil {
		subResult := &ValidationResult{Valid: true, Errors: []ValidationError{}}
		v.validateValue(value, schema.Not, path, subResult)

		if subResult.Valid {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: "Value must not match the 'not' schema",
				Value:   value,
			})

			return
		}
	}

	// Check type
	if schema.Type != "" {
		v.validateType(value, schema, path, result)
	}
}
