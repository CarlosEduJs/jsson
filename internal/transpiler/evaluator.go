package transpiler

import (
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"maps"
)

func (t *Transpiler) evalExpression(expr ast.Expression, ctx map[string]any) (any, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return e.Value, nil
	case *ast.FloatLiteral:
		return e.Value, nil
	case *ast.BooleanLiteral:
		return e.Value, nil
	case *ast.StringLiteral:
		return e.Value, nil
	case *ast.ValidatorExpression:
		return t.generateValidatorValue(e)
	case *ast.PresetReference:
		return t.evalPresetReference(e, ctx)
	case *ast.MapExpression:
		return t.evalMapExpression(e, ctx)
	case *ast.InterpolatedString:
		return t.evalInterpolatedString(e, ctx)
	case *ast.Identifier:
		return t.evalIdentifier(e, ctx)
	case *ast.ObjectLiteral:
		return t.evalObjectLiteral(e, ctx)
	case *ast.ArrayLiteral:
		return t.evalArrayLiteral(e, ctx)
	case *ast.RangeExpression:
		return t.evalRangeExpression(e, ctx)
	case *ast.ArrayTemplate:
		return t.evalArrayTemplate(e, ctx)
	case *ast.BinaryExpression:
		return t.evalBinaryExpression(e, ctx)
	case *ast.ConditionalExpression:
		return t.evalConditionalExpression(e, ctx)
	case *ast.MemberExpression:
		return t.evalMemberExpression(e, ctx)
	default:
		return nil, t.errfNode(expr, "unknown expression type: %T", expr)
	}
}

func (t *Transpiler) evalPresetReference(e *ast.PresetReference, ctx map[string]any) (any, error) {
	presetName := e.Name.Value

	presetBody, exists := t.presetTable[presetName]
	if !exists {
		return nil, t.errfNode(e, "preset %q not found — define it with @preset %q { ... }", presetName, presetName)
	}

	baseVal, err := t.evalExpression(presetBody, ctx)
	if err != nil {
		return nil, err
	}

	baseObj, ok := baseVal.(map[string]any)
	if !ok {
		return nil, t.errfNode(e, "preset %q did not evaluate to an object", presetName)
	}

	if e.Overrides != nil {
		overridesVal, err := t.evalExpression(e.Overrides, ctx)
		if err != nil {
			return nil, err
		}

		overridesObj, ok := overridesVal.(map[string]any)
		if !ok {
			return nil, t.errfNode(e, "preset overrides must be an object")
		}

		result := make(map[string]any, len(baseObj)+len(overridesObj))
		maps.Copy(result, baseObj)
		maps.Copy(result, overridesObj)

		return result, nil
	}

	result := make(map[string]any, len(baseObj))
	maps.Copy(result, baseObj)

	return result, nil
}

func (t *Transpiler) evalIdentifier(e *ast.Identifier, ctx map[string]any) (any, error) {
	if ctx != nil {
		if val, ok := ctx[e.Value]; ok {
			return val, nil
		}
	}

	if val, ok := t.symbolTable[e.Value]; ok {
		return val, nil
	}

	return nil, t.errfNodeMsg(e, ie.UndefinedVariable(e.Value))
}

func (t *Transpiler) evalBinaryExpression(e *ast.BinaryExpression, ctx map[string]any) (any, error) {
	left, err := t.evalExpression(e.Left, ctx)
	if err != nil {
		return nil, err
	}

	right, err := t.evalExpression(e.Right, ctx)
	if err != nil {
		return nil, err
	}

	return t.evalBinary(left, e.Operator, right)
}

func (t *Transpiler) evalConditionalExpression(e *ast.ConditionalExpression, ctx map[string]any) (any, error) {
	condition, err := t.evalExpression(e.Condition, ctx)
	if err != nil {
		return nil, err
	}

	isTruthy := t.isTruthy(condition)

	if isTruthy {
		return t.evalExpression(e.Consequence, ctx)
	}

	return t.evalExpression(e.Alternative, ctx)
}

func (t *Transpiler) evalMemberExpression(e *ast.MemberExpression, ctx map[string]any) (any, error) {
	leftVal, err := t.evalExpression(e.Left, ctx)
	if err != nil {
		return nil, err
	}

	if obj, ok := leftVal.(map[string]any); ok {
		if val, ok := obj[e.Property.Value]; ok {
			return val, nil
		}

		return nil, t.errfNode(e, "property %q not found — gremlin searched everywhere", e.Property.Value)
	}

	return nil, t.errfNodeMsg(e, ie.NotAnObject())
}
