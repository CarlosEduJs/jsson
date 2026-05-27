package transpiler

import (
	"crypto/rand"
	"fmt"
	"jsson/internal/ast"
	"math/big"
	"time"
)

// generateValidatorValue generates a value for a validator expression.
func (t *Transpiler) generateValidatorValue(v *ast.ValidatorExpression) (any, error) {
	switch v.Type {
	case "uuid":
		return generateUUID(), nil
	case "email":
		return fmt.Sprintf("user%d@example.com", time.Now().UnixNano()%10000), nil
	case "url":
		return "https://example.com", nil
	case "ipv4":
		return fmt.Sprintf("192.168.%d.%d", time.Now().UnixNano()%256, (time.Now().UnixNano()/256)%256), nil
	case "ipv6":
		return "2001:0db8:85a3:0000:0000:8a2e:0370:7334", nil
	case "filepath":
		return "/path/to/file.txt", nil
	case "date":
		return time.Now().Format("2006-01-02"), nil
	case "datetime":
		return time.Now().Format(time.RFC3339), nil
	case "regex":
		if v.Pattern != "" {
			return "matched-value", nil
		}

		return "sample-text", nil
	case "int":
		return generateInt(v.Args), nil
	case "float":
		return generateFloat(v.Args), nil
	case "bool":
		return generateBool(), nil
	default:
		return nil, t.errfNode(v, "unknown validator type: %s", v.Type)
	}
}

// generateUUID generates a random UUID v4.
func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// generateInt generates a random integer between min and max.
func generateInt(args []any) int64 {
	var minVal, maxVal int64 = 0, 100 // defaults

	if len(args) >= 2 {
		if v, ok := args[0].(int64); ok {
			minVal = v
		}

		if v, ok := args[1].(int64); ok {
			maxVal = v
		}
	}

	if minVal >= maxVal {
		return minVal
	}

	n, err := rand.Int(rand.Reader, big.NewInt(maxVal-minVal+1))
	if err != nil {
		return minVal
	}

	return minVal + n.Int64()
}

// generateFloat generates a random float between min and max.
func generateFloat(args []any) float64 {
	minVal, maxVal := 0.0, 1.0 // defaults

	if len(args) >= 2 {
		// Handle both int64 and float64
		switch v := args[0].(type) {
		case int64:
			minVal = float64(v)
		case float64:
			minVal = v
		}

		switch v := args[1].(type) {
		case int64:
			maxVal = float64(v)
		case float64:
			maxVal = v
		}
	}

	if minVal >= maxVal {
		return minVal
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1<<53))
	if err != nil {
		return minVal
	}

	return minVal + float64(n.Int64())/(1<<53)*(maxVal-minVal)
}

// generateBool generates a random boolean.
func generateBool() bool {
	n, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		return false
	}

	return n.Int64() == 1
}
