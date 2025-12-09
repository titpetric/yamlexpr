package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/handlers"
	"github.com/titpetric/yamlexpr/stack"
)

// TestContainsInterpolation tests the ContainsInterpolation function.
func TestContainsInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"no interpolation", "hello world", false},
		{"with interpolation", "hello ${name}", true},
		{"multiple interpolations", "${first} ${second}", true},
		{"only opening bracket", "hello ${name", false},
		{"only closing bracket", "hello name}", false},
		{"just brackets", "${}", true}, // Contains both ${ and }, technically
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.ContainsInterpolation(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestInterpolateString tests the InterpolateString function with lenient mode (from util.go).
func TestInterpolateString(t *testing.T) {
	st := stack.New(map[string]any{
		"name":  "world",
		"count": 42,
	})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"no interpolation", "hello world", "hello world"},
		{"single interpolation", "hello ${name}", "hello world"},
		{"multiple interpolations", "${name} has ${count} items", "world has 42 items"},
		{"undefined variable", "hello ${undefined}", "hello ${undefined}"}, // Unchanged in lenient mode
		{"partial match", "prefix ${name} suffix", "prefix world suffix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.InterpolateString(tt.input, st)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestInterpolateStringWithContext tests interpolation with strict error handling.
func TestInterpolateStringWithContext(t *testing.T) {
	st := stack.New(map[string]any{
		"name":  "world",
		"count": 42,
	})

	tests := []struct {
		name      string
		input     string
		path      string
		expected  string
		expectErr bool
	}{
		{"no interpolation", "hello world", "", "hello world", false},
		{"single interpolation", "hello ${name}", "", "hello world", false},
		{"with path context", "hello ${name}", "root.greeting", "hello world", false},
		{"undefined variable", "hello ${undefined}", "", "", true},
		{"multiple vars, one undefined", "${name} ${undefined}", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.InterpolateStringWithContext(tt.input, st, tt.path)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestInterpolateValue tests interpolation of various value types.
func TestInterpolateValue(t *testing.T) {
	st := stack.New(map[string]any{
		"name": "Alice",
	})

	tests := []struct {
		name      string
		input     any
		expected  any
		expectErr bool
	}{
		{"string without interpolation", "hello", "hello", false},
		{"string with interpolation", "hello ${name}", "hello Alice", false},
		{"integer unchanged", 42, 42, false},
		{"boolean unchanged", true, true, false},
		{"nil unchanged", nil, nil, false},
		{"slice unchanged", []any{1, 2}, []any{1, 2}, false},
		{"map unchanged", map[string]any{"key": "value"}, map[string]any{"key": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.InterpolateValue(tt.input, st, "")
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestNewInterpolationHandler tests that NewInterpolationHandler returns a valid handler.
func TestNewInterpolationHandler(t *testing.T) {
	handler := handlers.NewInterpolationHandler()
	require.NotNil(t, handler)
}
