package merge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titpetric/yamlexpr/merge"
)

func TestDeduplicate_SimpleSlice(t *testing.T) {
	data := map[string]any{
		"items": []any{
			map[string]any{"id": 1, "name": "Alice"},
			map[string]any{"id": 2, "name": "Bob"},
			map[string]any{"id": 1, "name": "Alice"}, // duplicate
		},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	assert.Equal(t, 2, len(items), "expected 2 unique items")
	assert.Equal(t, int(1), items[0].(map[string]any)["id"])
	assert.Equal(t, int(2), items[1].(map[string]any)["id"])
}

func TestDeduplicate_ComparisonKey(t *testing.T) {
	data := map[string]any{
		"items": []any{
			map[string]any{"name": "Josh", "id": 123},
			map[string]any{"id": 123, "name": "Josh"}, // same content, different key order
		},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	assert.Equal(t, 1, len(items), "expected 1 unique item (key order shouldn't matter)")
}

func TestDeduplicate_NestedMaps(t *testing.T) {
	data := map[string]any{
		"config": map[string]any{
			"database": map[string]any{
				"host": "localhost",
				"port": 5432,
			},
			"rules": []any{
				map[string]any{"path": "*.go", "action": "lint"},
				map[string]any{"path": "*.go", "action": "lint"}, // duplicate
			},
		},
	}

	result := merge.Deduplicate(data)
	config := result["config"].(map[string]any)
	rules := config["rules"].([]any)

	assert.Equal(t, 1, len(rules), "expected 1 unique rule")
}

func TestDeduplicate_DeepNesting(t *testing.T) {
	data := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": []any{
					map[string]any{"value": "a"},
					map[string]any{"value": "a"}, // duplicate
				},
			},
		},
	}

	result := merge.Deduplicate(data)
	level3 := result["level1"].(map[string]any)["level2"].(map[string]any)["level3"].([]any)

	assert.Equal(t, 1, len(level3), "expected 1 unique value in deeply nested slice")
}

func TestDeduplicate_PreservesFirstOccurrence(t *testing.T) {
	data := map[string]any{
		"items": []any{
			map[string]any{"id": 3},
			map[string]any{"id": 1},
			map[string]any{"id": 2},
			map[string]any{"id": 1}, // duplicate of second
		},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	assert.Equal(t, 3, len(items))
	assert.Equal(t, 3, items[0].(map[string]any)["id"])
	assert.Equal(t, 1, items[1].(map[string]any)["id"])
	assert.Equal(t, 2, items[2].(map[string]any)["id"])
}

func TestDeduplicate_EmptySlice(t *testing.T) {
	data := map[string]any{
		"items": []any{},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	assert.Equal(t, 0, len(items))
}

func TestDeduplicate_MultipleSlices(t *testing.T) {
	data := map[string]any{
		"linters": []any{
			map[string]any{"name": "gofmt"},
			map[string]any{"name": "gofmt"}, // duplicate
		},
		"rules": []any{
			map[string]any{"id": "rule1"},
			map[string]any{"id": "rule2"},
			map[string]any{"id": "rule1"}, // duplicate
		},
	}

	result := merge.Deduplicate(data)

	linters := result["linters"].([]any)
	assert.Equal(t, 1, len(linters))

	rules := result["rules"].([]any)
	assert.Equal(t, 2, len(rules))
}

func TestDeduplicate_ScalarValuesInSlice(t *testing.T) {
	data := map[string]any{
		"strings": []any{"a", "b", "a"}, // scalars are now deduplicated
	}

	result := merge.Deduplicate(data)
	strings := result["strings"].([]any)

	assert.Equal(t, 2, len(strings), "scalar values in slices are deduplicated")
	assert.Equal(t, []any{"a", "b"}, strings)
}

func TestDeduplicate_NonMapInput(t *testing.T) {
	result := merge.Deduplicate("string")
	assert.Equal(t, 0, len(result))

	result = merge.Deduplicate(123)
	assert.Equal(t, 0, len(result))

	result = merge.Deduplicate(nil)
	assert.Equal(t, 0, len(result))
}

func TestDeduplicate_EmptyMap(t *testing.T) {
	data := map[string]any{}
	result := merge.Deduplicate(data)

	assert.Equal(t, 0, len(result))
}

func TestDeduplicate_ComplexDifferentValues(t *testing.T) {
	data := map[string]any{
		"configs": []any{
			map[string]any{"env": "prod", "version": "1.0"},
			map[string]any{"env": "dev", "version": "1.0"},
			map[string]any{"env": "prod", "version": "1.0"}, // duplicate of first
		},
	}

	result := merge.Deduplicate(data)
	configs := result["configs"].([]any)

	assert.Equal(t, 2, len(configs))
}

func TestDeduplicate_MixedMapAndNonMapSlice(t *testing.T) {
	// Note: mixing maps and non-maps in same slice
	data := map[string]any{
		"mixed": []any{
			map[string]any{"id": 1},
			"string",
			map[string]any{"id": 1}, // duplicate map
			"string",                // duplicate string
		},
	}

	result := merge.Deduplicate(data)
	mixed := result["mixed"].([]any)

	// Both map and string duplicates should be removed
	assert.Equal(t, 2, len(mixed))
	assert.Equal(t, map[string]any{"id": 1}, mixed[0])
	assert.Equal(t, "string", mixed[1])
}

func TestDeduplicate_LargeNumberOfDuplicates(t *testing.T) {
	items := []any{
		map[string]any{"id": 1, "data": "test"},
	}
	// Create 99 duplicates
	for i := 0; i < 99; i++ {
		items = append(items, map[string]any{"id": 1, "data": "test"})
	}

	data := map[string]any{"items": items}
	result := merge.Deduplicate(data)
	deduped := result["items"].([]any)

	assert.Equal(t, 1, len(deduped), "100 identical items should deduplicate to 1")
}

func TestDeduplicate_NestedMapsWithDifferentStructures(t *testing.T) {
	data := map[string]any{
		"items": []any{
			map[string]any{
				"nested": map[string]any{"key": "value"},
				"id":     1,
			},
			map[string]any{
				"nested": map[string]any{"key": "value"},
				"id":     1, // same structure
			},
		},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	// Nested maps with same keys/values should be considered duplicates
	assert.Equal(t, 1, len(items))
}
