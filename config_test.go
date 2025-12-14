package yamlexpr_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

// TestNew_StandardHandlers verifies that New with standard options registers handlers.
func TestNew_StandardHandlers(t *testing.T) {
	e := yamlexpr.New()

	// Test that directives ARE processed with standard handlers
	doc := map[string]any{
		"if":    true,
		"name":  "test",
		"items": []any{"a", "b"},
	}

	result, err := e.Process(doc, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)

	// if directive should be processed and removed
	require.NotContains(t, m, "if")
	require.Contains(t, m, "name")
	require.Equal(t, "test", m["name"])
}

func TestNew_DefaultSyntax(t *testing.T) {
	e := yamlexpr.New()
	require.NotNil(t, e)

	// Test with default syntax
	doc := map[string]any{
		"for":   "item in items",
		"name":  "${item}",
		"items": []any{"a", "b"},
	}

	result, err := e.Process(doc, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestNew_CustomIfSyntax(t *testing.T) {
	e := yamlexpr.New(yamlexpr.WithSyntax(yamlexpr.Syntax{If: "v-if"}))

	doc := map[string]any{
		"name": "test",
		"v-if": true,
	}

	result, err := e.Process(doc, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "test", m["name"])
	require.NotContains(t, m, "v-if")
}

func TestNew_CustomForSyntax(t *testing.T) {
	e := yamlexpr.New(yamlexpr.WithSyntax(yamlexpr.Syntax{For: "v-for"}))

	doc := map[string]any{
		"users": []any{"alice", "bob"},
	}

	template := map[string]any{
		"v-for": "user in users",
		"name":  "${user}",
	}

	result, err := e.Process(template, doc)
	require.NoError(t, err)

	// Result should be a slice with 2 items (one for each user)
	s, ok := result.([]any)
	require.True(t, ok)
	require.Equal(t, 2, len(s))

	// Check first user
	first, ok := s[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "alice", first["name"])

	// Check second user
	second, ok := s[1].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "bob", second["name"])
}

func TestNew_CustomIncludeSyntax(t *testing.T) {
	// Create a mock filesystem with test files
	fs := fstest.MapFS{
		"_custom-base.yaml": &fstest.MapFile{Data: []byte("env: production\nport: 8080\n")},
	}

	// Use custom include directive name
	e := yamlexpr.New(yamlexpr.WithFS(fs), yamlexpr.WithSyntax(yamlexpr.Syntax{Include: "v-include"}))

	// Create a YAML template that includes another file using custom directive
	template := map[string]any{
		"v-include": "_custom-base.yaml",
		"debug":     true,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "production", m["env"])
	require.Equal(t, 8080, m["port"])
	require.Equal(t, true, m["debug"])
	require.NotContains(t, m, "v-include") // Directive should be removed
}

func TestNew_AllCustomSyntax(t *testing.T) {
	e := yamlexpr.New(yamlexpr.WithSyntax(yamlexpr.Syntax{
		If:      "v-if",
		For:     "v-for",
		Include: "v-include",
	}))

	doc := map[string]any{
		"items": []any{true, true},
	}

	template := map[string]any{
		"v-for": "item in items",
		"name":  "test",
		"v-if":  true,
	}

	result, err := e.Process(template, doc)
	require.NoError(t, err)

	s, ok := result.([]any)
	require.True(t, ok)
	require.Equal(t, 2, len(s))

	// Both items should have passed the condition
	for _, item := range s {
		m, ok := item.(map[string]any)
		require.True(t, ok)
		require.Equal(t, "test", m["name"])
		require.NotContains(t, m, "v-if")
	}
}

func TestNew_CustomForWithIfSyntax(t *testing.T) {
	e := yamlexpr.New(yamlexpr.WithSyntax(yamlexpr.Syntax{For: "v-for", If: "v-if"}))

	doc := map[string]any{
		"numbers": []any{1, 2, 3, 4, 5},
	}

	template := map[string]any{
		"v-for": "num in numbers",
		"value": "${num}",
		"v-if":  "num > 2", // Only include numbers greater than 2
	}

	result, err := e.Process(template, doc)
	require.NoError(t, err)

	s, ok := result.([]any)
	require.True(t, ok)
	require.Equal(t, 3, len(s)) // Should have 3, 4, 5

	// Check values - variable references preserve native types, so ${num} returns int not string
	expected := []int{3, 4, 5}
	for i, item := range s {
		m, ok := item.(map[string]any)
		require.True(t, ok)
		require.Equal(t, expected[i], m["value"])
	}
}

func TestNew_PartialCustomSyntax(t *testing.T) {
	// Test that empty fields keep defaults
	e := yamlexpr.New(
		yamlexpr.WithSyntax(yamlexpr.Syntax{
			If:  "v-if",
			For: "v-for",
			// Include remains "include"
		}),
	)

	doc := map[string]any{
		"active": []any{true},
	}

	template := map[string]any{
		"v-for": "flag in active",
		"name":  "enabled",
		"v-if":  true,
	}

	result, err := e.Process(template, doc)
	require.NoError(t, err)

	s, ok := result.([]any)
	require.True(t, ok)
	require.Equal(t, 1, len(s))

	m, ok := s[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "enabled", m["name"])
}

func TestNew_IfConditionWithCustomSyntax(t *testing.T) {
	e := yamlexpr.New(yamlexpr.WithSyntax(yamlexpr.Syntax{If: "v-if"}))

	doc := map[string]any{
		"enabled": true,
	}

	template := map[string]any{
		"name": "test",
		"v-if": "${enabled}",
	}

	result, err := e.Process(template, doc)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "test", m["name"])
	require.NotContains(t, m, "v-if")
}

func TestNew_FalseIfWithCustomSyntax(t *testing.T) {
	e := yamlexpr.New(yamlexpr.WithSyntax(yamlexpr.Syntax{If: "v-if"}))

	template := map[string]any{
		"name": "test",
		"v-if": false,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	// When if condition is false, the entire block should be omitted (return nil)
	require.Nil(t, result)
}

func TestNew_IncludeWithCustomSyntax(t *testing.T) {
	// Create a mock filesystem with test files
	fs := fstest.MapFS{
		"users.yaml": &fstest.MapFile{Data: []byte("users:\n  - alice\n  - bob\n")},
	}

	// Use all custom directive names
	e := yamlexpr.New(
		yamlexpr.WithFS(fs),
		yamlexpr.WithSyntax(yamlexpr.Syntax{
			If:      "v-if",
			For:     "v-for",
			Include: "v-include",
		}),
	)

	// Create a template that includes a file, then iterates and filters
	template := map[string]any{
		"v-include": "users.yaml",
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)

	// Check the included data is there
	users, ok := m["users"].([]any)
	require.True(t, ok)
	require.Equal(t, 2, len(users))
	require.Equal(t, "alice", users[0])
	require.Equal(t, "bob", users[1])

	// Now verify we can iterate over the included data
	// First process the template to get the users data, then pass it as root vars
	doc := map[string]any{
		"v-include": "users.yaml",
	}
	result2, err := e.Process(doc, nil)
	require.NoError(t, err)

	// Now create a for loop template and pass the result as root vars
	forTemplate := map[string]any{
		"items": map[string]any{
			"v-for": "user in users",
			"name":  "${user}",
		},
	}

	result3, err := e.Process(forTemplate, result2.(map[string]any))
	require.NoError(t, err)

	m3, ok := result3.(map[string]any)
	require.True(t, ok)

	// Should have expanded items
	items, ok := m3["items"].([]any)
	require.True(t, ok)
	require.Equal(t, 2, len(items))

	item1, ok := items[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "alice", item1["name"])

	item2, ok := items[1].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "bob", item2["name"])
}
