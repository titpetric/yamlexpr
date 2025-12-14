package yamlexpr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
	"github.com/titpetric/yamlexpr/stack"
)

func TestExpr_ProcessWithStack(t *testing.T) {
	e := yamlexpr.New()
	st := stack.NewStack(map[string]any{"name": "John"})

	// Test with primitive value
	t.Run("primitive value", func(t *testing.T) {
		result, err := e.ProcessWithStack(st, "hello")
		require.Nil(t, result)
		require.Error(t, err)
	})

	t.Run("with map", func(t *testing.T) {
		doc := map[string]any{"key": "value"}

		want := []any{
			map[string]any{
				"key": "value",
			},
		}

		result, err := e.ProcessWithStack(st, doc)
		require.NoError(t, err)
		require.Equal(t, want, result)
	})

	t.Run("with slice", func(t *testing.T) {
		// Test with slice
		sliceDoc := []any{"a", "b", "c"}

		want := []any{
			[]any{"a", "b", "c"},
		}

		result, err := e.ProcessWithStack(st, sliceDoc)
		require.NoError(t, err)
		require.Equal(t, want, result)
	})
}

func TestExpr_Process(t *testing.T) {
	e := yamlexpr.New()

	doc := map[string]any{
		"name":  "${user.name}",
		"items": []any{"a", "b"},
	}

	want := []any{
		map[string]any{
			"name":  "John",
			"items": []any{"a", "b"},
		},
	}

	got, err := e.Process(doc, map[string]any{"user": map[string]any{"name": "John"}})

	require.NoError(t, err)
	require.Equal(t, want, got)
}
