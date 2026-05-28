package ast

import (
	"bytes"
	"jsson/internal/token"
)

// Identifier represents a named identifier in the AST.
type Identifier struct {
	Token token.Token // the token.IDENT
	Value string
}

func (i *Identifier) expressionNode()           {}
func (i *Identifier) TokenLiteral() string      { return i.Token.Literal }
func (i *Identifier) Position() (line, col int) { return i.Token.Line, i.Token.Column }
func (i *Identifier) String() string            { return i.Value }

// ValidatorType represents the type of a validator expression.
type ValidatorType string

const (
	ValidatorUUID     ValidatorType = "uuid"
	ValidatorEmail    ValidatorType = "email"
	ValidatorURL      ValidatorType = "url"
	ValidatorIPv4     ValidatorType = "ipv4"
	ValidatorIPv6     ValidatorType = "ipv6"
	ValidatorFilepath ValidatorType = "filepath"
	ValidatorDate     ValidatorType = "date"
	ValidatorDatetime ValidatorType = "datetime"
	ValidatorRegex    ValidatorType = "regex"
	ValidatorInt      ValidatorType = "int"
	ValidatorFloat    ValidatorType = "float"
	ValidatorBool     ValidatorType = "bool"
)

// ValidatorExpression represents a validator call like @uuid, @int(min, max).
type ValidatorExpression struct {
	Token   token.Token
	Type    ValidatorType
	Pattern string
	Args    []Expression // For validators like @int(min, max), @float(min, max)
}

func (ve *ValidatorExpression) expressionNode()           {}
func (ve *ValidatorExpression) TokenLiteral() string      { return ve.Token.Literal }
func (ve *ValidatorExpression) Position() (line, col int) { return ve.Token.Line, ve.Token.Column }
func (ve *ValidatorExpression) String() string {
	if ve.Pattern != "" {
		return "@" + string(ve.Type) + "(\"" + ve.Pattern + "\")"
	}

	if len(ve.Args) > 0 {
		return "@" + string(ve.Type) + "(...)"
	}

	return "@" + string(ve.Type)
}

// BinaryExpression: x + y.
type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode()           {}
func (be *BinaryExpression) TokenLiteral() string      { return be.Operator }
func (be *BinaryExpression) Position() (line, col int) { return be.Token.Line, be.Token.Column }
func (be *BinaryExpression) String() string {
	return "(" + be.Left.String() + " " + be.Operator + " " + be.Right.String() + ")"
}

// MemberExpression: item.path.
type MemberExpression struct {
	Token    token.Token // The '.' token
	Left     Expression
	Property *Identifier
}

func (me *MemberExpression) expressionNode()           {}
func (me *MemberExpression) TokenLiteral() string      { return me.Token.Literal }
func (me *MemberExpression) Position() (line, col int) { return me.Token.Line, me.Token.Column }
func (me *MemberExpression) String() string {
	return me.Left.String() + "." + me.Property.String()
}

// ConditionalExpression: condition ? consequence : alternative.
type ConditionalExpression struct {
	Token       token.Token // The '?' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ce *ConditionalExpression) expressionNode()           {}
func (ce *ConditionalExpression) TokenLiteral() string      { return ce.Token.Literal }
func (ce *ConditionalExpression) Position() (line, col int) { return ce.Token.Line, ce.Token.Column }
func (ce *ConditionalExpression) String() string {
	return "(" + ce.Condition.String() + " ? " + ce.Consequence.String() + " : " + ce.Alternative.String() + ")"
}

// RangeExpression: start .. end [ step N ].
type RangeExpression struct {
	Token token.Token
	Start Expression
	End   Expression
	Step  Expression // optional
}

func (re *RangeExpression) expressionNode()           {}
func (re *RangeExpression) TokenLiteral() string      { return re.Token.Literal }
func (re *RangeExpression) Position() (line, col int) { return re.Token.Line, re.Token.Column }
func (re *RangeExpression) String() string {
	if re.Step != nil {
		return re.Start.String() + ".." + re.End.String() + " step " + re.Step.String()
	}

	return re.Start.String() + ".." + re.End.String()
}

// MapClause: map (x) = { ... }.
type MapClause struct {
	Token token.Token    // "map"
	Param *Identifier    // "x"
	Body  *ObjectLiteral // "{ ... }"
}

func (mc *MapClause) expressionNode()           {}
func (mc *MapClause) TokenLiteral() string      { return mc.Token.Literal }
func (mc *MapClause) Position() (line, col int) { return mc.Token.Line, mc.Token.Column }
func (mc *MapClause) String() string {
	return "map (" + mc.Param.String() + ") = " + mc.Body.String()
}

// MapExpression: (collection map(x) = expr).
type MapExpression struct {
	Token    token.Token // The 'map' token
	Left     Expression  // The array being mapped
	Iterator *Identifier // The variable name (e.g. 't')
	Body     Expression  // The transformation body
}

func (me *MapExpression) expressionNode()           {}
func (me *MapExpression) TokenLiteral() string      { return me.Token.Literal }
func (me *MapExpression) Position() (line, col int) { return me.Token.Line, me.Token.Column }
func (me *MapExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(me.Left.String())
	out.WriteString(" map (")
	out.WriteString(me.Iterator.String())
	out.WriteString(") = ")
	out.WriteString(me.Body.String())
	out.WriteString(")")

	return out.String()
}

// ArrayTemplate: users [ template { name, age } ... ].
type ArrayTemplate struct {
	Token    token.Token // The identifier token before '['
	Name     *Identifier
	Template *ObjectLiteral // The template definition
	Map      *MapClause     // Optional map clause
	Rows     [][]Expression // The data rows
}

func (at *ArrayTemplate) expressionNode()           {}
func (at *ArrayTemplate) TokenLiteral() string      { return at.Token.Literal }
func (at *ArrayTemplate) Position() (line, col int) { return at.Token.Line, at.Token.Column }
func (at *ArrayTemplate) String() string {
	var out bytes.Buffer
	if at.Name != nil {
		out.WriteString(at.Name.String())
		out.WriteString(" ")
	}

	out.WriteString("[ template ")

	if at.Template != nil {
		out.WriteString(at.Template.String())
	}

	if at.Map != nil {
		out.WriteString(" ")
		out.WriteString(at.Map.String())
	}

	out.WriteString(" ... ]")

	return out.String()
}

// PresetReference: @"name" or @"name" { overrides }
// References and optionally extends a preset.
type PresetReference struct {
	Token     token.Token    // The '@' token
	Name      *StringLiteral // Preset name to reference
	Overrides *ObjectLiteral // Optional overrides
}

func (pr *PresetReference) expressionNode()           {}
func (pr *PresetReference) TokenLiteral() string      { return pr.Token.Literal }
func (pr *PresetReference) Position() (line, col int) { return pr.Token.Line, pr.Token.Column }
func (pr *PresetReference) String() string {
	var out bytes.Buffer

	out.WriteString("@")
	out.WriteString(pr.Name.String())

	if pr.Overrides != nil {
		out.WriteString(" ")
		out.WriteString(pr.Overrides.String())
	}

	return out.String()
}
