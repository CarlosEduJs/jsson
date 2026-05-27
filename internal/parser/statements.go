package parser

import (
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT,
		token.UUID, token.EMAIL, token.URL,
		token.IPV4, token.IPV6, token.FILEPATH,
		token.DATE, token.DATETIME, token.REGEX:
		switch p.peekToken.Type {
		case token.DECLARE:
			return p.parseVariableDeclaration()
		case token.ASSIGN:
			return p.parseAssignment()
		case token.LBRACE:
			return p.parseObjectStatement()
		case token.LBRACKET:
			return p.parseArrayTemplateStatement()
		default:
			return nil
		}
	case token.AT:
		if p.peekToken.Type == token.PRESET {
			return p.parsePresetStatement()
		}
		return nil
	case token.INCLUDE:
		return p.parseIncludeStatement()
	default:
		return nil
	}
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}

	p.nextToken() // consume include

	if p.curToken.Type != token.STRING && p.curToken.Type != token.RAWSTRING {
		p.addError(ie.IncludePathExpected())

		return nil
	}

	stmt.Path = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

func (p *Parser) parsePresetStatement() ast.Statement {
	stmt := &ast.PresetStatement{Token: p.curToken}

	p.nextToken() // consume @

	if p.curToken.Type != token.PRESET {
		p.addError(ie.ExpectedToken(token.PRESET, p.curToken.Literal))

		return nil
	}

	p.nextToken() // consume 'preset'

	if p.curToken.Type != token.STRING && p.curToken.Type != token.RAWSTRING {
		p.addError("expected preset name as string")

		return nil
	}

	stmt.Name = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume preset name

	if p.curToken.Type != token.LBRACE {
		p.addError(ie.ExpectedToken(token.LBRACE, p.curToken.Literal))

		return nil
	}

	bodyExpr := p.parseObjectLiteral()
	if obj, ok := bodyExpr.(*ast.ObjectLiteral); ok {
		stmt.Body = obj
	} else {
		p.addError("expected object literal for preset body")

		return nil
	}

	return stmt
}

func (p *Parser) parseAssignment() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	p.nextToken() // consume ASSIGN

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	stmt := &ast.VariableDeclaration{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	p.nextToken() // consume DECLARE (:=)

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseObjectStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseArrayTemplateStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	stmt.Value = p.parseArrayTemplate()

	return stmt
}
