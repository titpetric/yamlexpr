package merge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titpetric/yamlexpr/merge"
)

func TestNewMergeMap(t *testing.T) {
	m := merge.NewMergeMap()
	assert.NotNil(t, m)
	assert.Equal(t, 0, len(m.Data()))
	assert.Equal(t, 0, len(m.Stats()))
	assert.Equal(t, 0, len(m.Distinct()))
}

func TestMergeMap_SimpleMaps(t *testing.T) {
	m := merge.NewMergeMap()

	// Merge {a, b, c}
	m.Merge(map[string]any{"a": 1, "b": 2, "c": 3})
	assert.Equal(t, map[string]any{"a": 1, "b": 2, "c": 3}, m.Data())
	assert.Equal(t, map[string]int{"a": 1, "b": 1, "c": 1}, m.Stats())

	// Merge {b, c}
	m.Merge(map[string]any{"b": 20, "c": 30})
	assert.Equal(t, map[string]any{"a": 1, "b": 20, "c": 30}, m.Data())
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 2}, m.Stats())

	// Merge {c}
	m.Merge(map[string]any{"c": 300})
	assert.Equal(t, map[string]any{"a": 1, "b": 20, "c": 300}, m.Data())
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, m.Stats())
}

func TestMergeMap_DistinctValues(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{"x": "val1"})
	m.Merge(map[string]any{"x": "val2"})
	m.Merge(map[string]any{"x": "val1"}) // Duplicate

	assert.Equal(t, map[string]any{"x": "val1"}, m.Data())
	assert.Equal(t, map[string]int{"x": 3}, m.Stats())
	assert.Equal(t, []any{"val1", "val2"}, m.Distinct()["x"])
}

func TestMergeMap_DeepStructures(t *testing.T) {
	m := merge.NewMergeMap()

	// Merge nested map
	m.Merge(map[string]any{
		"user": map[string]any{
			"name": "Alice",
			"age":  30,
		},
	})

	assert.Equal(t, map[string]any{
		"user": map[string]any{
			"name": "Alice",
			"age":  30,
		},
	}, m.Data())

	assert.Equal(t, map[string]int{
		"user":      1,
		"user.name": 1,
		"user.age":  1,
	}, m.Stats())
}

func TestMergeMap_DeeperNesting(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"db": map[string]any{
			"config": map[string]any{
				"host": "localhost",
				"port": 5432,
			},
		},
	})

	m.Merge(map[string]any{
		"db": map[string]any{
			"config": map[string]any{
				"port": 3306,
				"user": "admin",
			},
		},
	})

	expected := map[string]any{
		"db": map[string]any{
			"config": map[string]any{
				"host": "localhost",
				"port": 3306,
				"user": "admin",
			},
		},
	}
	assert.Equal(t, expected, m.Data())

	stats := m.Stats()
	assert.Equal(t, 2, stats["db"])
	assert.Equal(t, 2, stats["db.config"])
	assert.Equal(t, 1, stats["db.config.host"])
	assert.Equal(t, 2, stats["db.config.port"])
	assert.Equal(t, 1, stats["db.config.user"])
}

func TestMergeMap_SliceAppending(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"items": []any{1, 2, 3},
	})

	m.Merge(map[string]any{
		"items": []any{4, 5},
	})

	assert.Equal(t, map[string]any{
		"items": []any{1, 2, 3, 4, 5},
	}, m.Data())

	assert.Equal(t, map[string]int{"items": 2}, m.Stats())
}

func TestMergeMap_SliceOverwrite(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"values": "scalar",
	})

	m.Merge(map[string]any{
		"values": []any{1, 2, 3},
	})

	assert.Equal(t, map[string]any{
		"values": []any{1, 2, 3},
	}, m.Data())
}

func TestMergeMap_MixedTypes(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"str":   "hello",
		"num":   42,
		"float": 3.14,
		"bool":  true,
		"null":  nil,
	})

	m.Merge(map[string]any{
		"str":   "world",
		"num":   100,
		"float": 2.71,
		"bool":  false,
		"null":  nil,
	})

	expected := map[string]any{
		"str":   "world",
		"num":   100,
		"float": 2.71,
		"bool":  false,
		"null":  nil,
	}
	assert.Equal(t, expected, m.Data())

	assert.Equal(t, []any{"hello", "world"}, m.Distinct()["str"])
	assert.Equal(t, []any{42, 100}, m.Distinct()["num"])
	assert.Equal(t, []any{3.14, 2.71}, m.Distinct()["float"])
	assert.Equal(t, []any{true, false}, m.Distinct()["bool"])
	assert.Equal(t, []any{nil}, m.Distinct()["null"])
}

func TestMergeMap_NestedSlices(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"config": map[string]any{
			"servers": []any{"srv1", "srv2"},
		},
	})

	m.Merge(map[string]any{
		"config": map[string]any{
			"servers": []any{"srv3"},
		},
	})

	assert.Equal(t, map[string]any{
		"config": map[string]any{
			"servers": []any{"srv1", "srv2", "srv3"},
		},
	}, m.Data())

	assert.Equal(t, 2, m.Stats()["config.servers"])
}

func TestMergeMap_ComplexNestedStructure(t *testing.T) {
	m := merge.NewMergeMap()

	// First merge
	m.Merge(map[string]any{
		"app": map[string]any{
			"name":    "MyApp",
			"version": "1.0",
			"features": map[string]any{
				"auth":  true,
				"cache": false,
			},
			"plugins": []any{"plugin1", "plugin2"},
		},
	})

	// Second merge - partial overwrite
	m.Merge(map[string]any{
		"app": map[string]any{
			"version": "1.1",
			"features": map[string]any{
				"cache": true,
				"logs":  true,
			},
			"plugins": []any{"plugin3"},
		},
	})

	data := m.Data()
	appData := data["app"].(map[string]any)
	assert.Equal(t, "MyApp", appData["name"])
	assert.Equal(t, "1.1", appData["version"])

	features := appData["features"].(map[string]any)
	assert.Equal(t, true, features["auth"])
	assert.Equal(t, true, features["cache"])
	assert.Equal(t, true, features["logs"])

	plugins := appData["plugins"].([]any)
	assert.Equal(t, []any{"plugin1", "plugin2", "plugin3"}, plugins)
}

func TestMergeMap_EmptyMerge(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{"x": 1})
	m.Merge(map[string]any{})

	assert.Equal(t, map[string]any{"x": 1}, m.Data())
	assert.Equal(t, map[string]int{"x": 1}, m.Stats())
}

func TestMergeMap_DeepCopy(t *testing.T) {
	m := merge.NewMergeMap()

	source := map[string]any{
		"nested": map[string]any{
			"value": 42,
		},
	}

	m.Merge(source)

	// Modify original source
	source["nested"].(map[string]any)["value"] = 999

	// Merged data should be unchanged
	assert.Equal(t, 42, m.Data()["nested"].(map[string]any)["value"])
}

func TestMergeMap_MultipleDistinctValues(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"level": map[string]any{
			"setting": "a",
		},
	})

	m.Merge(map[string]any{
		"level": map[string]any{
			"setting": "b",
		},
	})

	m.Merge(map[string]any{
		"level": map[string]any{
			"setting": "c",
		},
	})

	assert.Equal(t, 3, m.Stats()["level.setting"])
	assert.Equal(t, []any{"a", "b", "c"}, m.Distinct()["level.setting"])
}

func TestMergeMap_PathNotation(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	})

	stats := m.Stats()
	assert.Contains(t, stats, "a")
	assert.Contains(t, stats, "a.b")
	assert.Contains(t, stats, "a.b.c")
}

func TestMergeMap_ScalarToMapConversion(t *testing.T) {
	m := merge.NewMergeMap()

	m.Merge(map[string]any{
		"config": "scalar",
	})

	m.Merge(map[string]any{
		"config": map[string]any{
			"option": "value",
		},
	})

	assert.Equal(t, map[string]any{
		"config": map[string]any{
			"option": "value",
		},
	}, m.Data())
}
