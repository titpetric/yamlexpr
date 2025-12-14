package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/handlers"
	"github.com/titpetric/yamlexpr/stack"
)

func TestEvaluateConditionWithPath_Boolean(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"true boolean", true, true},
		{"false boolean", false, false},
		{"string true", "true", true},
		{"string false", "false", false},
		{"string yes", "yes", true},
		{"string no", "no", false},
		{"empty string", "", false},
		{"int zero", 0, false},
		{"int one", 1, true},
		{"int negative", -1, true},
		{"float zero", 0.0, false},
		{"float positive", 1.5, true},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.EvaluateConditionWithPath(tt.value, stack.New(), "")
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEvaluateConditionWithPath_StringLiterals(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true", "true", true},
		{"false", "false", false},
		{"1", "1", true},
		{"0", "0", false},
		{"yes", "yes", true},
		{"no", "no", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.EvaluateConditionWithPath(tt.value, stack.New(), "")
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEvaluateConditionWithPath_InvalidType(t *testing.T) {
	_, err := handlers.EvaluateConditionWithPath([]string{"foo"}, stack.New(), "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported condition type")
}

func TestEvaluateConditionWithPath_InvalidString(t *testing.T) {
	_, err := handlers.EvaluateConditionWithPath("maybe", stack.New(), "test.if")
	require.Error(t, err)
	require.Contains(t, err.Error(), "at test.if")
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"true bool", true, true},
		{"false bool", false, false},
		{"int zero", int(0), false},
		{"int nonzero", int(5), true},
		{"empty string", "", false},
		{"nonempty string", "hello", true},
		{"empty slice", []any{}, false},
		{"nonempty slice", []any{1}, true},
		{"empty map", map[string]any{}, false},
		{"nonempty map", map[string]any{"a": 1}, true},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.IsTruthy(tt.value)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestQuoteUnquotedComparisons(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"already quoted",
			"'active' == 'active'",
			"'active' == 'active'",
		},
		{
			"unquoted identifiers",
			"active == test",
			"'active' == 'test'",
		},
		{
			"mixed quoted",
			"'active' == test",
			"'active' == 'test'",
		},
		{
			"with dots (variable path)",
			"item.status == 'active'",
			"item.status == 'active'",
		},
		{
			"with function call",
			"len(items) == 5",
			"len(items) == 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.QuoteUnquotedComparisons(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestIsQuoted(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"double quoted", `"hello"`, true},
		{"single quoted", `'hello'`, true},
		{"unquoted", `hello`, false},
		{"partial double", `"hello`, false},
		{"partial single", `'hello`, false},
		{"empty quoted double", `""`, true},
		{"empty quoted single", `''`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.IsQuoted(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestIfHandler(t *testing.T) {
	handler := handlers.IfHandler("if")
	require.NotNil(t, handler)
}
