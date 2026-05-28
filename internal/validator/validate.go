package validator

import (
	"encoding/json"
	"fmt"

	yamlv3 "gopkg.in/yaml.v3"
)

// Validate validates data against a schema.
func (v *Validator) Validate(data []byte, schema *Schema, format string) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
		Format: format,
	}

	var (
		parsedData any
		err        error
	)

	switch format {
	case "json":
		err = json.Unmarshal(data, &parsedData)
	case "yaml":
		err = yamlv3.Unmarshal(data, &parsedData)
		if err == nil {
			parsedData = normalizeData(parsedData)
		}
	case "toml":
		parsedData = parseTOML(string(data))
		err = nil
	case "typescript", "ts":
		parsedData, err = parseTypeScript(string(data))
	default:
		err = json.Unmarshal(data, &parsedData)
		if err != nil {
			err = yamlv3.Unmarshal(data, &parsedData)
			if err == nil {
				parsedData = normalizeData(parsedData)
			}
		}
	}

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    "$",
			Message: fmt.Sprintf("Failed to parse %s: %s", format, err.Error()),
		})

		return result
	}

	v.validateValue(parsedData, schema, "$", result)

	return result
}

// ValidateJSON validates JSON data against a schema.
func (v *Validator) ValidateJSON(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "json")
}

// ValidateYAML validates YAML data against a schema.
func (v *Validator) ValidateYAML(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "yaml")
}

// ValidateTOML validates TOML data against a schema.
func (v *Validator) ValidateTOML(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "toml")
}

// ValidateTypeScript validates TypeScript data against a schema.
func (v *Validator) ValidateTypeScript(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "typescript")
}
