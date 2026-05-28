package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (s *Server) Start(ctx context.Context) error {
	log.Println("JSSON Language Server starting...")

	defer log.Println("JSSON Language Server stopped")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

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

	const maxMessageSize = 100 * 1024 * 1024 // 100 MB

	if contentLength > maxMessageSize {
		return nil, fmt.Errorf("message too large: %d bytes exceeds maximum of %d", contentLength, maxMessageSize)
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

func (s *Server) getDocument(uri string) (*Document, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.documents[uri]

	return doc, ok
}
