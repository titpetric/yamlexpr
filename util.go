package yamlexpr

import (
	"fmt"
	"regexp"
	"strings"
)

// parseForExpr parses a for loop expression string into a ForLoopExpr.
// Supports multiple syntax patterns:
//   - "item in items" - iterates over items, binding each to 'item'
//   - "item in item.subitems" - iterate over nested path
//   - "(idx, item) in items" - iterates over items, binding index to 'idx' and item to 'item'
//   - "(key, value) in items" - for map iteration
//   - "_" can be used to omit a variable
//
// Returns an error if the expression syntax is invalid.
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

// isValidVarName checks if a string is a valid variable name or "_" (omit marker).
// Valid names start with a letter or underscore, followed by letters, digits, or underscores.
func isValidVarName(name string) bool {
	if name == "_" {
		return true
	}
	// Must start with letter or underscore, followed by letters, digits, or underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z_]\w*$`, name)
	return matched
}
