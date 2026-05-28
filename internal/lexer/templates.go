package lexer

import (
	ie "jsson/internal/errors"
	"jsson/internal/token"
)

func (l *Lexer) readTemplateString() (string, bool) {
	l.readChar()

	var runes []rune

	for {
		if l.ch == '`' {
			l.readChar()

			return string(runes), true
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

func (l *Lexer) readTripleBacktickString() (string, bool) {
	l.readChar()

	var runes []rune

	for {
		if l.ch == '`' && l.peekChar() == '`' {
			l.readChar()

			if l.peekChar() == '`' {
				l.readChar()
				l.readChar()

				return string(runes), true
			}

			runes = append(runes, '`', l.ch)
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

func (l *Lexer) lexBacktickString() token.Token {
	var tok token.Token

	if l.peekChar() == '`' {
		l.readChar()

		if l.peekChar() == '`' {
			l.readChar()
			lit, ok := l.readTripleBacktickString()
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

		tok.Type = token.TEMPLATESTR
		tok.Literal = ""
		tok.Line = l.line
		tok.Column = l.column

		return tok
	}

	lit, ok := l.readTemplateString()
	tok.Line = l.line
	tok.Column = l.column

	if !ok {
		msg := l.lexErrMsg(ie.UnterminatedString())
		l.errors = append(l.errors, msg)
		tok = l.newToken(token.ILLEGAL, msg)
	} else {
		tok.Type = token.TEMPLATESTR
		tok.Literal = lit
	}

	return tok
}
