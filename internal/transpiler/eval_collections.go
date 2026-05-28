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

		if rf, ok := val.(rangeFlattener); ok {
			arr = append(arr, rf.Flatten()...)
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
		evaluatedRow, err := t.evalTemplateRow(row, ctx)
		if err != nil {
			return nil, err
		}

		hasArrays, minLength := t.templateRowHasArrays(evaluatedRow)

		if hasArrays && minLength > 0 {
			items := t.expandArrays(evaluatedRow, keys, isImplicitTemplate, minLength)

			for _, item := range items {
				mapped, err := t.applyTemplateMap(e.Map, item, ctx)
				if err != nil {
					return nil, err
				}

				result = append(result, mapped)
			}
		} else {
			item := t.buildRowItem(evaluatedRow, keys, isImplicitTemplate)

			mapped, err := t.applyTemplateMap(e.Map, item, ctx)
			if err != nil {
				return nil, err
			}

			result = append(result, mapped)
		}
	}

	return result, nil
}

func (t *Transpiler) evalTemplateRow(row []ast.Expression, ctx map[string]any) ([]any, error) {
	evaluatedRow := make([]any, len(row))

	for i, expr := range row {
		val, err := t.evalExpression(expr, ctx)
		if err != nil {
			return nil, err
		}

		if rf, ok := val.(rangeFlattener); ok {
			evaluatedRow[i] = rf.Flatten()
		} else {
			evaluatedRow[i] = val
		}
	}

	return evaluatedRow, nil
}

func (t *Transpiler) templateRowHasArrays(evaluatedRow []any) (hasArrays bool, minLen int) {
	hasArrays = false
	minLen = -1

	for _, val := range evaluatedRow {
		arr, ok := val.([]any)
		if !ok {
			continue
		}

		isObjectArray := false

		if len(arr) > 0 {
			if _, isMap := arr[0].(map[string]any); isMap {
				isObjectArray = true
			}
		}

		if !isObjectArray {
			hasArrays = true

			if minLen == -1 || len(arr) < minLen {
				minLen = len(arr)
			}
		}
	}

	return
}

func (t *Transpiler) expandArrays(evaluatedRow []any, keys []string, isImplicitTemplate bool, count int) []any {
	items := make([]any, 0, count)

	for idx := range count {
		items = append(items, t.buildRowItemAt(evaluatedRow, keys, isImplicitTemplate, idx))
	}

	return items
}

func (t *Transpiler) buildRowItem(evaluatedRow []any, keys []string, isImplicitTemplate bool) any {
	if isImplicitTemplate {
		return evaluatedRow[0]
	}

	rowObj := make(map[string]any)

	for i, val := range evaluatedRow {
		if i >= len(keys) {
			break
		}

		rowObj[keys[i]] = val
	}

	return rowObj
}

func (t *Transpiler) buildRowItemAt(evaluatedRow []any, keys []string, isImplicitTemplate bool, idx int) any {
	if isImplicitTemplate {
		if arr, ok := evaluatedRow[0].([]any); ok {
			return arr[idx]
		}

		return evaluatedRow[0]
	}

	rowObj := make(map[string]any)

	for i, val := range evaluatedRow {
		if i >= len(keys) {
			break
		}

		if arr, ok := val.([]any); ok {
			rowObj[keys[i]] = arr[idx]
		} else {
			rowObj[keys[i]] = val
		}
	}

	return rowObj
}

func (t *Transpiler) applyTemplateMap(mapClause *ast.MapClause, itemValue any, ctx map[string]any) (any, error) {
	if mapClause == nil {
		return itemValue, nil
	}

	mapCtx := make(map[string]any, len(ctx))
	maps.Copy(mapCtx, ctx)

	mapCtx[mapClause.Param.Value] = itemValue

	return t.evalExpression(mapClause.Body, mapCtx)
}
