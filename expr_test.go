package yamlexpr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
	"github.com/titpetric/yamlexpr/stack"
)

func TestNew(t *testing.T) {
	e := yamlexpr.New(nil)
	require.NotNil(t, e)
}

func TestExpr_ProcessWithStack(t *testing.T) {
	e := yamlexpr.New(nil)
	st := stack.New(map[string]any{"name": "John"})

	// Test with primitive value
	result, err := e.ProcessWithStack("hello", st)
	require.NoError(t, err)
	require.Equal(t, "hello", result)

	// Test with map
	doc := map[string]any{"key": "value"}
	result, err = e.ProcessWithStack(doc, st)
	require.NoError(t, err)
	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "value", m["key"])

	// Test with slice
	sliceDoc := []any{"a", "b", "c"}
	result, err = e.ProcessWithStack(sliceDoc, st)
	require.NoError(t, err)
	s, ok := result.([]any)
	require.True(t, ok)
	require.Equal(t, 3, len(s))
}

func TestExpr_Process(t *testing.T) {
	e := yamlexpr.New(nil)

	doc := map[string]any{
		"name":  "test",
		"items": []any{"a", "b"},
	}
	result, err := e.Process(doc)
	require.NoError(t, err)
	require.NotNil(t, result)
}
