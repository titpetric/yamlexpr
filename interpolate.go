package yamlexpr

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/titpetric/yamlexpr/stack"
)

// interpolationPattern matches ${...} syntax in strings.
var interpolationPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// interpolateString replaces ${varname} placeholders with stack values.
// Non-string values are returned as-is.
func interpolateString(s string, st *stack.Stack) string {
	return interpolationPattern.ReplaceAllStringFunc(s, func(match string) string {
		// Extract variable name from ${varname}
		expr := match[2 : len(match)-1] // Remove ${ and }

		// Resolve the value from stack
		val, ok := st.Resolve(expr)
		if !ok || val == nil {
			// Return placeholder unchanged if not found
			return match
		}

		// Convert value to string
		str, ok := st.GetString(expr)
		if !ok {
			return match
		}
		return str
	})
}

// InterpolateString is a helper to interpolate a single string without a full Expr instance.
func InterpolateString(s string, st *stack.Stack) string {
	if st == nil {
		st = stack.New(nil)
	}
	return interpolateString(s, st)
}

// ContainsInterpolation checks if a string contains ${...} patterns.
func ContainsInterpolation(s string) bool {
	return strings.Contains(s, "${") && strings.Contains(s, "}")
}

// interpolateStringWithContext replaces ${varname} placeholders with stack values,
// returning an error if a referenced variable doesn't exist. Path is used for error context.
func interpolateStringWithContext(s string, st *stack.Stack, path string) (string, error) {
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

		// Resolve from stack
		val, ok := st.Resolve(varName)
		if !ok || val == nil {
			pathCtx := ""
			if path != "" {
				pathCtx = fmt.Sprintf(" at %s", path)
			}
			return "", fmt.Errorf("undefined variable '%s'%s", varName, pathCtx)
		}

		// Convert to string
		str, ok := st.GetString(varName)
		if !ok {
			pathCtx := ""
			if path != "" {
				pathCtx = fmt.Sprintf(" at %s", path)
			}
			return "", fmt.Errorf("cannot convert variable '%s' to string%s", varName, pathCtx)
		}

		result.WriteString(str)
		lastIdx = matchEnd
	}

	// Add remaining string
	result.WriteString(s[lastIdx:])
	return result.String(), nil
}
