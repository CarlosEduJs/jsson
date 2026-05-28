package transpiler

import (
	"jsson/internal/ast"
	"maps"
)

func (t *Transpiler) evalMapExpression(e *ast.MapExpression, ctx map[string]any) (any, error) {
	leftVal, err := t.evalExpression(e.Left, ctx)
	if err != nil {
		return nil, err
	}

	var items []any

	switch v := leftVal.(type) {
	case []any:
		items = v
	case rangeFlattener:
		items = v.Flatten()
	default:
		return nil, t.errfNode(e, "map target is not an array, it's a %T — gremlin is confused", leftVal)
	}

	var result []any

	for _, item := range items {
		newCtx := make(map[string]any, len(ctx))

		maps.Copy(newCtx, ctx)
		newCtx[e.Iterator.Value] = item

		mappedVal, err := t.evalExpression(e.Body, newCtx)
		if err != nil {
			return nil, err
		}

		result = append(result, mappedVal)
	}

	return result, nil
}
