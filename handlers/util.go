package handlers

import (
	"regexp"
	"strings"

	"github.com/titpetric/yamlexpr/stack"
)

// InterpolateString replaces ${varname} placeholders with stack values.
// Non-string values are returned as-is. This is a simplified version
// that doesn't error on missing variables.
func InterpolateString(st *stack.Stack, s string) (string, error) {
	if st == nil {
		st = stack.New()
	}

	interpolationPattern := regexp.MustCompile(`\$\{([^}]+)\}`)
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
			// Return placeholder unchanged if not found (lenient mode)
			result.WriteString(s[matchStart:matchEnd])
			lastIdx = matchEnd
			continue
		}

		// Convert to string
		str, ok := st.GetString(varName)
		if !ok {
			// Return placeholder unchanged if conversion fails
			result.WriteString(s[matchStart:matchEnd])
			lastIdx = matchEnd
			continue
		}

		result.WriteString(str)
		lastIdx = matchEnd
	}

	// Add remaining string
	result.WriteString(s[lastIdx:])
	return result.String(), nil
}
