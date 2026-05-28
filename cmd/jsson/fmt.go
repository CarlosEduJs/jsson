package main

import (
	"flag"
	"fmt"
	"jsson/internal/ast"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"os"
	"strings"
)

func runFormatter() {
	writePtr := flag.Bool("w", false, "Write result to source file instead of stdout")
	checkPtr := flag.Bool("check", false, "Check if file is already formatted, exit 1 if not")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: jsson fmt [-w] [-check] <file.jsson>")
		os.Exit(1)
	}

	for _, path := range args {
		formatFile(path, *writePtr, *checkPtr)
	}
}

func formatFile(path string, write, check bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
		os.Exit(1)
	}

	source := string(data)
	l := lexer.New(source)
	l.SetSourceFile(path)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "Parse errors in %s:\n", path)
		for _, e := range p.Errors() {
			fmt.Fprintf(os.Stderr, "\t%s\n", e.Error())
		}
		os.Exit(1)
	}

	formatted := formatProgram(program)

	if check {
		if formatted == source {
			return
		}
		fmt.Fprintf(os.Stderr, "%s: not formatted\n", path)
		os.Exit(1)
	}

	if write {
		if err := os.WriteFile(path, []byte(formatted), 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("%s: formatted\n", path)
	} else {
		fmt.Print(formatted)
	}
}

func formatProgram(program *ast.Program) string {
	var out strings.Builder

	for i, stmt := range program.Statements {
		if i > 0 {
			out.WriteString("\n")
		}

		formatStatement(&out, stmt, 0)
	}

	out.WriteString("\n")
	return out.String()
}

func formatStatement(out *strings.Builder, stmt ast.Statement, indent int) {
	indentStr := strings.Repeat("    ", indent)

	switch s := stmt.(type) {
	case *ast.PresetStatement:
		out.WriteString(indentStr)
		out.WriteString("@preset \"")
		out.WriteString(s.Name.Value)
		out.WriteString("\" ")
		formatObjectLiteral(out, s.Body, indent, true)
	case *ast.VariableDeclaration:
		out.WriteString(indentStr)
		out.WriteString(s.Name.Value)
		out.WriteString(" := ")
		formatExpression(out, s.Value, indent)
	case *ast.AssignmentStatement:
		out.WriteString(indentStr)
		out.WriteString(s.Name.Value)
		out.WriteString(" = ")
		formatExpression(out, s.Value, indent)
	case *ast.IncludeStatement:
		out.WriteString(indentStr)
		out.WriteString("@include \"")
		out.WriteString(s.Path.Value)
		out.WriteString("\"")
	}
}

func formatExpression(out *strings.Builder, expr ast.Expression, indent int) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.StringLiteral:
		if e.IsRaw {
			out.WriteString("`")
			out.WriteString(e.Value)
			out.WriteString("`")
		} else {
			out.WriteString("\"")
			out.WriteString(escapeString(e.Value))
			out.WriteString("\"")
		}
	case *ast.IntegerLiteral:
		out.WriteString(e.Token.Literal)
	case *ast.FloatLiteral:
		out.WriteString(e.Token.Literal)
	case *ast.BooleanLiteral:
		if e.Value {
			out.WriteString("true")
		} else {
			out.WriteString("false")
		}
	case *ast.Identifier:
		out.WriteString(e.Value)
	case *ast.ObjectLiteral:
		formatObjectLiteral(out, e, indent, false)
	case *ast.ArrayLiteral:
		formatArrayLiteral(out, e, indent)
	case *ast.ArrayTemplate:
		formatArrayTemplate(out, e, indent)
	case *ast.ValidatorExpression:
		out.WriteString("@")
		out.WriteString(string(e.Type))
		if e.Pattern != "" {
			out.WriteString("(\"")
			out.WriteString(escapeString(e.Pattern))
			out.WriteString("\")")
		} else if len(e.Args) > 0 {
			out.WriteString("(")
			for i, arg := range e.Args {
				if i > 0 {
					out.WriteString(", ")
				}
				formatExpression(out, arg, indent)
			}
			out.WriteString(")")
		}
	case *ast.PresetReference:
		out.WriteString("@use \"")
		out.WriteString(e.Name.Value)
		out.WriteString("\"")
		if e.Overrides != nil {
			out.WriteString(" ")
			formatObjectLiteral(out, e.Overrides, indent, true)
		}
	case *ast.RangeExpression:
		formatExpression(out, e.Start, indent)
		out.WriteString("..")
		formatExpression(out, e.End, indent)
		if e.Step != nil {
			out.WriteString(".")
			formatExpression(out, e.Step, indent)
		}
	case *ast.BinaryExpression:
		formatExpression(out, e.Left, indent)
		out.WriteString(" ")
		out.WriteString(e.Operator)
		out.WriteString(" ")
		formatExpression(out, e.Right, indent)
	case *ast.ConditionalExpression:
		formatExpression(out, e.Condition, indent)
		out.WriteString(" ? ")
		formatExpression(out, e.Consequence, indent)
		out.WriteString(" : ")
		formatExpression(out, e.Alternative, indent)
	case *ast.InterpolatedString:
		out.WriteString("`")
		for _, part := range e.Parts {
			switch p := part.(type) {
			case ast.TextPart:
				out.WriteString(p.Value)
			case ast.ExprPart:
				out.WriteString("{")
				formatExpression(out, p.Expr, indent)
				out.WriteString("}")
			}
		}
		out.WriteString("`")
	case *ast.MemberExpression:
		formatExpression(out, e.Left, indent)
		out.WriteString(".")
		out.WriteString(e.Property.Value)
	case *ast.MapExpression:
		formatExpression(out, e.Left, indent)
		out.WriteString(" -> ")
		out.WriteString(e.Iterator.Value)
		out.WriteString(" => ")
		formatExpression(out, e.Body, indent)
	default:
		out.WriteString(e.String())
	}
}

func formatObjectLiteral(out *strings.Builder, obj *ast.ObjectLiteral, indent int, inline bool) {
	hasDecls := len(obj.Declarations) > 0
	hasProps := len(obj.Properties) > 0

	if !hasDecls && !hasProps {
		out.WriteString("{ ")
		for i, key := range obj.Keys {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(key)
		}
		out.WriteString(" }")
		return
	}

	entries := len(obj.Declarations) + len(obj.Properties)

	if inline && entries <= 2 {
		out.WriteString("{ ")
		for _, decl := range obj.Declarations {
			out.WriteString(decl.Name.Value)
			out.WriteString(" := ")
			formatExpression(out, decl.Value, indent+1)
			out.WriteString("; ")
		}
		for i, key := range obj.Keys {
			if i > 0 || len(obj.Declarations) > 0 {
				out.WriteString("; ")
			}
			out.WriteString(key)
			if val := obj.Properties[key]; val != nil {
				out.WriteString(" = ")
				formatExpression(out, val, indent+1)
			}
		}
		out.WriteString("}")
		return
	}

	out.WriteString("{\n")
	nextIndent := indent + 1

	for _, decl := range obj.Declarations {
		out.WriteString(strings.Repeat("    ", nextIndent))
		out.WriteString(decl.Name.Value)
		out.WriteString(" := ")
		formatExpression(out, decl.Value, nextIndent)
		out.WriteString("\n")
	}

	for _, key := range obj.Keys {
		out.WriteString(strings.Repeat("    ", nextIndent))
		out.WriteString(key)

		if val := obj.Properties[key]; val != nil {
			out.WriteString(" = ")
			formatExpression(out, val, nextIndent)
		}

		out.WriteString(",\n")
	}

	out.WriteString(indentStr(indent))
	out.WriteString("}")
}

func formatArrayLiteral(out *strings.Builder, arr *ast.ArrayLiteral, indent int) {
	if len(arr.Elements) == 0 {
		out.WriteString("[]")
		return
	}

	out.WriteString("[\n")
	nextIndent := indent + 1

	for _, el := range arr.Elements {
		out.WriteString(strings.Repeat("    ", nextIndent))
		formatExpression(out, el, nextIndent)
		out.WriteString(",\n")
	}

	out.WriteString(strings.Repeat("    ", indent))
	out.WriteString("]")
}

func formatArrayTemplate(out *strings.Builder, at *ast.ArrayTemplate, indent int) {
	out.WriteString("[\n")
	nextIndent := indent + 1

	if at.Template != nil {
		out.WriteString(strings.Repeat("    ", nextIndent))
		formatObjectLiteral(out, at.Template, nextIndent, true)
		out.WriteString(",\n")
	}

	for _, row := range at.Rows {
		out.WriteString(strings.Repeat("    ", nextIndent))
		for j, expr := range row {
			if j > 0 {
				out.WriteString(", ")
			}
			formatExpression(out, expr, nextIndent)
		}
		out.WriteString(",\n")
	}

	if at.Map != nil {
		out.WriteString(strings.Repeat("    ", nextIndent))
		out.WriteString(at.Map.Param.Value)
		out.WriteString(" -> ")
		formatExpression(out, at.Map.Body, nextIndent)
		out.WriteString("\n")
	}

	out.WriteString(strings.Repeat("    ", indent))
	out.WriteString("]")
}

func indentStr(n int) string {
	return strings.Repeat("    ", n)
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
