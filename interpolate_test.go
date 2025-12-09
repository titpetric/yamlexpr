package yamlexpr

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/handlers"
	"github.com/titpetric/yamlexpr/stack"
)

// TestInterpolateStringHelper tests the handler's lenient interpolation function.
func TestInterpolateStringHelper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		stack map[string]any
		want  string
	}{
		{
			name:  "with-nil-stack",
			input: "hello world",
			stack: nil,
			want:  "hello world",
		},
		{
			name:  "with-stack",
			input: "hello ${name}",
			stack: map[string]any{"name": "alice"},
			want:  "hello alice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var st *stack.Stack
			if tt.stack != nil {
				st = stack.New(tt.stack)
			} else {
				st = stack.New(nil)
			}
			got, err := handlers.InterpolateString(tt.input, st)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestInterpolateStringWithContext tests interpolation with error handling from handlers.
func TestInterpolateStringWithContext(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		stack   map[string]any
		want    string
		wantErr bool
	}{
		{
			name:    "no-interpolation",
			input:   "hello world",
			stack:   map[string]any{},
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "single-variable",
			input:   "hello ${name}",
			stack:   map[string]any{"name": "world"},
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "missing-variable",
			input:   "hello ${missing}",
			stack:   map[string]any{},
			want:    "",
			wantErr: true, // Error when variable not found
		},
		{
			name:    "multiple-variables",
			input:   "${greeting} ${name}!",
			stack:   map[string]any{"greeting": "hello", "name": "world"},
			want:    "hello world!",
			wantErr: false,
		},
		{
			name:    "nested-path",
			input:   "user: ${user.name}",
			stack:   map[string]any{"user": map[string]any{"name": "alice"}},
			want:    "user: alice",
			wantErr: false,
		},
		{
			name:    "partially-missing",
			input:   "${found} and ${missing}",
			stack:   map[string]any{"found": "something"},
			want:    "",
			wantErr: true, // Error on first missing variable
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := stack.New(tt.stack)
			got, err := handlers.InterpolateStringWithContext(tt.input, st, "")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// TestInterpolateStringWithContext_PathInError tests that path is included in errors.
func TestInterpolateStringWithContext_PathInError(t *testing.T) {
	st := stack.New(map[string]any{})
	_, err := handlers.InterpolateStringWithContext("value: ${missing}", st, "config.item")
	require.Error(t, err)
	// Path should be included in error message
	require.Contains(t, err.Error(), "config.item")
}
