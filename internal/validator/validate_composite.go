package validator

import (
	"fmt"
	"reflect"
)

// validateArray validates an array value.
func (v *Validator) validateArray(value any, schema *Schema, path string, result *ValidationResult) {
	var arr []any

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Slice {
		arr = make([]any, rv.Len())
		for i := range rv.Len() {
			arr[i] = rv.Index(i).Interface()
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: "Expected array type",
			Value:   value,
		})

		return
	}

	if schema.MinItems != nil && len(arr) < *schema.MinItems {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Array must have at least %d items", *schema.MinItems),
			Value:   value,
		})
	}

	if schema.MaxItems != nil && len(arr) > *schema.MaxItems {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Array must have at most %d items", *schema.MaxItems),
			Value:   value,
		})
	}

	if schema.UniqueItems {
		seen := make(map[string]bool)

		for i, item := range arr {
			key := fmt.Sprintf("%v", item)
			if seen[key] {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Path:    fmt.Sprintf("%s[%d]", path, i),
					Message: "Array items must be unique",
					Value:   item,
				})
			}

			seen[key] = true
		}
	}

	if schema.Items != nil {
		for i, item := range arr {
			itemPath := fmt.Sprintf("%s[%d]", path, i)
			v.validateValue(item, schema.Items, itemPath, result)
		}
	}
}

// validateObject validates an object value.
func (v *Validator) validateObject(value any, schema *Schema, path string, result *ValidationResult) {
	obj, ok := value.(map[string]any)
	if !ok {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: "Expected object type",
			Value:   value,
		})

		return
	}

	for _, req := range schema.Required {
		if _, exists := obj[req]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Missing required property '%s'", req),
			})
		}
	}

	if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
		for key := range obj {
			if _, defined := schema.Properties[key]; !defined {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Path:    fmt.Sprintf("%s.%s", path, key),
					Message: fmt.Sprintf("Additional property '%s' is not allowed", key),
					Value:   obj[key],
				})
			}
		}
	}

	for key, propSchema := range schema.Properties {
		if propValue, exists := obj[key]; exists {
			propPath := fmt.Sprintf("%s.%s", path, key)
			v.validateValue(propValue, propSchema, propPath, result)
		}
	}
}
