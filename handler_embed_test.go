package yamlexpr_test

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

// TestNewEmbedHandler_SingleFile tests embedding a single file.
func TestNewEmbedHandler_SingleFile(t *testing.T) {
	fs := fstest.MapFS{
		"base.yaml": &fstest.MapFile{Data: []byte("env: production\nport: 8080\n")},
	}

	e := yamlexpr.NewExtended(yamlexpr.WithFS(fs))

	template := map[string]any{
		"embed": "base.yaml",
		"debug": true,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "production", m["env"])
	require.Equal(t, 8080, m["port"])
	require.Equal(t, true, m["debug"])
	require.NotContains(t, m, "embed") // Directive removed
}

// TestNewEmbedHandler_MultipleFiles tests embedding multiple files.
func TestNewEmbedHandler_MultipleFiles(t *testing.T) {
	fs := fstest.MapFS{
		"config.yaml":  &fstest.MapFile{Data: []byte("env: production\n")},
		"secrets.yaml": &fstest.MapFile{Data: []byte("api_key: secret123\n")},
	}

	e := yamlexpr.NewExtended(yamlexpr.WithFS(fs))

	template := map[string]any{
		"embed": []any{"config.yaml", "secrets.yaml"},
		"debug": true,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "production", m["env"])
	require.Equal(t, "secret123", m["api_key"])
	require.Equal(t, true, m["debug"])
	require.NotContains(t, m, "embed")
}

// TestNewEmbedHandler_MergeOverride tests that embedded files can override values.
func TestNewEmbedHandler_MergeOverride(t *testing.T) {
	fs := fstest.MapFS{
		"defaults.yaml":   &fstest.MapFile{Data: []byte("env: dev\nport: 3000\n")},
		"production.yaml": &fstest.MapFile{Data: []byte("env: production\nport: 8080\n")},
	}

	e := yamlexpr.NewExtended(yamlexpr.WithFS(fs))

	template := map[string]any{
		"embed": []any{"defaults.yaml", "production.yaml"},
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	// Production overrides defaults
	require.Equal(t, "production", m["env"])
	require.Equal(t, 8080, m["port"])
}

// TestNewEmbedHandler_WithCustomSyntax tests embed with custom directive name.
func TestNewEmbedHandler_WithCustomSyntax(t *testing.T) {
	fs := fstest.MapFS{
		"config.yaml": &fstest.MapFile{Data: []byte("env: production\n")},
	}

	e := yamlexpr.NewExtended(yamlexpr.WithFS(fs), yamlexpr.WithSyntax(yamlexpr.Syntax{Embed: "v-embed"}))

	template := map[string]any{
		"v-embed": "config.yaml",
		"debug":   true,
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "production", m["env"])
	require.Equal(t, true, m["debug"])
	require.NotContains(t, m, "v-embed")
}

// TestNewEmbedHandler_InvalidFile tests embedding a non-existent file.
func TestNewEmbedHandler_InvalidFile(t *testing.T) {
	fs := fstest.MapFS{
		"exists.yaml": &fstest.MapFile{Data: []byte("env: production\n")},
	}

	e := yamlexpr.NewExtended(yamlexpr.WithFS(fs))

	template := map[string]any{
		"embed": "missing.yaml",
	}

	_, err := e.Process(template, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing.yaml")
}

// TestNewEmbedHandler_InvalidType tests embedding with wrong type.
func TestNewEmbedHandler_InvalidType(t *testing.T) {
	e := yamlexpr.NewExtended()

	template := map[string]any{
		"embed": 123, // Should be string or []string
	}

	_, err := e.Process(template, nil)
	require.Error(t, err)
}

// TestNewEmbedHandler_NestedEmbeds tests embedding files that themselves embed other files.
func TestNewEmbedHandler_NestedEmbeds(t *testing.T) {
	fs := fstest.MapFS{
		"level1.yaml": &fstest.MapFile{Data: []byte("level: 1\nchild:\n  embed: level2.yaml\n")},
		"level2.yaml": &fstest.MapFile{Data: []byte("level: 2\n")},
	}

	e := yamlexpr.NewExtended(yamlexpr.WithFS(fs))

	template := map[string]any{
		"embed": "level1.yaml",
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

// TestNewEmbedHandler_WithRealFiles tests embedding with real filesystem.
func TestNewEmbedHandler_WithRealFiles(t *testing.T) {
	// Use the testdata/fixtures directory that should exist
	fs := os.DirFS("testdata/fixtures")

	e := yamlexpr.NewExtended(yamlexpr.WithFS(fs))

	// Try to load a fixture that exists
	result, err := e.Load("001-simple-pass-through.yaml")
	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestNewEmbedHandler_EmptyList tests embedding an empty list of files.
func TestNewEmbedHandler_EmptyList(t *testing.T) {
	e := yamlexpr.NewExtended()

	template := map[string]any{
		"embed": []any{},
		"name":  "test",
	}

	result, err := e.Process(template, nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "test", m["name"])
	require.NotContains(t, m, "embed")
}
