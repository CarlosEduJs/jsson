package parser

import (
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/token"
	"strconv"
	"strings"
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(ie.IntegerTooSpicy(p.curToken.Literal))

		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal))

		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	isRaw := p.curToken.Type == token.RAWSTRING
	isTemplate := p.curToken.Type == token.TEMPLATESTR
	value := p.curToken.Literal

	if isTemplate && strings.Contains(value, "${") {
		return p.parseTemplateString(value)
	}

	if isRaw && strings.Contains(value, "{") {
		return p.parseInterpolatedString(value)
	}

	return &ast.StringLiteral{
		Token: p.curToken,
		Value: value,
		IsRaw: isRaw || isTemplate,
	}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	value := p.curToken.Type == token.TRUE ||
		p.curToken.Type == token.YES ||
		p.curToken.Type == token.ON

	return &ast.BooleanLiteral{Token: p.curToken, Value: value}
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.curToken}
	obj.Properties = make(map[string]ast.Expression)
	obj.Keys = []string{}
	obj.Declarations = []*ast.VariableDeclaration{}

	p.nextToken() // consume {

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		if !p.isValidPropertyName() {
			p.nextToken()

			continue
		}

		key := p.curToken.Literal
		if p.curToken.Type == token.STRING || p.curToken.Type == token.RAWSTRING {
			if len(key) >= 2 && key[0] == '"' && key[len(key)-1] == '"' {
				key = key[1 : len(key)-1]
			}
		}

		p.nextToken() // consume key

		switch p.curToken.Type {
		case token.DECLARE:
			p.nextToken() // consume :=
			val := p.parseExpression(LOWEST)
			decl := &ast.VariableDeclaration{
				Token: p.curToken,
				Name:  &ast.Identifier{Value: key},
				Value: val,
			}
			obj.Declarations = append(obj.Declarations, decl)

			p.nextToken() // consume value
		case token.ASSIGN, token.COLON:
			obj.Keys = append(obj.Keys, key)

			p.nextToken() // consume = or :
			val := p.parseExpression(LOWEST)
			obj.Properties[key] = val

			p.nextToken() // consume value
		case token.LBRACE:
			obj.Keys = append(obj.Keys, key)
			val := p.parseExpression(LOWEST)
			obj.Properties[key] = val

			p.nextToken()
		case token.LBRACKET:
			obj.Keys = append(obj.Keys, key)
			val := p.parseArrayTemplate()
			obj.Properties[key] = val

			p.nextToken()
		default:
			obj.Keys = append(obj.Keys, key)
			obj.Properties[key] = nil
		}

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	if p.curToken.Type != token.RBRACE {
		p.addError(ie.MissingClosingBrace())
	}

	return obj
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = []ast.Expression{}

	p.nextToken() // consume [

	for p.curToken.Type != token.RBRACKET && p.curToken.Type != token.EOF {
		elem := p.parseExpression(LOWEST)
		if elem != nil {
			array.Elements = append(array.Elements, elem)
		}

		p.nextToken()

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	return array
}
