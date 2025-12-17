package yamlexpr

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/interpolation"
	"github.com/titpetric/yamlexpr/stack"
)

func interpolateStringHelper(s string, st *stack.Stack) string {
	v, _ := interpolation.InterpolateString(st, s)
	return v
}

func containsInterpolation(s string) bool {
	return interpolation.ContainsInterpolation(s)
}

// TestInterpolateString tests the internal interpolateStringHelper function.
func TestInterpolateString(t *testing.T) {
	st := stack.NewStack(map[string]any{
		"name": "John",
		"age":  30,
		"user": map[string]any{
			"email": "john@example.com",
		},
	})

	// Simple interpolation
	result := interpolateStringHelper("Hello ${name}", st)
	require.Equal(t, "Hello John", result)

	// Multiple interpolations
	result = interpolateStringHelper("${name} is ${age} years old", st)
	require.Equal(t, "John is 30 years old", result)

	// Nested path
	result = interpolateStringHelper("Email: ${user.email}", st)
	require.Equal(t, "Email: john@example.com", result)

	// Missing variable
	result = interpolateStringHelper("Hello ${missing}", st)
	require.Equal(t, "Hello ${missing}", result)

	// No interpolation
	result = interpolateStringHelper("No variables here", st)
	require.Equal(t, "No variables here", result)
}

// TestContainsInterpolation tests the internal containsInterpolation function.
func TestContainsInterpolation(t *testing.T) {
	require.True(t, containsInterpolation("Hello ${name}"))
	require.False(t, containsInterpolation("Hello name"))
	require.False(t, containsInterpolation("${unclosed"))
	require.False(t, containsInterpolation("${"))
}

// TestParseForExpr tests the internal parseForExpr function.
func TestParseForExpr(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantVars   []string
		wantSource string
		wantErr    bool
	}{
		{
			name:       "single variable",
			input:      "item in items",
			wantVars:   []string{"item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "single variable with whitespace",
			input:      "  item  in  items  ",
			wantVars:   []string{"item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "two variables with underscores",
			input:      "(_, item) in items",
			wantVars:   []string{"_", "item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "two variables index and item",
			input:      "(idx, item) in items",
			wantVars:   []string{"idx", "item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "two variables key and value",
			input:      "(key, value) in config",
			wantVars:   []string{"key", "value"},
			wantSource: "config",
			wantErr:    false,
		},
		{
			name:       "first underscore",
			input:      "(_, value) in data",
			wantVars:   []string{"_", "value"},
			wantSource: "data",
			wantErr:    false,
		},
		{
			name:       "both underscores",
			input:      "(_, _) in items",
			wantVars:   []string{"_", "_"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:    "invalid: no 'in'",
			input:   "item items",
			wantErr: true,
		},
		{
			name:    "invalid: empty variable name",
			input:   "(item, ) in items",
			wantErr: true,
		},
		{
			name:    "invalid: invalid variable name",
			input:   "(123item, value) in items",
			wantErr: true,
		},
		{
			name:    "invalid: trailing comma",
			input:   "item, in items",
			wantErr: true,
		},
		{
			name:       "three variables",
			input:      "(a, b, c) in items",
			wantVars:   []string{"a", "b", "c"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "variable with underscores",
			input:      "(first_var, second_var) in my_items",
			wantVars:   []string{"first_var", "second_var"},
			wantSource: "my_items",
			wantErr:    false,
		},
		{
			name:       "dotted path source",
			input:      "item in item.subitems",
			wantVars:   []string{"item"},
			wantSource: "item.subitems",
			wantErr:    false,
		},
		{
			name:       "nested dotted path",
			input:      "dept in item.departments.values",
			wantVars:   []string{"dept"},
			wantSource: "item.departments.values",
			wantErr:    false,
		},
		{
			name:       "tuple with dotted path",
			input:      "(idx, item) in config.items",
			wantVars:   []string{"idx", "item"},
			wantSource: "config.items",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forExpr, err := parseForExpr(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantVars, forExpr.Variables)
			require.Equal(t, tt.wantSource, forExpr.Source)
		})
	}
}

// TestMapMatchesSpec tests the MapMatchesSpec function.
func TestMapMatchesSpec(t *testing.T) {
	tests := []struct {
		name    string
		m       map[string]any
		spec    map[string]any
		matches bool
	}{
		{
			name: "exact match",
			m: map[string]any{
				"os":   "linux",
				"arch": "x86_64",
			},
			spec: map[string]any{
				"os":   "linux",
				"arch": "x86_64",
			},
			matches: true,
		},
		{
			name: "partial spec match",
			m: map[string]any{
				"os":      "linux",
				"arch":    "x86_64",
				"version": "14",
			},
			spec: map[string]any{
				"os": "linux",
			},
			matches: true,
		},
		{
			name: "spec mismatch",
			m: map[string]any{
				"os":   "linux",
				"arch": "x86_64",
			},
			spec: map[string]any{
				"os": "windows",
			},
			matches: false,
		},
		{
			name: "missing key in map",
			m: map[string]any{
				"os": "linux",
			},
			spec: map[string]any{
				"os":   "linux",
				"arch": "x86_64",
			},
			matches: false,
		},
		{
			name: "type coercion int to float",
			m: map[string]any{
				"os":      "linux",
				"version": float64(14),
			},
			spec: map[string]any{
				"version": 14,
			},
			matches: true,
		},
		{
			name: "empty spec always matches",
			m: map[string]any{
				"os": "linux",
			},
			spec:    map[string]any{},
			matches: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapMatchesSpec(tt.m, tt.spec)
			require.Equal(t, tt.matches, result)
		})
	}
}

// TestValuesEqual tests the ValuesEqual function.
func TestValuesEqual(t *testing.T) {
	tests := []struct {
		name  string
		a     any
		b     any
		equal bool
	}{
		{
			name:  "string equality",
			a:     "linux",
			b:     "linux",
			equal: true,
		},
		{
			name:  "string inequality",
			a:     "linux",
			b:     "windows",
			equal: false,
		},
		{
			name:  "int equality",
			a:     14,
			b:     14,
			equal: true,
		},
		{
			name:  "int to float coercion",
			a:     14,
			b:     float64(14),
			equal: true,
		},
		{
			name:  "float equality",
			a:     14.0,
			b:     14.0,
			equal: true,
		},
		{
			name:  "float to int coercion",
			a:     float64(14),
			b:     14,
			equal: true,
		},
		{
			name:  "bool equality",
			a:     true,
			b:     true,
			equal: true,
		},
		{
			name:  "bool inequality",
			a:     true,
			b:     false,
			equal: false,
		},
		{
			name:  "different types",
			a:     14,
			b:     "14",
			equal: false,
		},
		{
			name:  "nil vs nil",
			a:     nil,
			b:     nil,
			equal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValuesEqual(tt.a, tt.b)
			require.Equal(t, tt.equal, result)
		})
	}
}
