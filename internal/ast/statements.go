package ast

import (
	"bytes"
	"jsson/internal/token"
)

// Assignment: name = "value".
type AssignmentStatement struct {
	Token token.Token // the token.IDENT
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()            {}
func (as *AssignmentStatement) TokenLiteral() string      { return as.Token.Literal }
func (as *AssignmentStatement) Position() (line, col int) { return as.Token.Line, as.Token.Column }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	return out.String()
}

// VariableDeclaration: name := value.
type VariableDeclaration struct {
	Token token.Token // the ':=' token
	Name  *Identifier
	Value Expression
}

func (vd *VariableDeclaration) statementNode()            {}
func (vd *VariableDeclaration) TokenLiteral() string      { return vd.Token.Literal }
func (vd *VariableDeclaration) Position() (line, col int) { return vd.Token.Line, vd.Token.Column }
func (vd *VariableDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(vd.Name.String())
	out.WriteString(" := ")

	if vd.Value != nil {
		out.WriteString(vd.Value.String())
	}

	return out.String()
}

// IncludeStatement: include "file.jsson".
type IncludeStatement struct {
	Token token.Token // the 'include' token
	Path  *StringLiteral
}

func (is *IncludeStatement) statementNode()            {}
func (is *IncludeStatement) TokenLiteral() string      { return is.Token.Literal }
func (is *IncludeStatement) Position() (line, col int) { return is.Token.Line, is.Token.Column }
func (is *IncludeStatement) String() string {
	return "include " + is.Path.String()
}

// PresetStatement: @preset "name" { ... }
// Defines a reusable configuration preset.
type PresetStatement struct {
	Token token.Token    // The '@' token
	Name  *StringLiteral // Preset name
	Body  *ObjectLiteral // Preset contents
}

func (ps *PresetStatement) statementNode()            {}
func (ps *PresetStatement) TokenLiteral() string      { return ps.Token.Literal }
func (ps *PresetStatement) Position() (line, col int) { return ps.Token.Line, ps.Token.Column }
func (ps *PresetStatement) String() string {
	var out bytes.Buffer

	out.WriteString("@preset ")
	out.WriteString(ps.Name.String())
	out.WriteString(" ")
	out.WriteString(ps.Body.String())

	return out.String()
}
