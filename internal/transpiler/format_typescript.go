package transpiler

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

// TranspileToTypeScript converts the transpiled data to TypeScript format with types.
func (t *Transpiler) TranspileToTypeScript(ctx context.Context) ([]byte, error) {
	root, err := t.buildRootMap(ctx)
	if err != nil {
		return nil, err
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
