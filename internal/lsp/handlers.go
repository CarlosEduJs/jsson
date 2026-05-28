package lsp

import (
	"encoding/json"

	"jsson/internal/lexer"
	"jsson/internal/parser"
)

func (s *Server) handleInitialize(id, _ json.RawMessage) error {
	response := ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncOptions{
					OpenClose: true,
					Change:    1,
				},
				CompletionProvider: &CompletionOptions{
					TriggerCharacters: []string{".", ":", "="},
				},
				HoverProvider: true,
				SemanticTokensProvider: &SemanticTokensOptions{
					Legend: SemanticTokensLegend{
						TokenTypes: []string{
							"namespace", "type", "class", "enum", "interface",
							"struct", "typeParameter", "parameter", "variable", "property",
							"enumMember", "event", "function", "method", "macro",
							"keyword", "modifier", "comment", "string", "number",
							"regexp", "operator",
						},
						TokenModifiers: []string{},
					},
					Full: true,
				},
			},
			ServerInfo: ServerInfo{
				Name:    "jsson-lsp",
				Version: "0.0.6",
			},
		},
	}

	return s.writeMessage(response)
}

func (s *Server) handleDidOpen(raw json.RawMessage) error {
	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return err
	}

	s.mu.Lock()
	s.documents[params.TextDocument.URI] = &Document{
		URI:     params.TextDocument.URI,
		Content: params.TextDocument.Text,
		Version: params.TextDocument.Version,
	}
	s.mu.Unlock()

	return s.publishDiagnostics(params.TextDocument.URI, params.TextDocument.Text)
}

func (s *Server) handleDidChange(raw json.RawMessage) error {
	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return err
	}

	if len(params.ContentChanges) == 0 {
		return nil
	}

	s.mu.Lock()
	s.documents[params.TextDocument.URI] = &Document{
		URI:     params.TextDocument.URI,
		Content: params.ContentChanges[0].Text,
		Version: params.TextDocument.Version,
	}
	s.mu.Unlock()

	return s.publishDiagnostics(params.TextDocument.URI, params.ContentChanges[0].Text)
}

func (s *Server) handleDidClose(raw json.RawMessage) error {
	var params DidCloseTextDocumentParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return err
	}

	s.mu.Lock()
	delete(s.documents, params.TextDocument.URI)
	s.mu.Unlock()

	return nil
}

func (s *Server) handleShutdown(id json.RawMessage) error {
	response := ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  nil,
	}

	return s.writeMessage(response)
}

func (s *Server) parseDocument(content string) []string {
	l := lexer.New(content)
	p := parser.New(l)

	_ = p.ParseProgram()

	return p.Errors()
}
