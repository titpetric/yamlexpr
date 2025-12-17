package yamlexpr

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/stack"
)

// TestExpr_ProcessIfConditions tests if: directives.
func TestExpr_ProcessIfConditions(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "if-false-omits-key",
			input: map[string]any{
				"config": map[string]any{
					"name": "production",
					"debug": map[string]any{
						"if":      false,
						"enabled": true,
					},
					"version": "1.0",
				},
			},
			expected: map[string]any{
				"config": map[string]any{
					"name":    "production",
					"version": "1.0",
				},
			},
		},
		{
			name: "if-true-includes-key",
			input: map[string]any{
				"config": map[string]any{
					"name": "production",
					"debug": map[string]any{
						"if":      true,
						"enabled": true,
						"level":   "verbose",
					},
					"version": "1.0",
				},
			},
			expected: map[string]any{
				"config": map[string]any{
					"name": "production",
					"debug": map[string]any{
						"enabled": true,
						"level":   "verbose",
					},
					"version": "1.0",
				},
			},
		},
	}

	e := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs, err := e.Parse(Document(tt.input))
			require.NoError(t, err)
			require.Len(t, docs, 1)
			result := map[string]any(docs[0])
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestExpr_ProcessForLoops tests for: directives.
func TestExpr_ProcessForLoops(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "for-empty-array",
			input: map[string]any{
				"items": map[string]any{
					"for": []any{},
				},
				"servers": []any{
					map[string]any{"name": "main"},
				},
			},
			expected: map[string]any{
				"items": []any{},
				"servers": []any{
					map[string]any{"name": "main"},
				},
			},
		},
		{
			name: "for-single-item",
			input: map[string]any{
				"users": []any{
					map[string]any{
						"for":  []any{"alice"},
						"name": "${item}",
						"role": "admin",
					},
				},
				"metadata": map[string]any{
					"version": "1.0",
				},
			},
			expected: map[string]any{
				"users": []any{
					map[string]any{
						"name": "alice",
						"role": "admin",
					},
				},
				"metadata": map[string]any{
					"version": "1.0",
				},
			},
		},
		{
			name: "for-multiple-items",
			input: map[string]any{
				"users": []any{
					map[string]any{
						"for":    []any{"alice", "bob", "charlie"},
						"name":   "${item}",
						"active": true,
					},
				},
			},
			expected: map[string]any{
				"users": []any{
					map[string]any{
						"name":   "alice",
						"active": true,
					},
					map[string]any{
						"name":   "bob",
						"active": true,
					},
					map[string]any{
						"name":   "charlie",
						"active": true,
					},
				},
			},
		},
	}

	e := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputDoc, ok := tt.input.(map[string]any)
			require.True(t, ok, "input must be map[string]any")
			docs, err := e.Parse(Document(inputDoc))
			require.NoError(t, err)
			require.Len(t, docs, 1)
			result := map[string]any(docs[0])
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestExpr_ProcessForWithIf tests combined for: and if: directives.
func TestExpr_ProcessForWithIf(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "for-with-if-simple-boolean",
			input: map[string]any{
				"items": []any{
					map[string]any{
						"for": []any{
							map[string]any{"name": "api", "active": true},
							map[string]any{"name": "worker", "active": false},
							map[string]any{"name": "scheduler", "active": true},
						},
						"if":   "item.active",
						"name": "${item.name}",
					},
				},
			},
			expected: map[string]any{
				"items": []any{
					map[string]any{"name": "api"},
					map[string]any{"name": "scheduler"},
				},
			},
		},
	}

	e := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputDoc, ok := tt.input.(map[string]any)
			require.True(t, ok, "input must be map[string]any")
			docs, err := e.Parse(Document(inputDoc))
			require.NoError(t, err)
			require.Len(t, docs, 1)
			result := map[string]any(docs[0])
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestExpr_EvaluateCondition tests condition evaluation.
func TestExpr_EvaluateCondition(t *testing.T) {
	tests := []struct {
		name      string
		condition any
		st        *stack.Stack
		expected  bool
	}{
		{
			name:      "boolean-true",
			condition: true,
			st:        stack.New(),
			expected:  true,
		},
		{
			name:      "boolean-false",
			condition: false,
			st:        stack.New(),
			expected:  false,
		},
		{
			name:      "string-true",
			condition: "true",
			st:        stack.New(),
			expected:  true,
		},
		{
			name:      "string-false",
			condition: "false",
			st:        stack.New(),
			expected:  false,
		},
		{
			name:      "path-true",
			condition: "active",
			st:        stack.NewStack(map[string]any{"active": true}),
			expected:  true,
		},
		{
			name:      "path-false",
			condition: "active",
			st:        stack.NewStack(map[string]any{"active": false}),
			expected:  false,
		},
		{
			name:      "nested-path-true",
			condition: "item.active",
			st: func() *stack.Stack {
				st := stack.New()
				st.Push(map[string]any{"item": map[string]any{"active": true}})
				return st
			}(),
			expected: true,
		},
		{
			name:      "nested-path-false",
			condition: "item.active",
			st: func() *stack.Stack {
				st := stack.New()
				st.Push(map[string]any{"item": map[string]any{"active": false}})
				return st
			}(),
			expected: false,
		},
		{
			name:      "interpolation-true",
			condition: "${active}",
			st:        stack.NewStack(map[string]any{"active": true}),
			expected:  true,
		},
		{
			name:      "interpolation-false",
			condition: "${active}",
			st:        stack.NewStack(map[string]any{"active": false}),
			expected:  false,
		},
		{
			name:      "comparison-eq-true",
			condition: "${status} == 'active'",
			st:        stack.NewStack(map[string]any{"status": "active"}),
			expected:  true,
		},
		{
			name:      "comparison-eq-false",
			condition: "${status} == 'active'",
			st:        stack.NewStack(map[string]any{"status": "inactive"}),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluateConditionWithPath(tt.condition, tt.st, "")
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestProcessMapWithContext_EdgeCases tests map processing with edge cases.
func TestProcessMapWithContext_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		checkFn func(t *testing.T, docs []Document, err error)
	}{
		{
			name: "nested-maps-with-interpolation",
			input: map[string]any{
				"config": map[string]any{
					"db": map[string]any{
						"host": "${db_host}",
					},
				},
				"db_host": "localhost",
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				require.NotNil(t, docs[0]["config"])
			},
		},
		{
			name: "map-with-if-conditions",
			input: map[string]any{
				"settings": map[string]any{
					"debug": map[string]any{
						"if":    false,
						"level": "verbose",
					},
				},
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				// Verify if condition was processed
				require.NotNil(t, docs[0]["settings"])
			},
		},
		{
			name: "empty-nested-maps",
			input: map[string]any{
				"config": map[string]any{
					"empty": map[string]any{},
				},
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				require.NotNil(t, docs[0]["config"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(nil)
			docs, err := e.Parse(tt.input)
			tt.checkFn(t, docs, err)
		})
	}
}

// TestProcessSliceWithContext_EdgeCases tests slice processing with edge cases.
func TestProcessSliceWithContext_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		checkFn func(t *testing.T, docs []Document, err error)
	}{
		{
			name: "slice-with-interpolated-items",
			input: map[string]any{
				"hosts": []any{
					"${host1}",
					"${host2}",
				},
				"host1": "server1",
				"host2": "server2",
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				require.NotNil(t, docs[0]["hosts"])
			},
		},
		{
			name: "slice-with-nested-maps",
			input: map[string]any{
				"items": []any{
					map[string]any{
						"id":   1,
						"name": "${item_a}",
					},
				},
				"item_a": "First",
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				require.NotNil(t, docs[0]["items"])
			},
		},
		{
			name: "slice-with-mixed-types",
			input: map[string]any{
				"data": []any{
					"string",
					42,
					true,
				},
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				require.NotNil(t, docs[0]["data"])
			},
		},
		{
			name: "slice-with-if-filter",
			input: map[string]any{
				"items": []any{
					map[string]any{
						"if":    true,
						"value": "yes",
					},
				},
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				require.NotNil(t, docs[0]["items"])
			},
		},
		{
			name: "include-non-existent-file-produces-error",
			input: map[string]any{
				"config": map[string]any{
					"include": "nonexistent.yaml",
				},
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				// Should produce an error when trying to include non-existent file
				require.Error(t, err, "expected error for non-existent include file")
			},
		},
		{
			name: "empty-slice",
			input: map[string]any{
				"items": []any{},
			},
			checkFn: func(t *testing.T, docs []Document, err error) {
				require.NoError(t, err)
				require.Len(t, docs, 1)
				require.NotNil(t, docs[0]["items"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(nil)
			docs, err := e.Parse(tt.input)
			tt.checkFn(t, docs, err)
		})
	}
}
