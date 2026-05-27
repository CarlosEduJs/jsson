package parser

import (
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/token"
	"strconv"
)

func (p *Parser) parseRangeExpression(left ast.Expression) ast.Expression {
	expr := &ast.RangeExpression{Token: p.curToken, Start: left}

	p.nextToken()
	expr.End = p.parseExpression(MAP)

	if p.peekToken.Type == token.STEP {
		p.nextToken() // move to STEP
		p.nextToken() // move to step value

		switch p.curToken.Type {
		case token.MINUS:
			expr.Step = p.parsePrefixExpression()
		case token.INT:
			expr.Step = p.parseIntegerLiteral()
		default:
			expr.Step = p.parseExpression(LOWEST)
		}
	}

	return expr
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	expr := &ast.MemberExpression{Token: p.curToken, Left: left}
	p.nextToken() // consume .

	if p.curToken.Type != token.IDENT {
		p.addError(ie.ExpectedIdentifierAfterDot())

		return nil
	}

	expr.Property = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return expr
}

func (p *Parser) parseConditionalExpression(condition ast.Expression) ast.Expression {
	expr := &ast.ConditionalExpression{
		Token:     p.curToken,
		Condition: condition,
	}

	p.nextToken() // consume ?
	expr.Consequence = p.parseExpression(TERNARY - 1)

	if p.peekToken.Type != token.COLON {
		p.addError(ie.MissingColonInTernary())

		return nil
	}

	p.nextToken() // move to :
	p.nextToken() // consume :
	expr.Alternative = p.parseExpression(TERNARY - 1)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		p.addError(ie.MissingClosingParen())

		return nil
	}

	p.nextToken()

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	p.nextToken() // consume MINUS

	switch p.curToken.Type {
	case token.INT:
		lit := &ast.IntegerLiteral{Token: p.curToken}

		value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
		if err != nil {
			p.addError(ie.IntegerTooSpicy(p.curToken.Literal))

			return nil
		}

		lit.Value = -value

		return lit
	case token.FLOAT:
		lit := &ast.FloatLiteral{Token: p.curToken}

		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.addError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal))

			return nil
		}

		lit.Value = -value

		return lit
	}

	expr := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: "-",
		Left:     &ast.IntegerLiteral{Token: p.curToken, Value: 0},
		Right:    p.parseExpression(PREFIX),
	}

	return expr
}

func (p *Parser) parseMapExpression(left ast.Expression) ast.Expression {
	expression := &ast.MapExpression{Token: p.curToken, Left: left}

	if p.peekToken.Type != token.LPAREN {
		p.addError(ie.ExpectedToken(token.LPAREN, p.peekToken.Literal))

		return nil
	}

	p.nextToken() // consume map, now cur is (

	if p.peekToken.Type != token.IDENT {
		p.addError(ie.ExpectedToken(token.IDENT, p.peekToken.Literal))

		return nil
	}

	p.nextToken() // consume (, now cur is IDENT
	expression.Iterator = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != token.RPAREN {
		p.addError(ie.ExpectedToken(token.RPAREN, p.peekToken.Literal))

		return nil
	}

	p.nextToken() // consume IDENT, now cur is )

	if p.peekToken.Type != token.ASSIGN {
		p.addError(ie.ExpectedToken(token.ASSIGN, p.peekToken.Literal))

		return nil
	}

	p.nextToken() // consume ), now cur is =

	p.nextToken() // consume =, now cur is start of expression

	expression.Body = p.parseExpression(LOWEST)

	return expression
}
