package yamlexpr_test

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

// TestIncludeHandler_SingleFile tests including a single file.
func TestIncludeHandler_SingleFile(t *testing.T) {
	fs := fstest.MapFS{
		"base.yaml": &fstest.MapFile{Data: []byte("env: production\nport: 8080\n")},
	}

	e := yamlexpr.New(yamlexpr.WithFS(fs))

	template := map[string]any{
		"include": "base.yaml",
		"debug":   true,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "production", m["env"])
	require.Equal(t, 8080, m["port"])
	require.Equal(t, true, m["debug"])
	require.NotContains(t, m, "include") // Directive removed
}

// TestIncludeHandler_MultipleFiles tests including multiple files.
func TestIncludeHandler_MultipleFiles(t *testing.T) {
	fs := fstest.MapFS{
		"config.yaml":  &fstest.MapFile{Data: []byte("env: production\n")},
		"secrets.yaml": &fstest.MapFile{Data: []byte("api_key: secret123\n")},
	}

	e := yamlexpr.New(yamlexpr.WithFS(fs))

	template := map[string]any{
		"include": []any{"config.yaml", "secrets.yaml"},
		"debug":   true,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "production", m["env"])
	require.Equal(t, "secret123", m["api_key"])
	require.Equal(t, true, m["debug"])
	require.NotContains(t, m, "include")
}

// TestIncludeHandler_MergeOverride tests that embedded files can override values.
func TestIncludeHandler_MergeOverride(t *testing.T) {
	fs := fstest.MapFS{
		"defaults.yaml":   &fstest.MapFile{Data: []byte("env: dev\nport: 3000\n")},
		"production.yaml": &fstest.MapFile{Data: []byte("env: production\nport: 8080\n")},
	}

	e := yamlexpr.New(yamlexpr.WithFS(fs))

	template := map[string]any{
		"include": []any{"defaults.yaml", "production.yaml"},
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	// Production overrides defaults
	require.Equal(t, "production", m["env"])
	require.Equal(t, 8080, m["port"])
}

// TestIncludeHandler_WithCustomSyntax tests include with custom directive name.
func TestIncludeHandler_WithCustomSyntax(t *testing.T) {
	fs := fstest.MapFS{
		"config.yaml": &fstest.MapFile{Data: []byte("env: production\n")},
	}

	e := yamlexpr.New(yamlexpr.WithFS(fs), yamlexpr.WithSyntax(yamlexpr.Syntax{Include: "v-include"}))

	template := map[string]any{
		"v-include": "config.yaml",
		"debug":     true,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "production", m["env"])
	require.Equal(t, true, m["debug"])
	require.NotContains(t, m, "v-embed")
}

// TestIncludeHandler_InvalidFile tests including a non-existent file.
func TestIncludeHandler_InvalidFile(t *testing.T) {
	fs := fstest.MapFS{
		"exists.yaml": &fstest.MapFile{Data: []byte("env: production\n")},
	}

	e := yamlexpr.New(yamlexpr.WithFS(fs))

	template := map[string]any{
		"include": "missing.yaml",
	}

	_, err := e.Process(template, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing.yaml")
}

// TestIncludeHandler_InvalidType tests including with wrong type.
func TestIncludeHandler_InvalidType(t *testing.T) {
	e := yamlexpr.New()

	template := map[string]any{
		"include": 123, // Should be string or []string
	}

	_, err := e.Process(template, nil)
	require.Error(t, err)
}

// TestIncludeHandler_NestedEmbeds tests including files that themselves embed other files.
func TestIncludeHandler_NestedEmbeds(t *testing.T) {
	fs := fstest.MapFS{
		"level1.yaml": &fstest.MapFile{Data: []byte("level: 1\nchild:\n  include: level2.yaml\n")},
		"level2.yaml": &fstest.MapFile{Data: []byte("level: 2\n")},
	}

	e := yamlexpr.New(yamlexpr.WithFS(fs))

	template := map[string]any{
		"include": "level1.yaml",
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, 1, m["level"])

	// Child should have merged level2
	child, ok := m["child"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, 2, child["level"])
}

// TestIncludeHandler_WithRealFiles tests including with real filesystem.
func TestIncludeHandler_WithRealFiles(t *testing.T) {
	// Use the testdata/fixtures directory that should exist
	fs := os.DirFS("testdata/fixtures")

	e := yamlexpr.New(yamlexpr.WithFS(fs))

	// Try to load a fixture that exists
	docs, err := e.Load("001-simple-pass-through.yaml")
	require.NoError(t, err)
	require.NotNil(t, docs)
	require.NotEmpty(t, docs)
}

// TestIncludeHandler_EmptyList tests including an empty list of files.
func TestIncludeHandler_EmptyList(t *testing.T) {
	e := yamlexpr.New()

	template := map[string]any{
		"include": []any{},
		"name":    "test",
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "test", m["name"])
	require.NotContains(t, m, "include")
}
