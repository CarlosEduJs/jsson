package lexer

import (
	ie "jsson/internal/errors"
	"jsson/internal/token"
	"unicode/utf8"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
	errors       []string
	SourceFile   string
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0, errors: []string{}}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, width := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += width
	}

	l.column++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.ASSIGN, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.LT, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.GT, string(l.ch))
		}
	case '?':
		tok = l.newToken(token.QUESTION, string(l.ch))
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DECLARE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.COLON, string(l.ch))
		}
	case ',':
		tok = l.newToken(token.COMMA, string(l.ch))
	case '{':
		tok = l.newToken(token.LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(token.RBRACE, string(l.ch))
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.ch))
	case '(':
		tok = l.newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.ch))
	case '+':
		tok = l.newToken(token.PLUS, string(l.ch))
	case '-':
		tok = l.newToken(token.MINUS, string(l.ch))
	case '/':
		tok = l.newToken(token.SLASH, string(l.ch))
	case '*':
		tok = l.newToken(token.ASTERISK, string(l.ch))
	case '%':
		tok = l.newToken(token.MODULO, string(l.ch))
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LAND, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		}
	case '"':
		return l.lexQuotedString()
	case '`':
		return l.lexBacktickString()
	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.RANGE, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.DOT, string(l.ch))
		}
	case '@':
		tok = l.newToken(token.AT, string(l.ch))
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column

			return tok
		}

		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			if containsDot(tok.Literal) {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}

			tok.Line = l.line
			tok.Column = l.column

			return tok
		}

		msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
		l.errors = append(l.errors, msg)
		tok = l.newToken(token.ILLEGAL, msg)
	}

	l.readChar()

	return tok
}

func (l *Lexer) newToken(tokenType token.Type, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column}
}
