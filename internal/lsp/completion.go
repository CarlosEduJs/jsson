package lsp

import (
	"encoding/json"
	"strings"
)

type CompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind"`
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
}

const (
	CompletionItemKindText     = 1
	CompletionItemKindMethod   = 2
	CompletionItemKindFunction = 3
	CompletionItemKindKeyword  = 14
	CompletionItemKindVariable = 6
	CompletionItemKindProperty = 10
	CompletionItemKindSnippet  = 15
)

func (s *Server) handleCompletion(id, raw json.RawMessage) error {
	var params CompletionParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return s.sendError(id, "Invalid params")
	}

	doc, ok := s.getDocument(params.TextDocument.URI)
	if !ok {
		return s.sendError(id, "Document not found")
	}

	items := s.getCompletionItems(doc.Content, params.Position.Line, params.Position.Character)

	response := ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result: CompletionList{
			IsIncomplete: false,
			Items:        items,
		},
	}

	return s.writeMessage(response)
}

func (s *Server) getCompletionItems(content string, line, character int) []CompletionItem {
	items := []CompletionItem{}

	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return items
	}

	currentLine := lines[line]

	beforeCursor := ""
	if character <= len(currentLine) {
		beforeCursor = currentLine[:character]
	}

	if strings.HasSuffix(beforeCursor, "..") || strings.HasSuffix(beforeCursor, ".. ") {
		items = append(items, s.getRangeCompletions(content, line)...)

		return items
	}

	items = append(items, s.getKeywordCompletions()...)
	items = append(items, s.getSnippetCompletions()...)

	if strings.HasSuffix(strings.TrimSpace(beforeCursor), ".") {
		items = append(items, s.getPropertyCompletions()...)
	}

	items = append(items, s.getVariableCompletions(content)...)
	items = append(items, s.getScopeVariables(content, line)...)

	return items
}
