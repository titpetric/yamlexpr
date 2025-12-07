package yamlexpr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

// TestParseForExpr tests the for loop expression parser.
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
			forExpr, err := yamlexpr.ParseForExpr(tt.input)

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
