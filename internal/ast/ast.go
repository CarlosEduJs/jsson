package ast

import (
	"bytes"
)

// Node is the base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
	String() string
	// Position returns the line and column of the node in the source file
	Position() (line, col int)
}

// Statement is a node that represents a statement.
type Statement interface {
	Node
	statementNode()
}

// Expression is a node that represents an expression.
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) Position() (line, col int) {
	if len(p.Statements) > 0 {
		return p.Statements[0].Position()
	}

	return 0, 0
}
