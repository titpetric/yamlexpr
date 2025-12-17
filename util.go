package yamlexpr

import (
	"fmt"
	"regexp"
	"strings"
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

// MapMatchesSpec checks if a map contains all key:value pairs from a specification map.
// Used for matrix include/exclude matching and other spec-based filtering.
// Returns true only if every key in spec exists in the map with an equal value.
func MapMatchesSpec(m map[string]any, spec map[string]any) bool {
	for specKey, specVal := range spec {
		mapVal, exists := m[specKey]
		if !exists {
			return false
		}
		if !valuesEqual(mapVal, specVal) {
			return false
		}
	}
	return true
}

// ValuesEqual checks if two values are equal, handling primitives and type coercion.
// Used for comparing values in matrix specs where YAML may parse numbers as float64 or int.
func ValuesEqual(a, b any) bool {
	return valuesEqual(a, b)
}

// valuesEqual is the internal implementation of value comparison.
func valuesEqual(a, b any) bool {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case int:
		// Try both int and float64 (YAML may parse as float)
		switch bv := b.(type) {
		case int:
			return av == bv
		case float64:
			return float64(av) == bv
		}
		return false
	case float64:
		// Handle both float64 and int
		switch bv := b.(type) {
		case float64:
			return av == bv
		case int:
			return av == float64(bv)
		}
		return false
	case bool:
		bv, ok := b.(bool)
		return ok && av == bv
	default:
		// Fallback: compare string representations
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}
}
