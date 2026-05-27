package transpiler

import (
	"fmt"
	"jsson/internal/ast"
	"jsson/internal/token"
	"strings"
)

func (t *Transpiler) evalInterpolatedString(e *ast.InterpolatedString, ctx map[string]any) (any, error) {
	var result strings.Builder

	for _, part := range e.Parts {
		switch p := part.(type) {
		case ast.TextPart:
			result.WriteString(p.Value)
		case ast.ExprPart:
			expr := p.Expr

			if ident, ok := expr.(*ast.Identifier); ok {
				found := false

				if ctx != nil {
					_, found = ctx[ident.Value]
				}

				if !found {
					if e.Token.Type == token.TEMPLATESTR {
						result.WriteString("${")
						result.WriteString(ident.Value)
						result.WriteString("}")
					} else {
						result.WriteString("{")
						result.WriteString(ident.Value)
						result.WriteString("}")
					}

					continue
				}
			}

			val, err := t.evalExpression(expr, ctx)
			if err != nil {
				return nil, err
			}

			fmt.Fprintf(&result, "%v", val)
		}
	}

	return result.String(), nil
}
