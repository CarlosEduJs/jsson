package lsp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Server struct {
	reader    io.Reader
	writer    io.Writer
	mu        sync.Mutex
	documents map[string]*Document
}

type Document struct {
	URI     string
	Content string
	Version int
}

func NewServer(reader io.Reader, writer io.Writer) *Server {
	return &Server{
		reader:    reader,
		writer:    writer,
		documents: make(map[string]*Document),
	}
}

func (s *Server) Start() error {
	log.Println("JSSON Language Server starting...")

	for {
		msg, err := s.readMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			log.Printf("Error reading message: %v", err)

			return err
		}

		if err := s.handleMessage(msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

func (s *Server) readMessage() (*RequestMessage, error) {
	reader := bufio.NewReader(s.reader)

	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, errors.New("missing Content-Length header")
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil || contentLength == 0 {
		return nil, fmt.Errorf("invalid Content-Length: %s", contentLengthStr)
	}

	content := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, content); err != nil {
		return nil, err
	}

	var msg RequestMessage
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *Server) writeMessage(msg any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	if _, err := s.writer.Write([]byte(header)); err != nil {
		return err
	}

	if _, err := s.writer.Write(content); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleMessage(msg *RequestMessage) error {
	log.Printf("Received method: %s", msg.Method)

	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg.ID, msg.Params)
	case "initialized":
		return nil
	case "textDocument/didOpen":
		return s.handleDidOpen(msg.Params)
	case "textDocument/didChange":
		return s.handleDidChange(msg.Params)
	case "textDocument/didClose":
		return s.handleDidClose(msg.Params)
	case "textDocument/completion":
		return s.handleCompletion(msg.ID, msg.Params)
	case "textDocument/hover":
		return s.handleHover(msg.ID, msg.Params)
	case "textDocument/semanticTokens/full":
		return s.handleSemanticTokensFull(msg.ID, msg.Params)
	case "shutdown":
		return s.handleShutdown(msg.ID)
	case "exit":
		return io.EOF
	default:
		log.Printf("Unhandled method: %s", msg.Method)

		return nil
	}
}

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

// handleDidOpen handles textDocument/didOpen notification.
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

// handleDidChange handles textDocument/didChange notification.
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

// handleDidClose handles textDocument/didClose notification.
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

// handleShutdown handles the shutdown request.
func (s *Server) handleShutdown(id json.RawMessage) error {
	response := ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  nil,
	}

	return s.writeMessage(response)
}

// getDocument retrieves a document from the cache.
func (s *Server) getDocument(uri string) (*Document, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.documents[uri]

	return doc, ok
}

// parseDocument parses a JSSON document and returns any errors.
func (s *Server) parseDocument(content string) []string {
	l := lexer.New(content)
	p := parser.New(l)

	_ = p.ParseProgram()

	return p.Errors()
}
