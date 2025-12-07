package yamlexpr

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr/stack"
)

func TestDebugForWithIf(t *testing.T) {
	input := map[string]any{
		"items": []any{
			map[string]any{
				"for": []any{
					map[string]any{"name": "api", "active": true},
					map[string]any{"name": "worker", "active": false},
				},
				"if":   "item.active",
				"name": "${item.name}",
			},
		},
	}

	st := stack.New(nil)
	e := New(nil)

	// Test interpolation with item
	st.Push(map[string]any{"item": map[string]any{"name": "api", "active": true}})
	interp := interpolateStringHelper("${item.active}", st)
	t.Logf("Interpolated '${item.active}' -> '%s'", interp)

	cond, err := evaluateConditionWithPath("${item.active}", st, "")
	require.NoError(t, err)
	t.Logf("Condition '${item.active}' evaluated to: %v", cond)
	st.Pop()

	// Now test the full process
	result, err := e.Process(input, nil)
	require.NoError(t, err)
	t.Logf("Result: %v", result)
}
