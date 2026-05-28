package errors

import "fmt"

// LexError represents an error from the lexer/tokenizer.
type LexError struct {
	SourceFile string
	Line, Col  int
	Message    string
}

func (e *LexError) Error() string {
	if e.SourceFile != "" {
		return fmt.Sprintf("Lexer goblin: %s — %s",
			FormatContext(e.SourceFile, e.Line, e.Col), e.Message)
	}

	return fmt.Sprintf("Lexer goblin: %d:%d — %s", e.Line, e.Col, e.Message)
}

// ParseError represents an error from the parser.
type ParseError struct {
	SourceFile string
	Line, Col  int
	Message    string
}

func (e *ParseError) Error() string {
	if e.SourceFile != "" {
		return fmt.Sprintf("Syntax wizard: %s — %s",
			FormatContext(e.SourceFile, e.Line, e.Col), e.Message)
	}

	return fmt.Sprintf("Syntax wizard: %d:%d — %s", e.Line, e.Col, e.Message)
}

// TranspileError represents an error from the transpiler/evaluator.
type TranspileError struct {
	SourceFile string
	Line, Col  int
	Message    string
}

func (e *TranspileError) Error() string {
	if e.SourceFile != "" {
		return fmt.Sprintf("Transpile gremlin: %s — %s",
			FormatContext(e.SourceFile, e.Line, e.Col), e.Message)
	}

	return fmt.Sprintf("Transpile gremlin: %d:%d — %s", e.Line, e.Col, e.Message)
}
