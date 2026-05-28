package transpiler

import (
	"errors"
	"jsson/internal/ast"
	"maps"
)

func (t *Transpiler) evalObjectLiteral(e *ast.ObjectLiteral, ctx map[string]any) (any, error) {
	obj := make(map[string]any)

	localCtx := make(map[string]any, len(ctx))
	maps.Copy(localCtx, ctx)

	for _, decl := range e.Declarations {
		val, err := t.evalExpression(decl.Value, localCtx)
		if err != nil {
			return nil, err
		}

		localCtx[decl.Name.Value] = val
	}

	for _, key := range e.Keys {
		valExpr := e.Properties[key]
		if valExpr == nil {
			continue
		}

		val, err := t.evalExpression(valExpr, localCtx)
		if err != nil {
			return nil, err
		}

		obj[key] = val
	}

	return obj, nil
}

func (t *Transpiler) evalArrayLiteral(e *ast.ArrayLiteral, ctx map[string]any) (any, error) {
	arr := make([]any, 0, len(e.Elements))

	for _, el := range e.Elements {
		val, err := t.evalExpression(el, ctx)
		if err != nil {
			return nil, err
		}

		if rr, ok := val.(RangeResult); ok {
			arr = append(arr, rr.Values...)
		} else {
			arr = append(arr, val)
		}
	}

	return arr, nil
}

func (t *Transpiler) evalArrayTemplate(e *ast.ArrayTemplate, ctx map[string]any) (any, error) {
	if e.Template == nil {
		return nil, errors.New("array template has nil template")
	}

	result := make([]any, 0, len(e.Rows))
	keys := e.Template.Keys

	isImplicitTemplate := e.Map != nil && len(keys) == 1 && keys[0] == e.Map.Param.Value

	for _, row := range e.Rows {
		evaluatedRow := make([]any, len(row))

		for i, expr := range row {
			val, err := t.evalExpression(expr, ctx)
			if err != nil {
				return nil, err
			}

			if rr, ok := val.(RangeResult); ok {
				evaluatedRow[i] = rr.Values
			} else {
				evaluatedRow[i] = val
			}
		}

		hasArrays := false
		minArrayLength := -1

		for _, val := range evaluatedRow {
			if arr, ok := val.([]any); ok {
				isObjectArray := false

				if len(arr) > 0 {
					if _, isMap := arr[0].(map[string]any); isMap {
						isObjectArray = true
					}
				}

				if !isObjectArray {
					hasArrays = true

					if minArrayLength == -1 || len(arr) < minArrayLength {
						minArrayLength = len(arr)
					}
				}
			}
		}

		if hasArrays && minArrayLength > 0 {
			for idx := range minArrayLength {
				var itemValue any

				if isImplicitTemplate {
					if arr, ok := evaluatedRow[0].([]any); ok {
						itemValue = arr[idx]
					} else {
						itemValue = evaluatedRow[0]
					}
				} else {
					rowObj := make(map[string]any)

					for i, val := range evaluatedRow {
						if i >= len(keys) {
							break
						}

						key := keys[i]

						if arr, ok := val.([]any); ok {
							rowObj[key] = arr[idx]
						} else {
							rowObj[key] = val
						}
					}

					itemValue = rowObj
				}

				if e.Map != nil {
					mapCtx := make(map[string]any, len(ctx))
					maps.Copy(mapCtx, ctx)

					mapCtx[e.Map.Param.Value] = itemValue

					mappedVal, err := t.evalExpression(e.Map.Body, mapCtx)
					if err != nil {
						return nil, err
					}

					result = append(result, mappedVal)
				} else {
					result = append(result, itemValue)
				}
			}
		} else {
			var itemValue any

			if isImplicitTemplate {
				itemValue = evaluatedRow[0]
			} else {
				rowObj := make(map[string]any)

				for i, val := range evaluatedRow {
					if i >= len(keys) {
						break
					}

					key := keys[i]
					rowObj[key] = val
				}

				itemValue = rowObj
			}

			if e.Map != nil {
				mapCtx := make(map[string]any, len(ctx))
				maps.Copy(mapCtx, ctx)

				mapCtx[e.Map.Param.Value] = itemValue

				mappedVal, err := t.evalExpression(e.Map.Body, mapCtx)
				if err != nil {
					return nil, err
				}

				result = append(result, mappedVal)
			} else {
				result = append(result, itemValue)
			}
		}
	}

	return result, nil
}
