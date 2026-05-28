package transpiler

import (
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
)

// evalStringRange handles ranges of strings with numeric suffixes (e.g., IP addresses)
// Example: "192.168.1.100".."192.168.1.109" generates ["192.168.1.100", "192.168.1.101", ...].
func (t *Transpiler) evalStringRange(start, end string, stepV *int64, node ast.Node) (any, error) {
	// Find the numeric suffix in both strings
	// We'll look for the last sequence of digits
	var (
		startPrefix, endPrefix string
		startNum, endNum       int64
		foundStart, foundEnd   bool
	)

	// Extract numeric suffix from start

	for i := len(start) - 1; i >= 0; i-- {
		if start[i] < '0' || start[i] > '9' {
			// Found non-digit, extract number after this position
			if i < len(start)-1 {
				startPrefix = start[:i+1]

				numStr := start[i+1:]
				if n, err := fmt.Sscanf(numStr, "%d", &startNum); n == 1 && err == nil {
					foundStart = true
				}
			}

			break
		}

		if i == 0 {
			// Entire string is a number
			startPrefix = ""

			if n, err := fmt.Sscanf(start, "%d", &startNum); n == 1 && err == nil {
				foundStart = true
			}

			break
		}
	}

	// Extract numeric suffix from end
	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < '0' || end[i] > '9' {
			if i < len(end)-1 {
				endPrefix = end[:i+1]

				numStr := end[i+1:]
				if n, err := fmt.Sscanf(numStr, "%d", &endNum); n == 1 && err == nil {
					foundEnd = true
				}
			}

			break
		}

		if i == 0 {
			endPrefix = ""

			if n, err := fmt.Sscanf(end, "%d", &endNum); n == 1 && err == nil {
				foundEnd = true
			}

			break
		}
	}

	if !foundStart || !foundEnd {
		return nil, t.errfNode(node, "string range requires numeric suffix in both start and end (e.g., \"192.168.1.100\"..\"192.168.1.109\")")
	}

	if startPrefix != endPrefix {
		return nil, t.errfNode(node, "string range prefixes must match (start: %q, end: %q)", startPrefix, endPrefix)
	}

	// Determine step
	var step int64

	if stepV != nil {
		step = *stepV

		if step == 0 {
			return nil, t.errfNode(node, "step cannot be zero")
		}
	} else {
		if startNum > endNum {
			step = -1
		} else {
			step = 1
		}
	}

	if step == 0 {
		return nil, t.errfNode(node, "step cannot be zero")
	}

	// Calculate number of digits in original (for zero-padding)
	startStr := start[len(startPrefix):]
	padding := len(startStr)

	// Generate range
	res := make([]string, 0)

	checkCancelled := t.ctx != nil

	if step > 0 {
		for i := startNum; i <= endNum; i += step {
			if checkCancelled {
				select {
				case <-t.ctx.Done():
					return nil, t.ctx.Err()
				default:
				}
			}

			if padding > 1 && startStr[0] == '0' {
				res = append(res, fmt.Sprintf("%s%0*d", startPrefix, padding, i))
			} else {
				res = append(res, fmt.Sprintf("%s%d", startPrefix, i))
			}
		}
	} else {
		for i := startNum; i >= endNum; i += step {
			if checkCancelled {
				select {
				case <-t.ctx.Done():
					return nil, t.ctx.Err()
				default:
				}
			}

			if padding > 1 && startStr[0] == '0' {
				res = append(res, fmt.Sprintf("%s%0*d", startPrefix, padding, i))
			} else {
				res = append(res, fmt.Sprintf("%s%d", startPrefix, i))
			}
		}
	}

	return RangeResult[string]{Values: res}, nil
}

const maxRangeElements = 10_000_000

// evalIntegerRange evaluates an integer range expression.
func (t *Transpiler) evalIntegerRange(sInt, eInt int64, stepV *int64, node ast.Node) (any, error) {
	var step int64

	if stepV != nil {
		step = *stepV

		if step == 0 {
			return nil, t.errfNodeMsg(node, ie.StepCannotBeZero())
		}
	} else {
		if sInt > eInt {
			step = -1
		} else {
			step = 1
		}
	}

	if step == 0 {
		return nil, t.errfNodeMsg(node, ie.StepCannotBeZero())
	}

	// Prevent OOM: calculate approximate count before allocating
	var count int64

	if step > 0 && eInt >= sInt {
		count = ((eInt - sInt) / step) + 1
	} else if step < 0 && sInt >= eInt {
		count = ((sInt - eInt) / (-step)) + 1
	}

	if count > maxRangeElements {
		return nil, t.errfNode(node, "range with %d elements exceeds maximum of %d — that's too much data for me", count, maxRangeElements)
	}

	res := make([]int64, 0, count)

	checkCancelled := t.ctx != nil

	if step > 0 {
		for i := sInt; i <= eInt; i += step {
			if checkCancelled && i%1000 == 0 {
				select {
				case <-t.ctx.Done():
					return nil, t.ctx.Err()
				default:
				}
			}

			res = append(res, i)
		}
	} else {
		for i := sInt; i >= eInt; i += step {
			if checkCancelled && i%1000 == 0 {
				select {
				case <-t.ctx.Done():
					return nil, t.ctx.Err()
				default:
				}
			}

			res = append(res, i)
		}
	}

	return RangeResult[int64]{Values: res}, nil
}

// evalRangeExpression evaluates a range expression.
func (t *Transpiler) evalRangeExpression(e *ast.RangeExpression, ctx map[string]any) (any, error) {
	startV, err := t.evalExpression(e.Start, ctx)
	if err != nil {
		return nil, err
	}

	endV, err := t.evalExpression(e.End, ctx)
	if err != nil {
		return nil, err
	}

	var stepV *int64
	if e.Step != nil {
		sv, err := t.evalExpression(e.Step, ctx)
		if err != nil {
			return nil, err
		}

		var s int64

		switch v := sv.(type) {
		case int64:
			s = v
		case float64:
			s = int64(v)
		default:
			return nil, t.errfNode(e, "step must be a number, got %T", sv)
		}

		stepV = &s
	}

	if startStr, ok1 := startV.(string); ok1 {
		if endStr, ok2 := endV.(string); ok2 {
			return t.evalStringRange(startStr, endStr, stepV, e)
		}
	}

	sInt, ok1 := startV.(int64)

	eInt, ok2 := endV.(int64)
	if !ok1 || !ok2 {
		return nil, t.errfNodeMsg(e, ie.RangeBoundsNotIntegers(startV, endV))
	}

	return t.evalIntegerRange(sInt, eInt, stepV, e)
}
