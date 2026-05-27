package transpiler

import (
	"encoding/json"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"strings"
	"testing"
)

// Test boolean literals transpilation.
func TestBooleanLiteralsTranspilation(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]any
	}{
		{
			input: `test { enabled = yes }`,
			expected: map[string]any{
				"test": map[string]any{"enabled": true},
			},
		},
		{
			input: `test { disabled = no }`,
			expected: map[string]any{
				"test": map[string]any{"disabled": false},
			},
		},
		{
			input: `test { active = on }`,
			expected: map[string]any{
				"test": map[string]any{"active": true},
			},
		},
		{
			input: `test { inactive = off }`,
			expected: map[string]any{
				"test": map[string]any{"inactive": false},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Fatalf("parser errors for '%s': %v", tt.input, p.Errors())
		}

		tr := New(program, "", "keep", "")

		output, err := tr.Transpile()
		if err != nil {
			t.Fatalf("transpiler error for '%s': %v", tt.input, err)
		}

		var result map[string]any
		if err := json.Unmarshal(output, &result); err != nil {
			t.Fatalf("json unmarshal error for '%s': %v", tt.input, err)
		}

		testObj, _ := result["test"].(map[string]any)
		expectedObj, _ := tt.expected["test"].(map[string]any)

		for key, expectedVal := range expectedObj {
			actualVal, exists := testObj[key]
			if !exists {
				t.Errorf("key '%s' not found in output for '%s'", key, tt.input)

				continue
			}

			if actualVal != expectedVal {
				t.Errorf("value wrong for '%s'. key=%s expected=%v got=%v",
					tt.input, key, expectedVal, actualVal)
			}
		}
	}
}

// Test validator value generation.
func TestValidatorValueGeneration(t *testing.T) {
	tests := []struct {
		input        string
		key          string
		validateFunc func(any) bool
		description  string
	}{
		{
			input: `test { id = @uuid }`,
			key:   "id",
			validateFunc: func(v any) bool {
				s, ok := v.(string)
				if !ok {
					return false
				}
				// UUID format: 8-4-4-4-12 hex chars
				parts := strings.Split(s, "-")

				return len(parts) == 5 && len(parts[0]) == 8 && len(parts[1]) == 4
			},
			description: "UUID should have format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		},
		{
			input: `test { email = @email }`,
			key:   "email",
			validateFunc: func(v any) bool {
				s, ok := v.(string)

				return ok && strings.Contains(s, "@") && strings.Contains(s, ".")
			},
			description: "Email should contain @ and .",
		},
		{
			input: `test { website = @url }`,
			key:   "website",
			validateFunc: func(v any) bool {
				s, ok := v.(string)

				return ok && (strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://"))
			},
			description: "URL should start with http:// or https://",
		},
		{
			input: `test { ip = @ipv4 }`,
			key:   "ip",
			validateFunc: func(v any) bool {
				s, ok := v.(string)
				if !ok {
					return false
				}

				parts := strings.Split(s, ".")

				return len(parts) == 4
			},
			description: "IPv4 should have 4 octets",
		},
		{
			input: `test { created = @date }`,
			key:   "created",
			validateFunc: func(v any) bool {
				s, ok := v.(string)
				if !ok {
					return false
				}
				// Date format: YYYY-MM-DD
				return len(s) == 10 && s[4] == '-' && s[7] == '-'
			},
			description: "Date should have format YYYY-MM-DD",
		},
		{
			input: `test { timestamp = @datetime }`,
			key:   "timestamp",
			validateFunc: func(v any) bool {
				s, ok := v.(string)
				if !ok {
					return false
				}
				// DateTime should contain T for ISO8601
				return strings.Contains(s, "T") || strings.Contains(s, " ")
			},
			description: "DateTime should contain T or space separator",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Fatalf("parser errors for '%s': %v", tt.input, p.Errors())
		}

		tr := New(program, "", "keep", "")

		output, err := tr.Transpile()
		if err != nil {
			t.Fatalf("transpiler error for '%s': %v", tt.input, err)
		}

		var result map[string]any
		if err := json.Unmarshal(output, &result); err != nil {
			t.Fatalf("json unmarshal error for '%s': %v", tt.input, err)
		}

		testObj, ok := result["test"].(map[string]any)
		if !ok {
			t.Fatalf("result['test'] not a map for '%s'", tt.input)
		}

		value, exists := testObj[tt.key]
		if !exists {
			t.Errorf("key '%s' not found in output for '%s'", tt.key, tt.input)

			continue
		}

		if !tt.validateFunc(value) {
			t.Errorf("validation failed for '%s': %s. got value: %v",
				tt.input, tt.description, value)
		}
	}
}

// Test unique UUID generation.
func TestUniqueUUIDGeneration(t *testing.T) {
	input := `
users {
  user1 { id = @uuid }
  user2 { id = @uuid }
  user3 { id = @uuid }
}
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")

	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	users, _ := result["users"].(map[string]any)

	user1, _ := users["user1"].(map[string]any)
	user2, _ := users["user2"].(map[string]any)
	user3, _ := users["user3"].(map[string]any)
	uuid1, _ := user1["id"].(string)
	uuid2, _ := user2["id"].(string)
	uuid3, _ := user3["id"].(string)

	if uuid1 == uuid2 || uuid1 == uuid3 || uuid2 == uuid3 {
		t.Errorf("UUIDs should be unique. got: %s, %s, %s", uuid1, uuid2, uuid3)
	}
}

// Test keywords as property names.
func TestKeywordsAsPropertyNamesTranspilation(t *testing.T) {
	input := `
test {
  uuid = "custom-uuid"
  email = "custom-email"
  ipv4 = "custom-ip"
  date = "custom-date"
}
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")

	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	testObj, _ := result["test"].(map[string]any)

	expectedKeys := map[string]string{
		"uuid":  "custom-uuid",
		"email": "custom-email",
		"ipv4":  "custom-ip",
		"date":  "custom-date",
	}

	for key, expectedVal := range expectedKeys {
		actualVal, exists := testObj[key]
		if !exists {
			t.Errorf("key '%s' not found in output", key)

			continue
		}

		if actualVal != expectedVal {
			t.Errorf("value wrong for key '%s'. expected=%s got=%v",
				key, expectedVal, actualVal)
		}
	}
}

// Test mixed features (boolean literals + validators).
func TestMixedFeaturesV006(t *testing.T) {
	input := `
app {
  name = "TestApp"
  enabled = yes
  debug = on
  id = @uuid
  admin_email = @email
  active = true
  inactive = off
}
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")

	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	app, _ := result["app"].(map[string]any)

	// Check string
	if app["name"] != "TestApp" {
		t.Errorf("name wrong. expected=TestApp got=%v", app["name"])
	}

	// Check boolean literals
	if app["enabled"] != true {
		t.Errorf("enabled should be true, got=%v", app["enabled"])
	}

	if app["debug"] != true {
		t.Errorf("debug should be true, got=%v", app["debug"])
	}

	if app["active"] != true {
		t.Errorf("active should be true, got=%v", app["active"])
	}

	if app["inactive"] != false {
		t.Errorf("inactive should be false, got=%v", app["inactive"])
	}

	// Check validators generated values
	idVal, ok := app["id"].(string)
	if !ok || !strings.Contains(idVal, "-") {
		t.Errorf("id should be a UUID string, got=%v", app["id"])
	}

	emailVal, ok := app["admin_email"].(string)
	if !ok || !strings.Contains(emailVal, "@") {
		t.Errorf("admin_email should be an email string, got=%v", app["admin_email"])
	}
}
