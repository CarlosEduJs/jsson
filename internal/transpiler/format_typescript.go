package transpiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"os"
	"path/filepath"
	"strings"
)

// TranspileToTypeScript converts the transpiled data to TypeScript format with types.
func (t *Transpiler) TranspileToTypeScript() ([]byte, error) {
	root := make(map[string]any)

	for _, stmt := range t.program.Statements {
		switch s := stmt.(type) {
		case *ast.PresetStatement:
			// Preset definitions are stored in preset table but not added to output
			t.presetTable[s.Name.Value] = s.Body
		case *ast.VariableDeclaration:
			// Variable declarations are stored in symbol table but not added to output
			val, err := t.evalExpression(s.Value, nil)
			if err != nil {
				return nil, err
			}

			t.symbolTable[s.Name.Value] = val
		case *ast.AssignmentStatement:
			key := s.Name.Value

			val, err := t.evalExpression(s.Value, nil)
			if err != nil {
				return nil, err
			}

			root[key] = val
		case *ast.IncludeStatement:
			// Handle includes (same logic as other transpilers)
			includePath := s.Path.Value

			var includeAbs string
			if filepath.IsAbs(includePath) {
				includeAbs = filepath.Clean(includePath)
			} else {
				includeAbs = filepath.Clean(filepath.Join(t.baseDir, includePath))
			}

			if t.inProgress[includeAbs] {
				return nil, t.errfNodeMsg(s, ie.CyclicInclude(includeAbs))
			}

			if cached, ok := t.includeCache[includeAbs]; ok {
				for k, v := range cached {
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				}

				break
			}

			t.inProgress[includeAbs] = true

			data, err := os.ReadFile(includeAbs)
			if err != nil {
				t.inProgress[includeAbs] = false

				return nil, t.errfNode(s, "could not read include file %q — gremlin can't find it: %v", s.Path.Value, err)
			}

			l := lexer.New(string(data))
			l.SetSourceFile(includeAbs)
			p := parser.New(l)

			prog := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.inProgress[includeAbs] = false

				return nil, t.errfNode(s, "parser errors in included file %q — wizard got confused: %v", s.Path.Value, p.Errors())
			}

			incBase := filepath.Dir(includeAbs)
			incT := New(prog, incBase, t.mergeMode, includeAbs)
			incT.includeCache = t.includeCache
			incT.inProgress = t.inProgress

			incJSON, err := incT.Transpile()
			if err != nil {
				t.inProgress[includeAbs] = false

				return nil, t.errfNode(s, "transpile error in included file %q: %v", s.Path.Value, err)
			}

			var incRoot map[string]any
			if err := json.Unmarshal(incJSON, &incRoot); err != nil {
				t.inProgress[includeAbs] = false

				return nil, t.errfNode(s, "invalid json from include %q: %v", s.Path.Value, err)
			}

			t.includeCache[includeAbs] = incRoot
			t.inProgress[includeAbs] = false

			for k, v := range incRoot {
				switch t.mergeMode {
				case mergeModeKeep:
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				case mergeModeOverwrite:
					root[k] = v
				case mergeModeError:
					if _, exists := root[k]; exists {
						return nil, t.errfNode(s, "include merge conflict for key %q from %s", k, includeAbs)
					}

					root[k] = v
				default:
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				}
			}
		}
	}

	// Convert any RangeResult to plain arrays
	if r, ok := t.convertRangeResults(root).(map[string]any); ok {
		root = r
	}

	// Generate TypeScript code
	var buf bytes.Buffer

	// Write exports for each top-level key
	for key, value := range root {
		fmt.Fprintf(&buf, "export const %s = ", key)
		writeTypeScriptValue(&buf, value, 0)
		buf.WriteString(" as const;\n\n")
	}

	// Generate type exports
	buf.WriteString("// Generated types\n")

	for key := range root {
		fmt.Fprintf(&buf, "export type %s = typeof %s;\n", capitalize(key), key)
	}

	return buf.Bytes(), nil
}

func writeTypeScriptValue(buf *bytes.Buffer, value any, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch v := value.(type) {
	case string:
		// Escape quotes and write as string literal
		escaped := strings.ReplaceAll(v, "\"", "\\\"")
		fmt.Fprintf(buf, "\"%s\"", escaped)
	case int64:
		fmt.Fprintf(buf, "%d", v)
	case float64:
		fmt.Fprintf(buf, "%v", v)
	case bool:
		fmt.Fprintf(buf, "%t", v)
	case nil:
		buf.WriteString("null")
	case map[string]any:
		buf.WriteString("{\n")

		first := true
		for k, val := range v {
			if !first {
				buf.WriteString(",\n")
			}

			first = false

			fmt.Fprintf(buf, "%s  %s: ", indentStr, k)
			writeTypeScriptValue(buf, val, indent+1)
		}

		fmt.Fprintf(buf, "\n%s}", indentStr)
	case []any:
		buf.WriteString("[\n")

		for i, val := range v {
			buf.WriteString(indentStr + "  ")
			writeTypeScriptValue(buf, val, indent+1)

			if i < len(v)-1 {
				buf.WriteString(",")
			}

			buf.WriteString("\n")
		}

		buf.WriteString(indentStr + "]")
	default:
		fmt.Fprintf(buf, "%v", v)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
