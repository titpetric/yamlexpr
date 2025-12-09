package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/titpetric/yamlexpr/model"
	"github.com/titpetric/yamlexpr/stack"
)

// interpolationPattern matches ${...} syntax in strings.
var interpolationPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// ContainsInterpolation checks if a string contains interpolation patterns (${...}).
func ContainsInterpolation(s string) bool {
	return strings.Contains(s, "${") && strings.Contains(s, "}")
}

// InterpolateStringWithContext replaces ${varname} placeholders with stack values.
// It returns an error if a referenced variable doesn't exist or cannot be converted to string.
// The path parameter is used for error context information.
func InterpolateStringWithContext(s string, st *stack.Stack, path string) (string, error) {
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

// InterpolateValue interpolates a value (typically a string) using the given stack and path.
// Non-string values are returned unchanged. This is a strict version that errors on undefined variables.
func InterpolateValue(value any, st *stack.Stack, path string) (any, error) {
	switch v := value.(type) {
	case string:
		if ContainsInterpolation(v) {
			return InterpolateStringWithContext(v, st, path)
		}
		return v, nil
	default:
		// Return non-string values unchanged
		return value, nil
	}
}

// NewInterpolationHandler creates a handler for string interpolation.
// This handler processes ${variable} syntax in string values and replaces them
// with their resolved values from the stack.
//
// Note: Interpolation is automatically applied to all string values during
// document processing, but this handler allows explicit control over when
// interpolation occurs or can be used as a standalone directive.
//
// Priority: 50 (runs after include/for/if, before custom handlers)
func NewInterpolationHandler() *InterpolationHandlerImpl {
	return &InterpolationHandlerImpl{}
}

// InterpolationHandlerImpl implements the interpolation handler
type InterpolationHandlerImpl struct{}

// InterpolationHandlerBuiltin creates a handler for explicit interpolation control.
// This allows interpolation to be triggered as a directive in YAML.
func InterpolationHandlerBuiltin(interpolateDirective string) DirectiveHandler {
	return func(ctx *model.Context, block map[string]any, value any) (any, bool, error) {
		// Interpolation handler doesn't consume directives when used as a directive.
		// It's typically applied automatically during value processing.
		// If used explicitly as a directive, it returns nil to continue normal processing.
		return nil, false, nil
	}
}
