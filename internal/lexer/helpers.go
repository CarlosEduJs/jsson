package lexer

import (
	"fmt"
	ie "jsson/internal/errors"
	"unicode"
	"unicode/utf8"
)

func (l *Lexer) lexErrMsg(msg string) string {
	if l.SourceFile != "" {
		ctx := ie.FormatContext(l.SourceFile, l.line, l.column)

		return fmt.Sprintf("Lex goblin: %s — %s", ctx, msg)
	}

	ctx := fmt.Sprintf("%d:%d", l.line, l.column)

	return fmt.Sprintf("Lex goblin: %s — %s", ctx, msg)
}

func (l *Lexer) Errors() []string {
	return l.errors
}

func (l *Lexer) SetSourceFile(path string) {
	l.SourceFile = path
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || (l.ch == '/' && l.peekChar() == '/') {
		if l.ch == '/' && l.peekChar() == '/' {
			l.skipComment()

			continue
		}

		if l.ch == '\n' {
			l.line++
			l.column = 0
		}

		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.readChar()
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])

	return r
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func containsDot(s string) bool {
	for _, ch := range s {
		if ch == '.' {
			return true
		}
	}

	return false
}
