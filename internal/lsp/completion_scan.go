package lsp

import (
	"encoding/json"
	"strings"
)

func (s *Server) getScopeVariables(content string, currentLine int) []CompletionItem {
	items := []CompletionItem{}
	lines := strings.Split(content, "\n")

	if currentLine >= len(lines) {
		return items
	}

	params := make(map[string]bool)
	openParens := 0

	for i := currentLine; i >= 0; i-- {
		line := lines[i]

		for _, ch := range line {
			switch ch {
			case '(':
				openParens++
			case ')':
				openParens--
			}
		}

		if strings.Contains(line, " map (") || strings.Contains(line, " zip (") {
			startIdx := strings.Index(line, " map (")
			if startIdx == -1 {
				startIdx = strings.Index(line, " zip (")
			}

			if startIdx != -1 {
				rest := line[startIdx:]
				parenStart := strings.Index(rest, "(")

				parenEnd := strings.Index(rest, ")")
				if parenStart != -1 && parenEnd != -1 && parenEnd > parenStart {
					paramsStr := rest[parenStart+1 : parenEnd]
					for _, param := range strings.Split(paramsStr, ",") {
						param = strings.TrimSpace(param)
						if param != "" && isValidIdentifier(param) {
							params[param] = true
						}
					}
				}
			}
		}

		if openParens <= 0 && i < currentLine {
			break
		}
	}

	for param := range params {
		items = append(items, CompletionItem{
			Label:  param,
			Kind:   CompletionItemKindVariable,
			Detail: "Parameter (in scope)",
		})
	}

	return items
}

func (s *Server) getRangeCompletions(content string, _ int) []CompletionItem {
	items := []CompletionItem{}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":=") {
			parts := strings.Split(line, ":=")
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				if varName != "" && isValidIdentifier(varName) {
					isNumeric := false

					for _, ch := range value {
						if ch >= '0' && ch <= '9' {
							isNumeric = true

							break
						}
					}

					if isNumeric {
						items = append(items, CompletionItem{
							Label:  varName,
							Kind:   CompletionItemKindVariable,
							Detail: "Range end: " + value,
						})
					}
				}
			}
		}
	}

	items = append(items, CompletionItem{
		Label:      "10",
		Kind:       CompletionItemKindText,
		Detail:     "Range: 0..10",
		InsertText: "10",
	}, CompletionItem{
		Label:      "100",
		Kind:       CompletionItemKindText,
		Detail:     "Range: 0..100",
		InsertText: "100",
	})

	return items
}

func (s *Server) getVariableCompletions(content string) []CompletionItem {
	items := []CompletionItem{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if strings.Contains(line, ":=") {
			parts := strings.Split(line, ":=")
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				if varName != "" && isValidIdentifier(varName) {
					items = append(items, CompletionItem{
						Label:  varName,
						Kind:   CompletionItemKindVariable,
						Detail: "Variable",
					})
				}
			}
		}
	}

	return items
}

func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}

	if (s[0] < 'a' || s[0] > 'z') && (s[0] < 'A' || s[0] > 'Z') && s[0] != '_' {
		return false
	}

	for i := 1; i < len(s); i++ {
		c := s[i]
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}

	return true
}

func (s *Server) sendError(id json.RawMessage, message string) error {
	response := ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: &ErrorObject{
			Code:    -32602,
			Message: message,
		},
	}

	return s.writeMessage(response)
}
