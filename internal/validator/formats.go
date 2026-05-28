package validator

import (
	"regexp"
	"strings"
)

const ipv4Pattern = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

func matchString(pattern, value string) bool {
	matched, err := regexp.MatchString(pattern, value)

	return err == nil && matched
}

// validateJSSonFormat validates JSSON-specific format validators.
func (v *Validator) validateJSSonFormat(value, format, path string, result *ValidationResult) {
	var valid bool
	var message string

	switch format {
	case "email":
		valid = matchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, value)
		message = "Invalid email format"

	case "uri", "url":
		valid = matchString(`^https?://[^\s/$.?#].[^\s]*$`, value)
		message = "Invalid URL format"

	case "uuid":
		valid = matchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, value)
		message = "Invalid UUID format"

	case "date":
		valid = matchString(`^\d{4}-\d{2}-\d{2}$`, value) ||
			matchString(`^\d{2}/\d{2}/\d{4}$`, value) ||
			matchString(`^\d{2}-\d{2}-\d{4}$`, value)
		message = "Invalid date format"

	case "date-time", "datetime":
		valid = matchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})?$`, value) ||
			matchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, value) ||
			matchString(`^\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}$`, value)
		message = "Invalid datetime format"

	case "time":
		valid = matchString(`^([01]?[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?$`, value)
		message = "Invalid time format"

	case "ipv4":
		valid = matchString(ipv4Pattern, value)
		message = "Invalid IPv4 format"

	case "ipv6":
		valid = matchString(`^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`, value)
		message = "Invalid IPv6 format"

	case "hostname":
		valid = matchString(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`, value)
		message = "Invalid hostname format"

	case "file-path", "filepath":
		valid = value != "" && !strings.Contains(value, "\x00")
		message = "Invalid file path format"

	case "semver":
		valid = matchString(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-((0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$`, value)
		message = "Invalid semantic version format"

	case "hex-color", "hexcolor":
		valid = matchString(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`, value)
		message = "Invalid hex color format"

	case "rgb-color", "rgbcolor":
		valid = matchString(`^rgb\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)$`, value)
		message = "Invalid RGB color format"

	case "port":
		valid = matchString(`^([1-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`, value)
		message = "Invalid port number (must be 1-65535)"

	case "host":
		valid = matchString(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`, value) || matchString(ipv4Pattern, value)
		message = "Invalid host format"

	case "env-var", "envvar":
		valid = matchString(`^[A-Z_][A-Z0-9_]*$`, value)
		message = "Invalid environment variable name format"

	case "template-var", "templatevar":
		valid = matchString(`^\{\{[a-zA-Z_][a-zA-Z0-9_]*\}\}$`, value)
		message = "Invalid template variable format"

	case "json-pointer":
		valid = matchString(`^(/([^/~]|~[01])*)*$`, value)
		message = "Invalid JSON pointer format"

	case "regex":
		_, err := regexp.Compile(value)
		valid = err == nil
		message = "Invalid regex pattern"

	case "base64":
		valid = matchString(`^[A-Za-z0-9+/]*={0,2}$`, value) && len(value)%4 == 0
		message = "Invalid base64 format"

	case "phone":
		valid = matchString(`^[\+]?[(]?[0-9]{1,3}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,4}[-\s\.]?[0-9]{1,9}$`, value)
		message = "Invalid phone number format"

	case "credit-card", "creditcard":
		cleaned := strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", "")
		valid = matchString(`^\d{13,19}$`, cleaned)
		message = "Invalid credit card format"

	case "slug":
		valid = matchString(`^[a-z0-9]+(-[a-z0-9]+)*$`, value)
		message = "Invalid slug format"

	case "alpha":
		valid = matchString(`^[a-zA-Z]+$`, value)
		message = "Must contain only alphabetic characters"

	case "alphanumeric":
		valid = matchString(`^[a-zA-Z0-9]+$`, value)
		message = "Must contain only alphanumeric characters"

	case "macro-id":
		valid = matchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, value)
		message = "Invalid macro-id format (must start with letter or underscore)"

	default:
		return
	}

	if !valid {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: message,
			Value:   value,
		})
	}
}

// ValidateFormat checks if a value matches a specific format (exported helper).
func ValidateFormat(value, format string) bool {
	v := New()
	result := &ValidationResult{Valid: true, Errors: []ValidationError{}}
	v.validateJSSonFormat(value, format, "$", result)

	return result.Valid
}

// Common format validation helpers for direct use

// IsValidEmail checks if a string is a valid email.
func IsValidEmail(value string) bool {
	return ValidateFormat(value, "email")
}

// IsValidURL checks if a string is a valid URL.
func IsValidURL(value string) bool {
	return ValidateFormat(value, "url")
}

// IsValidUUID checks if a string is a valid UUID.
func IsValidUUID(value string) bool {
	return ValidateFormat(value, "uuid")
}

// IsValidIPv4 checks if a string is a valid IPv4 address.
func IsValidIPv4(value string) bool {
	return ValidateFormat(value, "ipv4")
}

// IsValidIPv6 checks if a string is a valid IPv6 address.
func IsValidIPv6(value string) bool {
	return ValidateFormat(value, "ipv6")
}

// IsValidDate checks if a string is a valid date.
func IsValidDate(value string) bool {
	return ValidateFormat(value, "date")
}

// IsValidDateTime checks if a string is a valid datetime.
func IsValidDateTime(value string) bool {
	return ValidateFormat(value, "datetime")
}

// IsValidSemVer checks if a string is a valid semantic version.
func IsValidSemVer(value string) bool {
	return ValidateFormat(value, "semver")
}

// IsValidHexColor checks if a string is a valid hex color.
func IsValidHexColor(value string) bool {
	return ValidateFormat(value, "hex-color")
}

// IsValidPort checks if a string is a valid port number.
func IsValidPort(value string) bool {
	return ValidateFormat(value, "port")
}
