package yamlexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestForLoopExpr_ParseSimple tests parsing simple for loop expressions.
func TestForLoopExpr_ParseSimple(t *testing.T) {
	tests := []struct {
		name       string
		expr       string
		wantVars   []string
		wantSource string
		wantErr    bool
	}{
		{
			name:       "single-variable",
			expr:       "item in items",
			wantVars:   []string{"item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "single-variable-nested-path",
			expr:       "item in config.items",
			wantVars:   []string{"item"},
			wantSource: "config.items",
			wantErr:    false,
		},
		{
			name:       "two-variables",
			expr:       "(idx, item) in items",
			wantVars:   []string{"idx", "item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "two-variables-with-omit",
			expr:       "(_, item) in items",
			wantVars:   []string{"_", "item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "whitespace-handling",
			expr:       "  item  in  items  ",
			wantVars:   []string{"item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:     "invalid-syntax",
			expr:     "item items",
			wantVars: nil,
			wantErr:  true,
		},
		{
			name:     "empty-expression",
			expr:     "",
			wantVars: nil,
			wantErr:  true,
		},
		{
			name:     "invalid-variable-name",
			expr:     "123item in items",
			wantVars: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseForExpr(tt.expr)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantVars, result.Variables)
			require.Equal(t, tt.wantSource, result.Source)
		})
	}
}

// TestForLoopExpr_ParseMultipleVariables tests parsing expressions with multiple variables.
func TestForLoopExpr_ParseMultipleVariables(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		wantVars  []string
		wantError bool
	}{
		{
			name:     "key-value",
			expr:     "(key, value) in items",
			wantVars: []string{"key", "value"},
		},
		{
			name:     "index-item",
			expr:     "(index, item) in collection",
			wantVars: []string{"index", "item"},
		},
		{
			name:     "omit-first",
			expr:     "(_, second) in items",
			wantVars: []string{"_", "second"},
		},
		{
			name:     "omit-second",
			expr:     "(first, _) in items",
			wantVars: []string{"first", "_"},
		},
		{
			name:      "empty-variable",
			expr:      "(key, , value) in items",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseForExpr(tt.expr)
			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantVars, result.Variables)
		})
	}
}

// TestIsValidVarName tests variable name validation.
func TestIsValidVarName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{
			name:      "simple-name",
			input:     "item",
			wantValid: true,
		},
		{
			name:      "underscore",
			input:     "_",
			wantValid: true,
		},
		{
			name:      "leading-underscore",
			input:     "_item",
			wantValid: true,
		},
		{
			name:      "uppercase",
			input:     "Item",
			wantValid: true,
		},
		{
			name:      "with-digits",
			input:     "item123",
			wantValid: true,
		},
		{
			name:      "all-digits",
			input:     "123",
			wantValid: false,
		},
		{
			name:      "leading-digit",
			input:     "1item",
			wantValid: false,
		},
		{
			name:      "hyphen",
			input:     "item-name",
			wantValid: false,
		},
		{
			name:      "space",
			input:     "item name",
			wantValid: false,
		},
		{
			name:      "empty",
			input:     "",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidVarName(tt.input)
			require.Equal(t, tt.wantValid, got)
		})
	}
}
