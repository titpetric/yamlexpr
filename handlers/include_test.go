package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/handlers"
	"github.com/titpetric/yamlexpr/model"
)

// mockProcessor implements the Processor interface for testing
type mockProcessor struct {
	loadAndMergeCalls []struct {
		filename string
		result   map[string]any
	}
	onLoadAndMerge func(ctx *model.Context, filename string, result map[string]any) error
}

func (m *mockProcessor) ProcessWithContext(ctx *model.Context, doc any) ([]any, error) {
	return nil, nil
}

func (m *mockProcessor) ProcessMapWithContext(ctx *model.Context, ma map[string]any) ([]any, error) {
	return nil, nil
}

func (m *mockProcessor) LoadAndMergeFileWithContext(ctx *model.Context, filename string, result map[string]any) error {
	m.loadAndMergeCalls = append(m.loadAndMergeCalls, struct {
		filename string
		result   map[string]any
	}{filename, result})

	if m.onLoadAndMerge != nil {
		return m.onLoadAndMerge(ctx, filename, result)
	}

	// Default: populate result with test data
	result["embedded"] = true
	result["file"] = filename
	return nil
}

func TestInclude_SingleFile(t *testing.T) {
	proc := &mockProcessor{
		onLoadAndMerge: func(ctx *model.Context, filename string, result map[string]any) error {
			result["config"] = map[string]any{"host": "localhost"}
			return nil
		},
	}

	handler := handlers.Include(proc, "embed")
	ctx := model.NewContext(nil)
	result, consumed, err := handler(ctx, map[string]any{}, "config.yaml")

	require.NoError(t, err)
	require.False(t, consumed)
	require.NotNil(t, result)
	require.IsType(t, []any{}, result)

	want := map[string]any{
		"config": map[string]any{
			"host": "localhost",
		},
	}

	require.Equal(t, want, result[0])
}

func TestInclude_MultipleFiles(t *testing.T) {
	proc := &mockProcessor{
		onLoadAndMerge: func(ctx *model.Context, filename string, result map[string]any) error {
			switch filename {
			case "file1.yaml":
				result["key1"] = "value1"
			case "file2.yaml":
				result["key2"] = "value2"
			}
			return nil
		},
	}

	handler := handlers.Include(proc, "embed")
	ctx := model.NewContext(nil)
	files := []any{"file1.yaml", "file2.yaml"}
	result, consumed, err := handler(ctx, map[string]any{}, files)

	require.NoError(t, err)
	require.False(t, consumed)
	require.NotNil(t, result)

	want := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	require.Equal(t, want, result[0])
	require.Equal(t, 2, len(proc.loadAndMergeCalls))
}

func TestInclude_InvalidType(t *testing.T) {
	proc := &mockProcessor{}
	handler := handlers.Include(proc, "include")
	ctx := model.NewContext(nil)

	_, _, err := handler(ctx, map[string]any{}, 123)

	require.Error(t, err)
	require.Contains(t, err.Error(), "include must be a string or list of strings")
}

func TestInclude_MixedTypes(t *testing.T) {
	proc := &mockProcessor{
		onLoadAndMerge: func(ctx *model.Context, filename string, result map[string]any) error {
			result[filename] = true
			return nil
		},
	}

	handler := handlers.Include(proc, "embed")
	ctx := model.NewContext(nil)

	// List with mixed types: string and int
	files := []any{"file1.yaml", 123, "file2.yaml"}
	res, consumed, err := handler(ctx, map[string]any{}, files)

	require.NoError(t, err)
	require.False(t, consumed)
	require.NotNil(t, res)

	// Should have processed only the string files
	require.Equal(t, 2, len(proc.loadAndMergeCalls))
	require.Equal(t, "file1.yaml", proc.loadAndMergeCalls[0].filename)
	require.Equal(t, "file2.yaml", proc.loadAndMergeCalls[1].filename)
}

func TestInclude_EmptyList(t *testing.T) {
	proc := &mockProcessor{}
	handler := handlers.Include(proc, "embed")
	ctx := model.NewContext(nil)

	res, consumed, err := handler(ctx, map[string]any{}, []any{})

	require.NoError(t, err)
	require.False(t, consumed)
	require.NotNil(t, res)
	require.Equal(t, 0, len(proc.loadAndMergeCalls))
}
