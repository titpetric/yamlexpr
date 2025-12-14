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
	result, err := e.ProcessWithStack(st, "hello")
	require.NoError(t, err)
	require.Equal(t, "hello", result)

	// Test with map
	doc := map[string]any{"key": "value"}
	result, err = e.ProcessWithStack(st, doc)
	require.NoError(t, err)
	require.Equal(t, doc, result[0])

	// Test with slice
	sliceDoc := []any{"a", "b", "c"}
	result, err = e.ProcessWithStack(st, sliceDoc)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, sliceDoc, result[0])
}

func TestExpr_Process(t *testing.T) {
	e := yamlexpr.New()

	doc := map[string]any{
		"name":  "${user.name}",
		"items": []any{"a", "b"},
	}

	want := map[string]any{
		"name":  "John",
		"items": []any{"a", "b"},
	}

	got, err := e.Process(doc, map[string]any{"user": map[string]any{"name": "John"}})

	require.NoError(t, err)
	require.Equal(t, want, got)
}
