package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/handlers"
)

// TestBuildScope tests the BuildScope helper function.
func TestBuildScope(t *testing.T) {
	tests := []struct {
		name     string
		varNames []string
		idx      int
		item     any
		want     map[string]any
	}{
		{
			name:     "single-variable-item",
			varNames: []string{"item"},
			idx:      0,
			item:     "alice",
			want:     map[string]any{"item": "alice"},
		},
		{
			name:     "single-variable-at-index-2",
			varNames: []string{"item"},
			idx:      2,
			item:     "bob",
			want:     map[string]any{"item": "bob"},
		},
		{
			name:     "two-variables-index-and-item",
			varNames: []string{"idx", "item"},
			idx:      5,
			item:     "charlie",
			want:     map[string]any{"idx": 5, "item": "charlie"},
		},
		{
			name:     "omit-index",
			varNames: []string{"_", "item"},
			idx:      0,
			item:     "alice",
			want:     map[string]any{"item": "alice"},
		},
		{
			name:     "omit-item",
			varNames: []string{"idx", "_"},
			idx:      10,
			item:     "bob",
			want:     map[string]any{"idx": 10},
		},
		{
			name:     "omit-both",
			varNames: []string{"_", "_"},
			idx:      0,
			item:     "charlie",
			want:     map[string]any{},
		},
		{
			name:     "complex-item",
			varNames: []string{"item"},
			idx:      0,
			item:     map[string]any{"name": "alice", "active": true},
			want:     map[string]any{"item": map[string]any{"name": "alice", "active": true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlers.BuildScope(tt.varNames, tt.idx, tt.item)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestNewForHandler tests that NewForHandler returns a valid handler.
func TestNewForHandler(t *testing.T) {
	handler := handlers.NewForHandler()
	require.NotNil(t, handler)
	require.IsType(t, &handlers.ForHandlerImpl{}, handler)
}

// TestParseForExpr tests the ParseForExpr function.
func TestParseForExpr(t *testing.T) {
	tests := []struct {
		name       string
		expr       string
		wantVars   []string
		wantSource string
		wantErr    bool
	}{
		{
			name:       "single variable",
			expr:       "item in items",
			wantVars:   []string{"item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "two variables",
			expr:       "(idx, item) in items",
			wantVars:   []string{"idx", "item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "nested source path",
			expr:       "item in parent.children",
			wantVars:   []string{"item"},
			wantSource: "parent.children",
			wantErr:    false,
		},
		{
			name:       "omit index",
			expr:       "(_, item) in items",
			wantVars:   []string{"_", "item"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:       "omit item",
			expr:       "(idx, _) in items",
			wantVars:   []string{"idx", "_"},
			wantSource: "items",
			wantErr:    false,
		},
		{
			name:    "no 'in' keyword",
			expr:    "item items",
			wantErr: true,
		},
		{
			name:    "missing closing paren",
			expr:    "(idx, item in items",
			wantErr: true,
		},
		{
			name:    "trailing comma",
			expr:    "(idx, item,) in items",
			wantErr: true,
		},
		{
			name:    "empty variable name",
			expr:    "(idx,, item) in items",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.ParseForExpr(tt.expr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.wantVars, result.Variables)
				require.Equal(t, tt.wantSource, result.Source)
			}
		})
	}
}

// TestForHandler tests the for handler builtin.
func TestForHandler(t *testing.T) {
	handler := handlers.ForHandler(nil, "for")
	require.NotNil(t, handler)
}
