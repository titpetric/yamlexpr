package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"

	"github.com/titpetric/yamlexpr/stack"
)

// interpolationPattern matches ${...} syntax in strings.
var interpolationPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// ContainsInterpolation checks if a string contains interpolation patterns (${...}).
func ContainsInterpolation(s string) bool {
	return strings.Contains(s, "${") && strings.Contains(s, "}")
}

// InterpolateStringWithContext replaces ${...} placeholders with values.
// Supports both simple variable references (${varname}) and expressions (${item * 2}).
// It returns an error if a referenced variable doesn't exist or cannot be converted to string.
// The path parameter is used for error context information.
func InterpolateStringWithContext(st *stack.Stack, s string, path string) (string, error) {
	var result strings.Builder
	lastIdx := 0

	for _, match := range interpolationPattern.FindAllStringSubmatchIndex(s, -1) {
		// match has: [fullMatchStart, fullMatchEnd, group1Start, group1End, ...]
		matchStart, matchEnd := match[0], match[1]
		var captureStart, captureEnd int

		// Check if we have capture groups
		if len(match) >= 4 {
			captureStart, captureEnd = match[2], match[3]
		} else {
			// No capture group, shouldn't happen with our pattern
			continue
		}

		// Add the string before this match
		result.WriteString(s[lastIdx:matchStart])

		// Extract expression from capture group
		exprStr := s[captureStart:captureEnd]

		// Try to evaluate as an expression first using expr-lang
		// This allows both simple variables and complex expressions like "item * 2"
		env := st.All()
		program, err := expr.Compile(exprStr, expr.Env(env))
		if err == nil {
			// Expression compiled successfully, evaluate it
			val, err := expr.Run(program, env)
			if err != nil {
				pathCtx := ""
				if path != "" {
					pathCtx = fmt.Sprintf(" at %s", path)
				}
				return "", fmt.Errorf("error evaluating expression '%s'%s: %w", exprStr, pathCtx, err)
			}

			// Convert result to string, handling null values as YAML literal
			if val == nil {
				result.WriteString("null")
			} else {
				str := fmt.Sprintf("%v", val)
				result.WriteString(str)
			}
		} else {
			// Expression compilation failed, try simple variable lookup (backwards compat)
			val, ok := st.Resolve(exprStr)
			if !ok || val == nil {
				pathCtx := ""
				if path != "" {
					pathCtx = fmt.Sprintf(" at %s", path)
				}
				return "", fmt.Errorf("undefined variable '%s'%s", exprStr, pathCtx)
			}

			// Convert to string
			str, ok := st.GetString(exprStr)
			if !ok {
				pathCtx := ""
				if path != "" {
					pathCtx = fmt.Sprintf(" at %s", path)
				}
				return "", fmt.Errorf("cannot convert variable '%s' to string%s", exprStr, pathCtx)
			}

			result.WriteString(str)
		}

		lastIdx = matchEnd
	}

	// Add remaining string
	result.WriteString(s[lastIdx:])
	return result.String(), nil
}

// InterpolateValue interpolates a value (typically a string) using the given stack and path.
// For strings with interpolation, evaluates expressions and returns native types when possible.
// Non-string values are returned unchanged. This is a strict version that errors on undefined variables.
func InterpolateValue(value any, st *stack.Stack, path string) (any, error) {
	switch v := value.(type) {
	case string:
		if ContainsInterpolation(v) {
			// Try to interpolate as a single expression first
			// This allows ${item * 2} to return an integer instead of a string
			// Also preserves null values (${xcode} with xcode=null returns null, not string "null")
			if isSingleInterpolation(v) {
				exprStr := extractSingleExpression(v)
				env := st.All()
				program, err := expr.Compile(exprStr, expr.Env(env))
				if err == nil {
					// Expression compiled successfully, return the native type
					result, err := expr.Run(program, env)
					if err != nil {
						pathCtx := ""
						if path != "" {
							pathCtx = fmt.Sprintf(" at %s", path)
						}
						return nil, fmt.Errorf("error evaluating expression '%s'%s: %w", exprStr, pathCtx, err)
					}
					// Return native type, including nil for null values
					return result, nil
				}
			}
			// Fall back to string interpolation
			return InterpolateStringWithContext(st, v, path)
		}
		return v, nil
	default:
		// Return non-string values unchanged
		return value, nil
	}
}

// InterpolateValueWithContext is like InterpolateValue but works with ExprContext.
// For single interpolations, preserves the native type of the value.
// For multiple interpolations or mixed text, returns a string.
func InterpolateValueWithContext(s string, st *stack.Stack, path string) (any, error) {
	if !ContainsInterpolation(s) {
		return s, nil
	}

	// Try to interpolate as a single expression/variable first
	// Also preserves null values (${xcode} with xcode=null returns null, not string "null")
	if isSingleInterpolation(s) {
		exprStr := extractSingleExpression(s)
		env := st.All()
		program, err := expr.Compile(exprStr, expr.Env(env))
		if err == nil {
			// Expression compiled successfully, return the native type
			result, err := expr.Run(program, env)
			if err != nil {
				pathCtx := ""
				if path != "" {
					pathCtx = fmt.Sprintf(" at %s", path)
				}
				return nil, fmt.Errorf("error evaluating expression '%s'%s: %w", exprStr, pathCtx, err)
			}
			// Return native type, including nil for null values
			return result, nil
		}
	}

	// Fall back to string interpolation
	return InterpolateStringWithContext(st, s, path)
}

// isSingleInterpolation checks if a string contains exactly one ${...} expression
func isSingleInterpolation(s string) bool {
	matches := interpolationPattern.FindAllStringIndex(s, -1)
	if len(matches) != 1 {
		return false
	}
	// Check if the expression spans the entire string (no text before or after)
	match := matches[0]
	return match[0] == 0 && match[1] == len(s)
}

// extractSingleExpression extracts the expression from a ${...} pattern
func extractSingleExpression(s string) string {
	matches := interpolationPattern.FindAllStringSubmatchIndex(s, -1)
	if len(matches) == 0 || len(matches[0]) < 4 {
		return ""
	}
	captureStart, captureEnd := matches[0][2], matches[0][3]
	return s[captureStart:captureEnd]
}

// InterpolateStringPermissive replaces ${varname} placeholders with stack values.
// If a variable is undefined or nil, the entire interpolated string returns nil.
// This is useful for optional matrix dimensions that may not be set.
func InterpolateStringPermissive(s string, st *stack.Stack) (any, error) {
	var result strings.Builder
	lastIdx := 0

	for _, match := range interpolationPattern.FindAllStringSubmatchIndex(s, -1) {
		// match has: [fullMatchStart, fullMatchEnd, group1Start, group1End, ...]
		matchStart, matchEnd := match[0], match[1]
		var captureStart, captureEnd int

		// Check if we have capture groups
		if len(match) >= 4 {
			captureStart, captureEnd = match[2], match[3]
		} else {
			// No capture group, shouldn't happen with our pattern
			continue
		}

		// Add the string before this match
		result.WriteString(s[lastIdx:matchStart])

		// Extract variable name from capture group
		varName := s[captureStart:captureEnd]

		// Resolve from stack - if undefined or nil, return nil for entire interpolation
		val, ok := st.Resolve(varName)
		if !ok || val == nil {
			return nil, nil
		}

		// Convert to string - if cannot convert, return nil
		str, ok := st.GetString(varName)
		if !ok {
			return nil, nil
		}

		result.WriteString(str)
		lastIdx = matchEnd
	}

	// Add remaining string
	result.WriteString(s[lastIdx:])
	return result.String(), nil
}

// InterpolateValuePermissive interpolates a value with permissive handling of undefined variables.
// Strings with interpolation return nil if any referenced variable is undefined.
// Non-string values are returned unchanged.
func InterpolateValuePermissive(value any, st *stack.Stack) (any, error) {
	switch v := value.(type) {
	case string:
		if ContainsInterpolation(v) {
			return InterpolateStringPermissive(v, st)
		}
		return v, nil
	default:
		// Return non-string values unchanged
		return value, nil
	}
}
