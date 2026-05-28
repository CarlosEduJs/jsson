package parser

import (
	"jsson/internal/ast"
	"jsson/internal/lexer"
	"jsson/internal/token"
	"strconv"
	"strings"
)

func (p *Parser) parseAtExpression() ast.Expression {
	p.nextToken() // consume @

	switch p.curToken.Type {
	case token.STRING, token.RAWSTRING:
		return p.parsePresetReferenceAfterAt()
	case token.USE:
		p.nextToken()

		return p.parsePresetReferenceAfterAt()
	case token.UUID, token.EMAIL, token.URL, token.IPV4, token.IPV6,
		token.FILEPATH, token.DATE, token.DATETIME, token.REGEX,
		token.VINT, token.VFLOAT, token.VBOOL:
		return p.parseValidator()
	default:
		p.addError("unexpected token after @: " + p.curToken.Literal)

		return nil
	}
}

func (p *Parser) parsePresetReferenceAfterAt() ast.Expression {
	ref := &ast.PresetReference{Token: p.curToken}
	ref.Name = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type == token.LBRACE {
		p.nextToken() // move to {

		overridesExpr := p.parseObjectLiteral()
		if obj, ok := overridesExpr.(*ast.ObjectLiteral); ok {
			ref.Overrides = obj
		}
	}

	return ref
}

func (p *Parser) parseValidator() ast.Expression {
	validator := &ast.ValidatorExpression{
		Token: p.curToken,
		Type:  strings.ToLower(p.curToken.Literal),
	}

	if p.peekToken.Type == token.LPAREN {
		p.nextToken() // consume validator name
		p.nextToken() // consume (

		args := []ast.Expression{}

		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			switch p.curToken.Type {
			case token.STRING, token.RAWSTRING:
				validator.Pattern = p.curToken.Literal
				args = append(args, &ast.StringLiteral{
					Token: p.curToken,
					Value: p.curToken.Literal,
				})
			case token.INT:
				val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
				if err == nil {
					args = append(args, &ast.IntegerLiteral{
						Token: p.curToken,
						Value: val,
					})
				}
			case token.FLOAT:
				val, err := strconv.ParseFloat(p.curToken.Literal, 64)
				if err == nil {
					args = append(args, &ast.FloatLiteral{
						Token: p.curToken,
						Value: val,
					})
				}
			}

			p.nextToken()

			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}

		validator.Args = args
	}

	return validator
}

func (p *Parser) parseArrayTemplate() ast.Expression {
	at := &ast.ArrayTemplate{Token: p.curToken}
	p.nextToken() // consume [

	hasTemplate := p.curToken.Type == token.TEMPLATE

	if hasTemplate {
		p.nextToken() // consume template
		templateObj := p.parseObjectLiteral()

		var templateOk bool

		at.Template, templateOk = templateObj.(*ast.ObjectLiteral)
		if !templateOk {
			return at
		}

		p.nextToken() // consume }
	}

	if p.curToken.Type == token.MAP {
		at.Map = p.parseMapClause()

		if !hasTemplate && at.Map != nil {
			at.Template = &ast.ObjectLiteral{
				Token:      at.Map.Token,
				Properties: make(map[string]ast.Expression),
				Keys:       []string{at.Map.Param.Value},
			}
			at.Template.Properties[at.Map.Param.Value] = at.Map.Param
		}
	}

	if at.Template == nil {
		p.addError("array must have either 'template' definition or 'map' clause")

		return at
	}

	at.Rows = [][]ast.Expression{}
	expectedCols := len(at.Template.Keys)

	for p.curToken.Type != token.RBRACKET && p.curToken.Type != token.EOF {
		for p.curToken.Type == token.RBRACE {
			p.nextToken()
		}

		row := []ast.Expression{}

		for range expectedCols {
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}

			if p.curToken.Type == token.RBRACKET {
				break
			}

			expr := p.parseExpression(LOWEST)
			if expr != nil {
				row = append(row, expr)
			}

			p.nextToken()
		}

		if len(row) > 0 {
			at.Rows = append(at.Rows, row)
		}

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	return at
}

func (p *Parser) parseMapClause() *ast.MapClause {
	mc := &ast.MapClause{Token: p.curToken}
	p.nextToken() // consume map
	p.nextToken() // consume (
	mc.Param = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // consume param
	p.nextToken() // consume )
	p.nextToken() // consume =
	bodyExpr := p.parseObjectLiteral()

	var bodyOk bool

	mc.Body, bodyOk = bodyExpr.(*ast.ObjectLiteral)
	if !bodyOk {
		return mc
	}

	p.nextToken() // consume }

	return mc
}

func (p *Parser) parseInterpolatedCommon(content string, isTemplate bool) ast.Expression {
	interp := &ast.InterpolatedString{
		Token: p.curToken,
		Parts: []ast.InterpolatedPart{},
	}

	var currentText strings.Builder

	i := 0

	for i < len(content) {
		found := false
		openLen := 0

		if isTemplate && i < len(content)-1 && content[i] == '$' && content[i+1] == '{' {
			found = true
			openLen = 2
		} else if !isTemplate && content[i] == '{' {
			found = true
			openLen = 1
		}

		if found {
			if currentText.Len() > 0 {
				interp.Parts = append(interp.Parts, ast.TextPart{Value: currentText.String()})
				currentText.Reset()
			}

			depth := 1
			start := i + openLen
			i += openLen

			for i < len(content) && depth > 0 {
				switch content[i] {
				case '{':
					depth++
				case '}':
					depth--
				}

				i++
			}

			if depth == 0 {
				exprText := content[start:i]
				exprLexer := lexer.New(exprText)
				exprParser := New(exprLexer)
				expr := exprParser.parseExpression(LOWEST)

				if expr != nil && len(exprParser.Errors()) == 0 {
					interp.Parts = append(interp.Parts, ast.ExprPart{Expr: expr})
				} else {
					if isTemplate {
						currentText.WriteString("${")
					} else {
						currentText.WriteString("{")
					}

					currentText.WriteString(exprText)
					currentText.WriteString("}")
				}
			} else {
				currentText.WriteString(content[start-openLen : i])
			}
		} else {
			currentText.WriteByte(content[i])
			i++
		}
	}

	if currentText.Len() > 0 {
		interp.Parts = append(interp.Parts, ast.TextPart{Value: currentText.String()})
	}

	return interp
}

func (p *Parser) parseTemplateString(content string) ast.Expression {
	return p.parseInterpolatedCommon(content, true)
}

func (p *Parser) parseInterpolatedString(content string) ast.Expression {
	return p.parseInterpolatedCommon(content, false)
}
