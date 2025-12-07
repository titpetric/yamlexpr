package yamlexpr

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/titpetric/yamlexpr/stack"
)

// parseForExpr parses a for loop expression string.
// Supports:
//   - "item in items" - iterates over items, binding each to 'item'
//   - "item in item.subitems" - iterate over nested path
//   - "(idx, item) in items" - iterates over items, binding index to 'idx' and item to 'item'
//   - "(key, value) in items" - for map iteration
//   - "_" can be used to omit a variable
func parseForExpr(expr string) (*ForLoopExpr, error) {
	expr = strings.TrimSpace(expr)

	// Pattern 1: (var1, var2, ...) in source
	// Source can be a simple name or a dotted path (e.g., item.subitem.array)
	tuplePattern := regexp.MustCompile(`^\((.*?)\)\s+in\s+([\w.]+)$`)
	if matches := tuplePattern.FindStringSubmatch(expr); matches != nil {
		varsPart := strings.TrimSpace(matches[1])
		source := strings.TrimSpace(matches[2])

		// Split variables by comma
		vars := strings.Split(varsPart, ",")
		for i := range vars {
			vars[i] = strings.TrimSpace(vars[i])
			if vars[i] == "" {
				return nil, fmt.Errorf("empty variable name in for expression: %q", expr)
			}
			// Validate variable names (allow "_" for omitted variables)
			if vars[i] != "_" && !isValidVarName(vars[i]) {
				return nil, fmt.Errorf("invalid variable name %q in for expression: %q", vars[i], expr)
			}
		}

		return &ForLoopExpr{
			Variables: vars,
			Source:    source,
		}, nil
	}

	// Pattern 2: var in source (single variable)
	// Source can be a simple name or a dotted path (e.g., item.subitem.array)
	simplePattern := regexp.MustCompile(`^(\w+)\s+in\s+([\w.]+)$`)
	if matches := simplePattern.FindStringSubmatch(expr); matches != nil {
		varName := strings.TrimSpace(matches[1])
		source := strings.TrimSpace(matches[2])

		if !isValidVarName(varName) {
			return nil, fmt.Errorf("invalid variable name %q in for expression: %q", varName, expr)
		}

		return &ForLoopExpr{
			Variables: []string{varName},
			Source:    source,
		}, nil
	}

	return nil, fmt.Errorf("invalid for expression syntax: %q (expected 'var in source' or '(var1, var2) in source', source can be a path like 'item.subitem')", expr)
}

// isValidVarName checks if a string is a valid variable name or "_".
func isValidVarName(name string) bool {
	if name == "_" {
		return true
	}
	// Must start with letter or underscore, followed by letters, digits, or underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z_]\w*$`, name)
	return matched
}

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

// interpolateStringHelper is a helper to interpolate a single string without a full Expr instance.
func interpolateStringHelper(s string, st *stack.Stack) string {
	if st == nil {
		st = stack.New(nil)
	}
	return interpolateString(s, st)
}

// containsInterpolation checks if a string contains ${...} patterns.
func containsInterpolation(s string) bool {
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
