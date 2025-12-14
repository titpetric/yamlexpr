package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"

	"github.com/titpetric/yamlexpr/model"
	"github.com/titpetric/yamlexpr/stack"
)

// EvaluateConditionWithPath evaluates an if condition with path context for error messages.
// Supports:
// - Boolean values: true/false
// - Interpolated expressions: "${item.active}"
// - Direct variable paths: item.active (converted to expressions via go-expr)
// - Complex expressions: item.status == 'active', item.count > 5, etc.
// Returns errors with variable context and path if referenced variables don't exist.
func EvaluateConditionWithPath(st *stack.Stack, condition any, path string) (bool, error) {
	switch v := condition.(type) {
	case bool:
		return v, nil
	case string:
		// Check for literal true/false first
		switch v {
		case "true", "1", "yes":
			return true, nil
		case "false", "0", "no", "":
			return false, nil
		}

		// Handle interpolated expressions like "${item.active}"
		if strings.Contains(v, "${") {
			str, err := InterpolateString(st, v)
			if err != nil {
				return false, err
			}
			// After interpolation, try to parse as boolean
			switch str {
			case "true", "1", "yes":
				return true, nil
			case "false", "0", "no", "":
				return false, nil
			default:
				// The interpolated result couldn't be parsed as boolean
				// If it's a string comparison (like "active == 'active'"), quote the left side
				v = QuoteUnquotedComparisons(str)
			}
		}

		// Use go-expr to evaluate the expression
		env := st.All()
		program, err := expr.Compile(v, expr.Env(env))
		if err != nil {
			pathCtx := ""
			if path != "" {
				pathCtx = fmt.Sprintf(" at %s", path)
			}
			return false, fmt.Errorf("error compiling expression '%s'%s: %w", v, pathCtx, err)
		}

		result, err := expr.Run(program, env)
		if err != nil {
			pathCtx := ""
			if path != "" {
				pathCtx = fmt.Sprintf(" at %s", path)
			}
			return false, fmt.Errorf("error evaluating expression '%s'%s: %w", v, pathCtx, err)
		}

		// Convert result to boolean
		return IsTruthy(result), nil

	case int, int8, int16, int32, int64:
		// Non-zero is true
		return v != 0, nil
	case float32, float64:
		// Non-zero is true
		return v != 0.0, nil
	case nil:
		// nil is always false
		return false, nil
	default:
		pathCtx := ""
		if path != "" {
			pathCtx = fmt.Sprintf(" at %s", path)
		}
		return false, fmt.Errorf("unsupported condition type: %T%s", condition, pathCtx)
	}
}

// QuoteUnquotedComparisons adds quotes around unquoted string literals in comparisons.
// For example: "active == 'active'" stays the same, but "active == test" becomes "'active' == 'test'"
// Numeric values and function calls are left unchanged.
func QuoteUnquotedComparisons(expr string) string {
	// Check for comparison operators
	operators := []string{"==", "!=", "<", ">", "<=", ">="}

	for _, op := range operators {
		if strings.Contains(expr, op) {
			parts := strings.Split(expr, op)
			if len(parts) == 2 {
				left := strings.TrimSpace(parts[0])
				right := strings.TrimSpace(parts[1])

				// Quote unquoted parts (but not if they're already quoted, numeric, have dots, or function calls)
				if !IsQuoted(left) && !strings.Contains(left, ".") && !strings.Contains(left, "(") && !isNumeric(left) {
					left = "'" + left + "'"
				}
				if !IsQuoted(right) && !strings.Contains(right, ".") && !strings.Contains(right, "(") && !isNumeric(right) {
					right = "'" + right + "'"
				}

				return left + " " + op + " " + right
			}
		}
	}

	return expr
}

// isNumeric checks if a string represents a numeric value.
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsQuoted checks if a string is already quoted
func IsQuoted(s string) bool {
	s = strings.TrimSpace(s)
	return (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}

// IsTruthy returns true for non-empty/non-zero values.
func IsTruthy(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int, int8, int16, int32, int64:
		return val != 0
	case uint, uint8, uint16, uint32, uint64:
		return val != 0
	case float32, float64:
		return val != 0.0
	case string:
		return val != ""
	case []any:
		return len(val) > 0
	case map[string]any:
		return len(val) > 0
	default:
		return v != nil
	}
}

// If creates a handler for the "if" directive.
// Implements conditional block inclusion.
func If(ifDirective string) DirectiveHandler {
	return func(ctx *model.Context, block map[string]any, value any) ([]any, bool, error) {
		ok, err := EvaluateConditionWithPath(ctx.Stack(), value, ctx.Path()+"."+ifDirective)
		if err != nil {
			return nil, false, err
		}
		if !ok {
			// Return nil if condition is false (omit the entire block)
			return nil, true, nil
		}
		// Continue processing with normal flow
		return nil, false, nil
	}
}
