package validator

import "fmt"

// validateType validates the type of a value.
func (v *Validator) validateType(value any, schema *Schema, path string, result *ValidationResult) {
	actualType := getType(value)

	if value == nil {
		if schema.Type != "null" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Expected type '%s', got 'null'", schema.Type),
				Value:   value,
			})
		}

		return
	}

	validType := false

	switch schema.Type {
	case typeString:
		validType = actualType == typeString
	case typeNumber:
		validType = actualType == typeNumber || actualType == typeInteger
	case typeInteger:
		validType = actualType == typeInteger || (actualType == typeNumber && isInteger(value))
	case "boolean":
		validType = actualType == "boolean"
	case typeArray:
		validType = actualType == typeArray
	case typeObject:
		validType = actualType == typeObject
	case typeNull:
		validType = value == nil
	}

	if !validType {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Expected type '%s', got '%s'", schema.Type, actualType),
			Value:   value,
		})

		return
	}

	switch schema.Type {
	case typeString:
		strValue, ok := value.(string)
		if !ok {
			return
		}

		v.validateString(strValue, schema, path, result)
	case typeNumber, typeInteger:
		v.validateNumber(value, schema, path, result)
	case typeArray:
		v.validateArray(value, schema, path, result)
	case typeObject:
		v.validateObject(value, schema, path, result)
	}
}
