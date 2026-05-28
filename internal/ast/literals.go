package ast

import (
	"bytes"
	"jsson/internal/token"
)

// IntegerLiteral represents an integer literal in the AST.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()           {}
func (il *IntegerLiteral) TokenLiteral() string      { return il.Token.Literal }
func (il *IntegerLiteral) Position() (line, col int) { return il.Token.Line, il.Token.Column }
func (il *IntegerLiteral) String() string            { return il.Token.Literal }

// FloatLiteral represents a float literal in the AST.
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()           {}
func (fl *FloatLiteral) TokenLiteral() string      { return fl.Token.Literal }
func (fl *FloatLiteral) Position() (line, col int) { return fl.Token.Line, fl.Token.Column }
func (fl *FloatLiteral) String() string            { return fl.Token.Literal }

// BooleanLiteral represents a boolean literal in the AST.
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()           {}
func (b *BooleanLiteral) TokenLiteral() string      { return b.Token.Literal }
func (b *BooleanLiteral) Position() (line, col int) { return b.Token.Line, b.Token.Column }
func (b *BooleanLiteral) String() string            { return b.Token.Literal }

// StringLiteral represents a string literal in the AST.
type StringLiteral struct {
	Token token.Token
	Value string
	IsRaw bool
}

func (sl *StringLiteral) expressionNode()           {}
func (sl *StringLiteral) TokenLiteral() string      { return sl.Token.Literal }
func (sl *StringLiteral) Position() (line, col int) { return sl.Token.Line, sl.Token.Column }
func (sl *StringLiteral) String() string            { return sl.Token.Literal }

// InterpolatedPart is either a literal text string or an Expression.
type InterpolatedPart interface {
	isPart()
}

// TextPart is a literal text segment within an interpolated string.
type TextPart struct {
	Value string
}

func (TextPart) isPart() {}

// ExprPart is an expression segment within an interpolated string.
type ExprPart struct {
	Expr Expression
}

func (ExprPart) isPart() {}

// InterpolatedString represents a string with ${var} or {var} interpolations.
type InterpolatedString struct {
	Token token.Token
	Parts []InterpolatedPart
}

func (is *InterpolatedString) expressionNode()           {}
func (is *InterpolatedString) TokenLiteral() string      { return is.Token.Literal }
func (is *InterpolatedString) Position() (line, col int) { return is.Token.Line, is.Token.Column }
func (is *InterpolatedString) String() string {
	var out bytes.Buffer

	for _, part := range is.Parts {
		switch p := part.(type) {
		case TextPart:
			out.WriteString(p.Value)
		case ExprPart:
			out.WriteString("{")
			out.WriteString(p.Expr.String())
			out.WriteString("}")
		}
	}

	return out.String()
}

// ObjectLiteral: { key = value }.
type ObjectLiteral struct {
	Token        token.Token            // '{'
	Declarations []*VariableDeclaration // Local variables (key := value)
	Properties   map[string]Expression  // Properties (key = value)
	Keys         []string               // Para manter a ordem das chaves
}

func (o *ObjectLiteral) expressionNode()           {}
func (o *ObjectLiteral) TokenLiteral() string      { return o.Token.Literal }
func (o *ObjectLiteral) Position() (line, col int) { return o.Token.Line, o.Token.Column }
func (o *ObjectLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")

	for _, key := range o.Keys {
		out.WriteString(key)

		if val := o.Properties[key]; val != nil {
			out.WriteString(" = ")
			out.WriteString(val.String())
		}

		out.WriteString(", ")
	}

	out.WriteString(" }")

	return out.String()
}

// ArrayLiteral: [ 1, 2, 3 ].
type ArrayLiteral struct {
	Token    token.Token // '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()           {}
func (al *ArrayLiteral) TokenLiteral() string      { return al.Token.Literal }
func (al *ArrayLiteral) Position() (line, col int) { return al.Token.Line, al.Token.Column }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	for i, el := range al.Elements {
		out.WriteString(el.String())

		if i < len(al.Elements)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString("]")

	return out.String()
}
