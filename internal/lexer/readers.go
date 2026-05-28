package lexer

import (
	ie "jsson/internal/errors"
	"jsson/internal/token"
)

func (l *Lexer) readIdentifier() string {
	start := l.position

	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	var runes []rune
	for isDigit(l.ch) {
		runes = append(runes, l.ch)
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		runes = append(runes, l.ch)
		l.readChar()

		for isDigit(l.ch) {
			runes = append(runes, l.ch)
			l.readChar()
		}
	}

	return string(runes)
}

func (l *Lexer) readString() (string, bool) {
	l.readChar()

	var runes []rune

	for {
		if l.ch == '"' {
			l.readChar()

			return string(runes), true
		}

		if l.ch == '\\' {
			l.readChar()

			switch l.ch {
			case 'n':
				runes = append(runes, '\n')
			case 't':
				runes = append(runes, '\t')
			case '"':
				runes = append(runes, '"')
			case '\\':
				runes = append(runes, '\\')
			default:
				runes = append(runes, '\\', l.ch)
			}

			l.readChar()

			continue
		}

		if l.ch == 0 {
			return "", false
		}

		runes = append(runes, l.ch)
		l.readChar()
	}
}

func (l *Lexer) readRawString() (string, bool) {
	l.readChar()

	var runes []rune

	for {
		if l.ch == '"' && l.peekChar() == '"' {
			l.readChar()

			if l.peekChar() == '"' {
				l.readChar()
				l.readChar()

				return string(runes), true
			}

			runes = append(runes, '"', l.ch)
			l.readChar()

			continue
		}

		if l.ch == 0 {
			return "", false
		}

		if l.ch == '\n' {
			l.line++
			l.column = 0
		}

		runes = append(runes, l.ch)
		l.readChar()
	}
}

func (l *Lexer) lexQuotedString() token.Token {
	var tok token.Token

	if l.peekChar() == '"' {
		l.readChar()

		if l.peekChar() == '"' {
			l.readChar()
			lit, ok := l.readRawString()
			tok.Line = l.line
			tok.Column = l.column

			if !ok {
				msg := l.lexErrMsg(ie.UnterminatedString())
				l.errors = append(l.errors, msg)
				tok = l.newToken(token.ILLEGAL, msg)
			} else {
				tok.Type = token.RAWSTRING
				tok.Literal = lit
			}

			return tok
		}

		tok.Type = token.STRING
		tok.Literal = ""
		tok.Line = l.line
		tok.Column = l.column
		l.readChar()

		return tok
	}

	lit, ok := l.readString()
	tok.Line = l.line
	tok.Column = l.column

	if !ok {
		msg := l.lexErrMsg(ie.UnterminatedString())
		l.errors = append(l.errors, msg)
		tok = l.newToken(token.ILLEGAL, msg)
	} else {
		tok.Type = token.STRING
		tok.Literal = lit
	}

	return tok
}
