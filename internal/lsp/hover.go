package lsp

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (s *Server) handleHover(id, raw json.RawMessage) error {
	var params HoverParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return s.sendError(id, "Invalid params")
	}

	doc, ok := s.getDocument(params.TextDocument.URI)
	if !ok {
		return s.sendError(id, "Document not found")
	}

	hoverInfo := s.getHoverInfo(doc.Content, params.Position.Line, params.Position.Character)

	response := ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  hoverInfo,
	}

	return s.writeMessage(response)
}

func (s *Server) getHoverInfo(content string, line, character int) *Hover {
	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return nil
	}

	currentLine := lines[line]
	if character >= len(currentLine) {
		return nil
	}

	word := s.getWordAtPosition(currentLine, character)
	if word == "" {
		return nil
	}

	paramInfo := s.getParameterInfo(currentLine, word, character)
	if paramInfo != "" {
		return &Hover{
			Contents: MarkupContent{
				Kind:  "markdown",
				Value: paramInfo,
			},
		}
	}

	varInfo := s.getVariableInfo(content, word, line)
	if varInfo != "" {
		return &Hover{
			Contents: MarkupContent{
				Kind:  "markdown",
				Value: varInfo,
			},
		}
	}

	objInfo := s.getObjectInfo(content, word, line)
	if objInfo != "" {
		return &Hover{
			Contents: MarkupContent{
				Kind:  "markdown",
				Value: objInfo,
			},
		}
	}

	doc := s.getDocumentation(word)
	if doc == "" {
		return nil
	}

	return &Hover{
		Contents: MarkupContent{
			Kind:  "markdown",
			Value: doc,
		},
	}
}

func (s *Server) getParameterInfo(line, paramName string, character int) string {
	beforeCursor := line[:character]

	if strings.Contains(beforeCursor, "map (") || strings.Contains(beforeCursor, "zip (") {
		startParen := strings.LastIndex(beforeCursor, "(")
		if startParen != -1 {
			afterParen := line[startParen:]

			endParen := strings.Index(afterParen, ")")
			if endParen != -1 {
				paramsSection := afterParen[1:endParen]
				if strings.Contains(paramsSection, paramName) {
					if strings.Contains(beforeCursor, "map (") {
						return fmt.Sprintf("# Parameter: `%s`\n\n**Map parameter** - represents each item in the mapped collection.\n\n**Example:**\n```jsson\nitems = (1..5 map (%s) = %s * 2)\n```",
							paramName, paramName, paramName)
					} else if strings.Contains(beforeCursor, "zip (") {
						return fmt.Sprintf("# Parameter: `%s`\n\n**Zip parameter** - represents corresponding items from parallel ranges.\n\n**Example:**\n```jsson\npairs [\n  template { a, b }\n  1..3, 10..12\n]\n```",
							paramName)
					}
				}
			}
		}
	}

	return ""
}

func (s *Server) getVariableInfo(content, varName string, _ int) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.Contains(line, varName+" :=") || strings.Contains(line, varName+":=") {
			parts := strings.Split(line, ":=")
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])

				return fmt.Sprintf("# Variable: `%s`\n\n**Declared at line %d**\n\n```jsson\n%s\n```\n\n**Value:** `%s`",
					varName, i+1, strings.TrimSpace(line), value)
			}
		}
	}

	return ""
}

func (s *Server) getObjectInfo(content, name string, _ int) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, name+" {") || strings.HasPrefix(trimmed, name+"{") {
			return fmt.Sprintf("# Object: `%s`\n\n**Defined at line %d**\n\n```jsson\n%s\n```",
				name, i+1, strings.TrimSpace(line))
		}

		if strings.HasPrefix(trimmed, name+" [") || strings.HasPrefix(trimmed, name+"[") {
			return fmt.Sprintf("# Array: `%s`\n\n**Defined at line %d**\n\n```jsson\n%s\n```",
				name, i+1, strings.TrimSpace(line))
		}

		if strings.HasPrefix(trimmed, name+" =") || strings.HasPrefix(trimmed, name+"=") {
			return fmt.Sprintf("# Property: `%s`\n\n**Defined at line %d**\n\n```jsson\n%s\n```",
				name, i+1, strings.TrimSpace(line))
		}
	}

	return ""
}

func (s *Server) getWordAtPosition(line string, character int) string {
	if character >= len(line) {
		return ""
	}

	start := character
	for start > 0 && isIdentifierChar(line[start-1]) {
		start--
	}

	end := character
	for end < len(line) && isIdentifierChar(line[end]) {
		end++
	}

	if start >= end {
		return ""
	}

	return line[start:end]
}

func isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
