package parser

import (
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/lexer"
	"jsson/internal/token"
)

const (
	_ int = iota
	LOWEST
	TERNARY     // ? :
	LOGICAL     // && ||
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SUM         // + -
	PRODUCT     // * / %
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	RANGEEND    // Used internally for parsing range end (allows arithmetic, stops before map)
	RANGE       // .. (higher than arithmetic so i..i+2 works)
	MAP         // map (higher than RANGE so range map works)
	INDEX       // array[index] or obj.prop
)

var precedences = map[token.Type]int{
	token.LAND:     LOGICAL,
	token.LOR:      LOGICAL,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MODULO:   PRODUCT,
	token.QUESTION: TERNARY,
	token.DOT:      INDEX,
	token.RANGE:    RANGE,
	token.MAP:      MAP,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []error
}

func (p *Parser) addError(msg string) {
	sourceFile := ""
	if p.l != nil {
		sourceFile = p.l.SourceFile
	}

	p.errors = append(p.errors, &ie.ParseError{
		SourceFile: sourceFile,
		Line:       p.curToken.Line,
		Col:        p.curToken.Column,
		Message:    msg,
	})
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []error{}}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.parsePrefix()
	if prefix == nil {
		return nil
	}

	for p.peekToken.Type != token.EOF && precedence < p.peekPrecedence() {
		infix := p.parseInfix(prefix)
		if infix == nil {
			return prefix
		}

		prefix = infix
	}

	return prefix
}

func (p *Parser) parsePrefix() ast.Expression {
	switch p.curToken.Type {
	case token.IDENT:
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.INT:
		return p.parseIntegerLiteral()
	case token.FLOAT:
		return p.parseFloatLiteral()
	case token.STRING, token.RAWSTRING, token.TEMPLATESTR:
		return p.parseStringLiteral()
	case token.TRUE, token.FALSE, token.YES, token.NO, token.ON, token.OFF:
		return p.parseBooleanLiteral()
	case token.LPAREN:
		return p.parseGroupedExpression()
	case token.LBRACKET:
		return p.parseArrayLiteral()
	case token.LBRACE:
		return p.parseObjectLiteral()
	case token.MINUS:
		return p.parsePrefixExpression()
	case token.AT:
		return p.parseAtExpression()
	default:
		return nil
	}
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	switch p.peekToken.Type {
	case token.PLUS, token.MINUS, token.SLASH, token.ASTERISK, token.MODULO,
		token.EQ, token.NEQ, token.LT, token.GT, token.LTE, token.GTE,
		token.LAND, token.LOR:
		p.nextToken()

		return p.parseBinaryExpression(left)
	case token.QUESTION:
		p.nextToken()

		return p.parseConditionalExpression(left)
	case token.DOT:
		p.nextToken()

		return p.parseMemberExpression(left)
	case token.RANGE:
		p.nextToken()

		return p.parseRangeExpression(left)
	case token.MAP:
		p.nextToken()

		return p.parseMapExpression(left)
	default:
		return nil
	}
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) isValidPropertyName() bool {
	switch p.curToken.Type {
	case token.IDENT,
		token.STRING, token.RAWSTRING,
		token.TRUE, token.FALSE,
		token.YES, token.NO, token.ON, token.OFF,
		token.TEMPLATE, token.MAP, token.INCLUDE, token.STEP,
		token.PRESET, token.USE,
		token.UUID, token.EMAIL, token.URL,
		token.IPV4, token.IPV6, token.FILEPATH,
		token.DATE, token.DATETIME, token.REGEX:
		return true
	default:
		return false
	}
}
