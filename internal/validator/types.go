package validator

// ValidationResult contains the result of a validation operation.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
	Format string            `json:"format,omitempty"`
}

// ValidationError represents a single validation error.
type ValidationError struct {
	Path       string `json:"path"`
	Message    string `json:"message"`
	Value      any    `json:"value,omitempty"`
	SchemaPath string `json:"schemaPath,omitempty"`
	Expected   string `json:"expected,omitempty"`
}

// Schema represents a JSON Schema for validation.
type Schema struct {
	Ref                  string             `json:"$ref,omitempty"                 yaml:"$ref,omitempty"`
	ID                   string             `json:"$id,omitempty"                  yaml:"$id,omitempty"`
	Type                 string             `json:"type,omitempty"                 yaml:"type,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"           yaml:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"             yaml:"required,omitempty"`
	Items                *Schema            `json:"items,omitempty"                yaml:"items,omitempty"`
	MinLength            *int               `json:"minLength,omitempty"            yaml:"minLength,omitempty"`
	MaxLength            *int               `json:"maxLength,omitempty"            yaml:"maxLength,omitempty"`
	Minimum              *float64           `json:"minimum,omitempty"              yaml:"minimum,omitempty"`
	Maximum              *float64           `json:"maximum,omitempty"              yaml:"maximum,omitempty"`
	ExclusiveMinimum     *float64           `json:"exclusiveMinimum,omitempty"     yaml:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum     *float64           `json:"exclusiveMaximum,omitempty"     yaml:"exclusiveMaximum,omitempty"`
	MultipleOf           *float64           `json:"multipleOf,omitempty"           yaml:"multipleOf,omitempty"`
	Pattern              string             `json:"pattern,omitempty"              yaml:"pattern,omitempty"`
	Enum                 []any              `json:"enum,omitempty"                 yaml:"enum,omitempty"`
	Format               string             `json:"format,omitempty"               yaml:"format,omitempty"`
	JSSonFormat          string             `json:"jssonFormat,omitempty"          yaml:"jssonFormat,omitempty"`
	AdditionalProperties *bool              `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	MinItems             *int               `json:"minItems,omitempty"             yaml:"minItems,omitempty"`
	MaxItems             *int               `json:"maxItems,omitempty"             yaml:"maxItems,omitempty"`
	UniqueItems          bool               `json:"uniqueItems,omitempty"          yaml:"uniqueItems,omitempty"`
	Const                any                `json:"const,omitempty"                yaml:"const,omitempty"`
	Default              any                `json:"default,omitempty"              yaml:"default,omitempty"`
	Description          string             `json:"description,omitempty"          yaml:"description,omitempty"`
	Title                string             `json:"title,omitempty"                yaml:"title,omitempty"`
	If                   *Schema            `json:"if,omitempty"                   yaml:"if,omitempty"`
	Then                 *Schema            `json:"then,omitempty"                 yaml:"then,omitempty"`
	Else                 *Schema            `json:"else,omitempty"                 yaml:"else,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"                yaml:"oneOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty"                yaml:"anyOf,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty"                yaml:"allOf,omitempty"`
	Not                  *Schema            `json:"not,omitempty"                  yaml:"not,omitempty"`
}

// Validator is the main validation engine.
type Validator struct {
	schemas     map[string]*Schema
	customRules map[string]ValidationRule
}

// ValidationRule is a custom validation function.
type ValidationRule func(value any, params map[string]any) bool
