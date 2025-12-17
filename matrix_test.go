package yamlexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestParseMatrixDirective tests the parseMatrixDirective function with various inputs.
func TestParseMatrixDirective(t *testing.T) {
	t.Run("BasicMatrix", func(t *testing.T) {
		m := map[string]any{
			"os":      []any{"linux", "windows"},
			"version": []any{1, 2},
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.NotNil(t, directive)
		require.Len(t, directive.Dimensions, 2)
		require.Len(t, directive.Variables, 0)
		require.Len(t, directive.Include, 0)
		require.Len(t, directive.Exclude, 0)
	})

	t.Run("MatrixWithVariables", func(t *testing.T) {
		m := map[string]any{
			"os":      []any{"linux"},
			"timeout": 300,
			"retries": 3,
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Dimensions, 1)
		require.Len(t, directive.Variables, 2)
		require.Equal(t, 300, directive.Variables["timeout"])
		require.Equal(t, 3, directive.Variables["retries"])
	})

	t.Run("MatrixWithExclude", func(t *testing.T) {
		m := map[string]any{
			"os": []any{"linux", "windows"},
			"exclude": []any{
				map[string]any{"os": "windows"},
			},
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Exclude, 1)
		require.Equal(t, "windows", directive.Exclude[0]["os"])
	})

	t.Run("MatrixWithInclude", func(t *testing.T) {
		m := map[string]any{
			"os": []any{"linux"},
			"include": []any{
				map[string]any{"os": "macos", "arch": "arm64"},
			},
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Include, 1)
		require.Equal(t, "macos", directive.Include[0]["os"])
	})

	t.Run("InvalidIncludeType", func(t *testing.T) {
		m := map[string]any{
			"os": []any{"linux"},
			"include": []any{
				"invalid_string", // Should be map
			},
		}

		_, err := parseMatrixDirective(m)
		require.Error(t, err)
	})

	t.Run("InvalidExcludeType", func(t *testing.T) {
		m := map[string]any{
			"os": []any{"linux"},
			"exclude": []any{
				123, // Should be map
			},
		}

		_, err := parseMatrixDirective(m)
		require.Error(t, err)
	})
}

// TestExpandMatrixBase tests the expandMatrixBase function.
func TestExpandMatrixBase(t *testing.T) {
	t.Run("SimpleMatrix", func(t *testing.T) {
		md := &MatrixDirective{
			Dimensions: map[string][]any{
				"os": {"linux", "windows"},
			},
		}

		combinations := expandMatrixBase(md)
		require.Len(t, combinations, 2)
		require.Equal(t, "linux", combinations[0]["os"])
		require.Equal(t, "windows", combinations[1]["os"])
	})

	t.Run("TwoDimensions", func(t *testing.T) {
		md := &MatrixDirective{
			Dimensions: map[string][]any{
				"os":   {"linux", "windows"},
				"arch": {"x86_64", "arm64"},
			},
		}

		combinations := expandMatrixBase(md)
		require.Len(t, combinations, 4) // 2 OS × 2 arch
	})

	t.Run("EmptyDimensions", func(t *testing.T) {
		md := &MatrixDirective{
			Dimensions: map[string][]any{},
		}
		combinations := expandMatrixBase(md)
		require.Len(t, combinations, 0)
	})

	t.Run("EmptyArray", func(t *testing.T) {
		md := &MatrixDirective{
			Dimensions: map[string][]any{
				"os": {},
			},
		}

		combinations := expandMatrixBase(md)
		require.Len(t, combinations, 0)
	})
}

// TestApplyExcludes tests the applyExcludes function.
func TestApplyExcludes(t *testing.T) {
	t.Run("ExcludeSingleItem", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux", "arch": "x86_64"},
			{"os": "windows", "arch": "x86_64"},
		}

		exclude := []map[string]any{
			{"os": "windows"},
		}

		result := applyExcludes(combinations, exclude)
		require.Len(t, result, 1)
		require.Equal(t, "linux", result[0]["os"])
	})

	t.Run("ExcludeMultipleFields", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux", "arch": "x86_64"},
			{"os": "linux", "arch": "arm64"},
			{"os": "windows", "arch": "x86_64"},
		}

		exclude := []map[string]any{
			{"os": "windows", "arch": "x86_64"},
		}

		result := applyExcludes(combinations, exclude)
		require.Len(t, result, 2)
	})

	t.Run("ExcludeAll", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux"},
		}

		exclude := []map[string]any{
			{"os": "linux"},
		}

		result := applyExcludes(combinations, exclude)
		require.Len(t, result, 0)
	})

	t.Run("NoExcludes", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux"},
			{"os": "windows"},
		}

		result := applyExcludes(combinations, nil)
		require.Len(t, result, 2)
	})
}

// TestApplyIncludes tests the applyIncludes function.
func TestApplyIncludes(t *testing.T) {
	t.Run("AddSingleInclude", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux"},
		}

		include := []map[string]any{
			{"os": "macos", "arch": "arm64"},
		}

		result, err := applyIncludes(combinations, include)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, "macos", result[1]["os"])
		require.Equal(t, "arm64", result[1]["arch"])
	})

	t.Run("AddMultipleIncludes", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux"},
		}

		include := []map[string]any{
			{"os": "macos"},
			{"os": "freebsd"},
		}

		result, err := applyIncludes(combinations, include)
		require.NoError(t, err)
		require.Len(t, result, 3)
	})

	t.Run("NoIncludes", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux"},
		}

		result, err := applyIncludes(combinations, nil)
		require.NoError(t, err)
		require.Len(t, result, 1)
	})

	t.Run("IncludeWithVariables", func(t *testing.T) {
		combinations := []map[string]any{
			{"os": "linux", "timeout": 30},
		}

		include := []map[string]any{
			{"os": "macos", "timeout": 60},
		}

		result, err := applyIncludes(combinations, include)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, 60, result[1]["timeout"])
	})
}

// TestMergeRecursive tests the mergeRecursive function.
func TestMergeRecursive(t *testing.T) {
	t.Run("MergeSimpleMaps", func(t *testing.T) {
		target := map[string]any{"a": 1}
		source := map[string]any{"b": 2}

		mergeRecursive(target, source)
		require.Equal(t, 1, target["a"])
		require.Equal(t, 2, target["b"])
	})

	t.Run("MergeNestedMaps", func(t *testing.T) {
		target := map[string]any{
			"config": map[string]any{"a": 1},
		}
		source := map[string]any{
			"config": map[string]any{"b": 2},
		}

		mergeRecursive(target, source)
		config := target["config"].(map[string]any)
		require.Equal(t, 1, config["a"])
		require.Equal(t, 2, config["b"])
	})

	t.Run("OverwriteValues", func(t *testing.T) {
		target := map[string]any{"a": 1}
		source := map[string]any{"a": 2}

		mergeRecursive(target, source)
		require.Equal(t, 2, target["a"])
	})

	t.Run("MergeEmptySource", func(t *testing.T) {
		target := map[string]any{"a": 1}
		source := map[string]any{}

		mergeRecursive(target, source)
		require.Equal(t, 1, target["a"])
	})

	t.Run("MergeWithArrays", func(t *testing.T) {
		target := map[string]any{
			"items": []any{1, 2},
		}
		source := map[string]any{
			"items": []any{3, 4},
		}

		// Arrays are overwritten, not merged
		mergeRecursive(target, source)
		items := target["items"].([]any)
		require.Len(t, items, 2)
		require.Equal(t, 3, items[0])
	})
}

// TestApplyIncludes_EdgeCases tests include application with multiple scenarios.
func TestApplyIncludes_EdgeCases(t *testing.T) {
	t.Run("append-single-include", func(t *testing.T) {
		jobs := []map[string]any{
			{"os": "linux"},
			{"os": "windows"},
		}
		includes := []map[string]any{
			{"os": "macos", "special": true},
		}

		result, err := applyIncludes(jobs, includes)
		require.NoError(t, err)
		require.Len(t, result, 3)
		require.Equal(t, map[string]any{"os": "macos", "special": true}, result[2])
	})

	t.Run("append-multiple-includes", func(t *testing.T) {
		jobs := []map[string]any{
			{"python": "3.9"},
		}
		includes := []map[string]any{
			{"python": "3.11", "experimental": true},
			{"python": "3.12", "beta": true},
		}

		result, err := applyIncludes(jobs, includes)
		require.NoError(t, err)
		require.Len(t, result, 3)
		require.Equal(t, "3.11", result[1]["python"])
		require.Equal(t, "3.12", result[2]["python"])
	})

	t.Run("empty-includes", func(t *testing.T) {
		jobs := []map[string]any{
			{"os": "linux"},
		}

		result, err := applyIncludes(jobs, []map[string]any{})
		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, jobs, result)
	})

	t.Run("empty-jobs-with-includes", func(t *testing.T) {
		jobs := []map[string]any{}
		includes := []map[string]any{
			{"os": "linux"},
		}

		result, err := applyIncludes(jobs, includes)
		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, map[string]any{"os": "linux"}, result[0])
	})
}

// TestApplyExcludes_EdgeCases tests exclude application with multiple scenarios.
func TestApplyExcludes_EdgeCases(t *testing.T) {
	t.Run("exclude-by-multiple-fields", func(t *testing.T) {
		jobs := []map[string]any{
			{"os": "linux", "arch": "amd64"},
			{"os": "linux", "arch": "arm64"},
			{"os": "windows", "arch": "amd64"},
		}
		excludes := []map[string]any{
			{"os": "linux", "arch": "arm64"},
		}

		result := applyExcludes(jobs, excludes)
		require.Len(t, result, 2)
		require.Equal(t, map[string]any{"os": "linux", "arch": "amd64"}, result[0])
		require.Equal(t, map[string]any{"os": "windows", "arch": "amd64"}, result[1])
	})

	t.Run("exclude-all-jobs", func(t *testing.T) {
		jobs := []map[string]any{
			{"version": "1.0"},
			{"version": "2.0"},
		}
		excludes := []map[string]any{
			{"version": "1.0"},
			{"version": "2.0"},
		}

		result := applyExcludes(jobs, excludes)
		require.Len(t, result, 0)
	})

	t.Run("partial-spec-match-excludes", func(t *testing.T) {
		jobs := []map[string]any{
			{"os": "linux", "arch": "amd64"},
			{"os": "linux", "arch": "arm64"},
		}
		excludes := []map[string]any{
			{"os": "linux"}, // Matches all jobs where os=linux
		}

		result := applyExcludes(jobs, excludes)
		require.Len(t, result, 0) // Both should be excluded
	})

	t.Run("empty-excludes", func(t *testing.T) {
		jobs := []map[string]any{
			{"os": "linux"},
		}

		result := applyExcludes(jobs, []map[string]any{})
		require.Len(t, result, 1)
		require.Equal(t, jobs, result)
	})

	t.Run("exclude-with-superset-exclusion", func(t *testing.T) {
		jobs := []map[string]any{
			{"os": "linux", "arch": "amd64", "version": "1.0"},
			{"os": "linux", "arch": "amd64", "version": "2.0"},
		}
		excludes := []map[string]any{
			{"os": "linux", "arch": "amd64", "version": "1.0"},
		}

		result := applyExcludes(jobs, excludes)
		require.Len(t, result, 1)
		require.Equal(t, "2.0", result[0]["version"])
	})
}

// TestExpandMatrixBase_Cartesian tests cartesian product generation.
func TestExpandMatrixBase_Cartesian(t *testing.T) {
	t.Run("single-dimension-three-values", func(t *testing.T) {
		md := &MatrixDirective{
			Dimensions: map[string][]any{
				"python": {"3.9", "3.10", "3.11"},
			},
		}

		result := expandMatrixBase(md)
		require.Len(t, result, 3)
		require.Equal(t, "3.9", result[0]["python"])
		require.Equal(t, "3.10", result[1]["python"])
		require.Equal(t, "3.11", result[2]["python"])
	})

	t.Run("three-dimensional-matrix", func(t *testing.T) {
		md := &MatrixDirective{
			Dimensions: map[string][]any{
				"os":     {"linux", "windows"},
				"arch":   {"amd64"},
				"python": {"3.9", "3.10"},
			},
		}

		result := expandMatrixBase(md)
		// 2 * 1 * 2 = 4 combinations
		require.Len(t, result, 4)
		// Check first combination (arch, os, python - alphabetical)
		require.Equal(t, "amd64", result[0]["arch"])
		require.Equal(t, "linux", result[0]["os"])
		require.Equal(t, "3.9", result[0]["python"])
	})

	t.Run("all-single-value-dimensions", func(t *testing.T) {
		md := &MatrixDirective{
			Dimensions: map[string][]any{
				"os":     {"linux"},
				"python": {"3.10"},
			},
		}

		result := expandMatrixBase(md)
		require.Len(t, result, 1)
		require.Equal(t, "linux", result[0]["os"])
		require.Equal(t, "3.10", result[0]["python"])
	})
}

// TestParseMatrixDirective_EdgeCases tests matrix parsing with edge cases.
func TestParseMatrixDirective_EdgeCases(t *testing.T) {
	t.Run("all-fields-empty", func(t *testing.T) {
		m := map[string]any{
			"dimensions": map[string][]any{},
			"include":    []any{},
			"exclude":    []any{},
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Dimensions, 0)
		require.Len(t, directive.Include, 0)
		require.Len(t, directive.Exclude, 0)
	})

	t.Run("dimensions-with-variable-values", func(t *testing.T) {
		m := map[string]any{
			"os":   []any{"${os_name}"},
			"arch": []any{"amd64"},
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Dimensions, 2)
	})

	t.Run("mixed-dimensions-and-variables", func(t *testing.T) {
		m := map[string]any{
			"os":      []any{"linux", "windows"},
			"timeout": 300,
			"retries": 3,
			"name":    "test-matrix",
			"enabled": true,
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Dimensions, 1)
		require.Len(t, directive.Variables, 4)
	})

	t.Run("include-with-multiple-maps", func(t *testing.T) {
		m := map[string]any{
			"os": []any{"linux"},
			"include": []any{
				map[string]any{"os": "macos"},
				map[string]any{"os": "windows"},
			},
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Include, 2)
	})

	t.Run("exclude-with-complex-specs", func(t *testing.T) {
		m := map[string]any{
			"os":   []any{"linux", "windows", "macos"},
			"arch": []any{"amd64", "arm64"},
			"exclude": []any{
				map[string]any{"os": "macos", "arch": "arm64"},
				map[string]any{"os": "windows", "arch": "arm64"},
			},
		}

		directive, err := parseMatrixDirective(m)
		require.NoError(t, err)
		require.Len(t, directive.Exclude, 2)
	})
}
