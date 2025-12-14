package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/handlers"
	"github.com/titpetric/yamlexpr/model"
)

func TestDiscardHandler_True(t *testing.T) {
	handler := handlers.DiscardHandlerBuiltin()
	ctx := model.NewContext(nil)
	result, consumed, err := handler(ctx, map[string]any{}, true)

	require.NoError(t, err)
	require.Nil(t, result)
	require.True(t, consumed)
}

func TestDiscardHandler_False(t *testing.T) {
	handler := handlers.DiscardHandlerBuiltin()
	ctx := model.NewContext(nil)
	result, consumed, err := handler(ctx, map[string]any{}, false)

	require.NoError(t, err)
	require.Nil(t, result)
	require.False(t, consumed)
}

func TestDiscardHandler_StringTrue(t *testing.T) {
	handler := handlers.DiscardHandlerBuiltin()
	ctx := model.NewContext(nil)
	result, consumed, err := handler(ctx, map[string]any{}, "true")

	require.NoError(t, err)
	require.Nil(t, result)
	require.True(t, consumed)
}

func TestDiscardHandler_StringFalse(t *testing.T) {
	handler := handlers.DiscardHandlerBuiltin()
	ctx := model.NewContext(nil)
	result, consumed, err := handler(ctx, map[string]any{}, "false")

	require.NoError(t, err)
	require.Nil(t, result)
	require.False(t, consumed)
}

func TestDiscardHandler_InvalidString(t *testing.T) {
	handler := handlers.DiscardHandlerBuiltin()
	ctx := model.NewContext(nil)
	_, _, err := handler(ctx, map[string]any{}, "maybe")

	require.Error(t, err)
	require.Contains(t, err.Error(), "must be boolean")
}

func TestDiscardHandler_InvalidType(t *testing.T) {
	handler := handlers.DiscardHandlerBuiltin()
	ctx := model.NewContext(nil)
	_, _, err := handler(ctx, map[string]any{}, []string{"foo"})

	require.Error(t, err)
	require.Contains(t, err.Error(), "must be boolean")
}

func TestDiscardHandler_Nil(t *testing.T) {
	handler := handlers.DiscardHandlerBuiltin()
	ctx := model.NewContext(nil)
	result, consumed, err := handler(ctx, map[string]any{}, nil)

	require.NoError(t, err)
	require.Nil(t, result)
	require.False(t, consumed)
}

func TestDiscardHandler_Integer(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		consumed bool
	}{
		{"zero is falsy", 0, false},
		{"one is truthy", 1, true},
		{"negative is truthy", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.DiscardHandlerBuiltin()
			ctx := model.NewContext(nil)
			result, consumed, err := handler(ctx, map[string]any{}, tt.value)

			require.NoError(t, err)
			require.Nil(t, result)
			require.Equal(t, tt.consumed, consumed)
		})
	}
}
