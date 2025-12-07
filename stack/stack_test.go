package stack_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/stack"
)

func TestNew(t *testing.T) {
	s := stack.New(nil)
	require.NotNil(t, s)

	data := map[string]any{"key": "value"}
	s = stack.New(data)
	v, ok := s.Lookup("key")
	require.True(t, ok)
	require.Equal(t, "value", v)
}

func TestStack_Set(t *testing.T) {
	s := stack.New(nil)
	s.Set("key", "value")
	v, ok := s.Lookup("key")
	require.True(t, ok)
	require.Equal(t, "value", v)
}

func TestStack_Lookup(t *testing.T) {
	s := stack.New(map[string]any{"root": "value"})
	v, ok := s.Lookup("root")
	require.True(t, ok)
	require.Equal(t, "value", v)

	// Non-existent key
	v, ok = s.Lookup("missing")
	require.False(t, ok)
	require.Nil(t, v)
}

func TestStack_LookupStackOrder(t *testing.T) {
	s := stack.New(map[string]any{"key": "root"})
	s.Push(map[string]any{"key": "pushed"})
	v, ok := s.Lookup("key")
	require.True(t, ok)
	require.Equal(t, "pushed", v)
}

func TestStack_Resolve(t *testing.T) {
	s := stack.New(map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
		},
		"items": []any{
			map[string]any{"id": 1, "title": "First"},
			map[string]any{"id": 2, "title": "Second"},
		},
	})

	// Simple key
	v, ok := s.Resolve("user")
	require.True(t, ok)
	require.NotNil(t, v)

	// Nested map
	v, ok = s.Resolve("user.name")
	require.True(t, ok)
	require.Equal(t, "John", v)

	// Array index
	v, ok = s.Resolve("items[0]")
	require.True(t, ok)
	require.NotNil(t, v)

	// Array index nested
	v, ok = s.Resolve("items[0].title")
	require.True(t, ok)
	require.Equal(t, "First", v)

	// Non-existent path
	v3, ok3 := s.Resolve("user.missing")
	require.False(t, ok3)
	require.Nil(t, v3)
}

func TestStack_GetString(t *testing.T) {
	s := stack.New(map[string]any{
		"str":  "hello",
		"num":  42,
		"flt":  3.14,
		"bool": true,
	})

	v, ok := s.GetString("str")
	require.True(t, ok)
	require.Equal(t, "hello", v)

	v, ok = s.GetString("num")
	require.True(t, ok)
	require.Equal(t, "42", v)

	v, ok = s.GetString("flt")
	require.True(t, ok)
	require.Equal(t, "3.14", v)

	v, ok = s.GetString("bool")
	require.True(t, ok)
	require.Equal(t, "true", v)

	v2, ok2 := s.GetString("missing")
	require.False(t, ok2)
	require.Empty(t, v2)
}

func TestStack_GetInt(t *testing.T) {
	s := stack.New(map[string]any{
		"int": 42,
		"str": "100",
	})

	v, ok := s.GetInt("int")
	require.True(t, ok)
	require.Equal(t, 42, v)

	v, ok = s.GetInt("str")
	require.True(t, ok)
	require.Equal(t, 100, v)

	v2, ok2 := s.GetInt("missing")
	require.False(t, ok2)
	require.Equal(t, 0, v2)
}

func TestStack_GetSlice(t *testing.T) {
	s := stack.New(map[string]any{
		"items": []any{"a", "b", "c"},
	})

	v, ok := s.GetSlice("items")
	require.True(t, ok)
	require.Equal(t, 3, len(v))
	require.Equal(t, "a", v[0])
}

func TestStack_GetMap(t *testing.T) {
	s := stack.New(map[string]any{
		"user": map[string]any{"name": "John", "age": 30},
	})

	v, ok := s.GetMap("user")
	require.True(t, ok)
	require.Equal(t, "John", v["name"])
	require.Equal(t, 30, v["age"])
}

func TestStack_Push_Pop(t *testing.T) {
	s := stack.New(map[string]any{"root": "value"})
	s.Push(map[string]any{"local": "scoped"})

	// Find in pushed scope
	v, ok := s.Lookup("local")
	require.True(t, ok)
	require.Equal(t, "scoped", v)

	// Root still accessible
	v, ok = s.Lookup("root")
	require.True(t, ok)
	require.Equal(t, "value", v)

	// Pop and verify local scope is gone, root remains
	s.Pop()
	v2, ok2 := s.Lookup("local")
	require.False(t, ok2)
	require.Nil(t, v2)

	v3, ok3 := s.Lookup("root")
	require.True(t, ok3)
	require.Equal(t, "value", v3)
}

func TestStack_All(t *testing.T) {
	s := stack.New(map[string]any{"root": "value"})
	s.Push(map[string]any{"local": "scoped"})

	all := s.All()
	require.Equal(t, "value", all["root"])
	require.Equal(t, "scoped", all["local"])
}

func TestStack_Copy(t *testing.T) {
	s := stack.New(map[string]any{"key": "value"})
	s.Push(map[string]any{"local": "scoped"})

	copy := s.Copy()
	all := copy.All()
	require.Equal(t, 2, len(all))
	require.Equal(t, "value", all["key"])
	require.Equal(t, "scoped", all["local"])
}

func TestStack_ForEach(t *testing.T) {
	s := stack.New(map[string]any{
		"items": []any{"a", "b", "c"},
	})

	var results []any
	err := s.ForEach("items", func(idx int, val any) error {
		results = append(results, val)
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 3, len(results))
	require.Equal(t, "a", results[0])
	require.Equal(t, "b", results[1])
	require.Equal(t, "c", results[2])
}
