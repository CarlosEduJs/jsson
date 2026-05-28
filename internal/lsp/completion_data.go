package lsp

// getKeywordCompletions returns keyword completions.
func (s *Server) getKeywordCompletions() []CompletionItem {
	keywords := []struct {
		label  string
		detail string
		doc    string
	}{
		{"template", "Template definition", "Define a template for structured data"},
		{"map", "Map transformation", "Transform data with a map function"},
		{"zip", "Zip ranges", "Combine multiple ranges into tuples"},
		{"include", "Include file", "Include another JSSON file"},
		{"step", "Range step", "Define step size for ranges"},
		{"@preset", "Preset definition", "Define a reusable preset configuration"},
		{"true", "Boolean true", "Boolean true value"},
		{"false", "Boolean false", "Boolean false value"},
		{"yes", "Boolean true", "Boolean true (alternative syntax)"},
		{"no", "Boolean false", "Boolean false (alternative syntax)"},
		{"on", "Boolean true", "Boolean true (alternative syntax)"},
		{"off", "Boolean false", "Boolean false (alternative syntax)"},
		{"null", "Null value", "Null value"},
		{"@uuid", "UUID validator", "Validates/generates UUID"},
		{"@email", "Email validator", "Validates/generates email"},
		{"@url", "URL validator", "Validates/generates URL"},
		{"@ipv4", "IPv4 validator", "Validates/generates IPv4"},
		{"@ipv6", "IPv6 validator", "Validates/generates IPv6"},
		{"@filepath", "File path validator", "Validates file path"},
		{"@date", "Date validator", "Validates/generates date"},
		{"@datetime", "DateTime validator", "Validates/generates datetime"},
		{"@regex", "Regex validator", "Validates with regex pattern"},
		{"@int", "Integer validator", "Generates random integer with min/max range"},
		{"@float", "Float validator", "Generates random float with min/max range"},
		{"@bool", "Boolean validator", "Generates random boolean value"},
		{"and", "Logical AND", "Logical AND operator"},
		{"or", "Logical OR", "Logical OR operator"},
		{"not", "Logical NOT", "Logical NOT operator"},
	}

	items := make([]CompletionItem, 0, len(keywords))
	for _, kw := range keywords {
		items = append(items, CompletionItem{
			Label:         kw.label,
			Kind:          CompletionItemKindKeyword,
			Detail:        kw.detail,
			Documentation: kw.doc,
		})
	}

	return items
}

// getSnippetCompletions returns snippet completions.
func (s *Server) getSnippetCompletions() []CompletionItem {
	return []CompletionItem{
		{
			Label:         "template",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Template array",
			Documentation: "Create an array with template definition",
			InsertText:    "[\n  template { $1 }\n\n  $2\n]",
		},
		{
			Label:         "map",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Map transformation",
			Documentation: "Transform data with map",
			InsertText:    "($1 map ($2) = $3)",
		},
		{
			Label:         "variable",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Variable declaration",
			Documentation: "Declare a variable with :=",
			InsertText:    "$1 := $2",
		},
		{
			Label:         "range",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Numeric range",
			Documentation: "Create a numeric range",
			InsertText:    "$1..$2",
		},
		{
			Label:         "ternary",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Ternary operator",
			Documentation: "Conditional expression",
			InsertText:    "$1 ? $2 : $3",
		},
		{
			Label:         "preset",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Preset definition",
			Documentation: "Define a reusable preset",
			InsertText:    "@preset \"$1\" {\n  $2\n}",
		},
		{
			Label:         "use preset",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Use preset",
			Documentation: "Apply a preset configuration",
			InsertText:    "@\"$1\" {\n  $2\n}",
		},
		{
			Label:         "object",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Object definition",
			Documentation: "Create an object",
			InsertText:    "$1 {\n  $2\n}",
		},
	}
}

// getPropertyCompletions returns property completions (for obj.prop).
func (s *Server) getPropertyCompletions() []CompletionItem {
	return []CompletionItem{
		{Label: "id", Kind: CompletionItemKindProperty},
		{Label: "name", Kind: CompletionItemKindProperty},
		{Label: "value", Kind: CompletionItemKindProperty},
		{Label: "type", Kind: CompletionItemKindProperty},
		{Label: "age", Kind: CompletionItemKindProperty},
		{Label: "email", Kind: CompletionItemKindProperty},
	}
}
