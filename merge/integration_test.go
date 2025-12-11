//go:build integration

package merge_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/yamlexpr/merge"
)

// TestMergeGolangciYAML finds and merges all .golangci.yml files in the workspace.
func TestMergeGolangciYAML(t *testing.T) {
	// Find all .golangci.yml files from /root
	cmd := exec.Command("find", "/root", "-name", ".golangci.yml", "-type", "f")
	output, err := cmd.Output()
	assert.NoError(t, err)

	files := parseFileList(string(output))

	// Skip test if no files found (expected in test environment)
	if len(files) == 0 {
		t.Skip("no .golangci.yml files found in /root")
	}

	// Create merger and merge all configs
	merger := merge.NewMergeMap()
	for _, file := range files {
		data, err := os.ReadFile(file)
		assert.NoError(t, err, "failed to read %s", file)

		var config map[string]any
		err = yaml.Unmarshal(data, &config)
		assert.NoError(t, err, "failed to parse %s", file)

		merger.Merge(config)
	}

	// Verify merged data
	merged := merger.Data()
	assert.NotNil(t, merged)
	assert.NotEmpty(t, merged)

	// Verify stats track all encountered files
	stats := merger.Stats()
	assert.Greater(t, len(stats), 0, "no stats recorded")

	// Log merge statistics
	t.Logf("Merged %d files with %d unique paths", len(files), len(stats))
	t.Logf("Top-level keys: %v", topLevelKeys(merged))

	// Write merged result to disk
	result := map[string]any{
		"data":  merged,
		"stats": stats,
	}

	resultYAML, err := yaml.Marshal(result)
	assert.NoError(t, err)

	err = os.WriteFile("/root/merged.yml", resultYAML, 0o644)
	assert.NoError(t, err, "failed to write merged.yml")

	// Write stats file
	statsData := map[string]any{
		"files_merged":      len(files),
		"unique_paths":      len(stats),
		"total_keys_seen":   countTotalKeys(stats),
		"output_file":       "/root/merged.yml",
		"output_size_bytes": len(resultYAML),
		"output_lines":      countLines(resultYAML),
	}
	statsYAML, _ := yaml.Marshal(statsData)
	os.WriteFile("/root/merged.stats.yml", statsYAML, 0o644)

	t.Logf("Results written to /root/merged.yml")
	t.Logf("Stats written to /root/merged.stats.yml")
}

// TestDeduplicateMergedGolangciYAML deduplicates the previously merged .golangci.yml files.
func TestDeduplicateMergedGolangciYAML(t *testing.T) {
	// Read the merged result
	mergedData, err := os.ReadFile("/root/merged.yml")
	assert.NoError(t, err, "failed to read /root/merged.yml")

	var mergedResult map[string]any
	err = yaml.Unmarshal(mergedData, &mergedResult)
	assert.NoError(t, err, "failed to parse /root/merged.yml")

	// Extract data section
	dataSection := mergedResult["data"].(map[string]any)

	// Deduplicate with stats
	deduplicated, dedupeStats := merge.DeduplicateWithStats(dataSection)

	// Log deduplication stats
	totalDups := 0
	for _, itemStats := range dedupeStats {
		for _, count := range itemStats {
			totalDups += count
		}
	}
	t.Logf("Deduplication removed %d total duplicate items", totalDups)

	// Create result with deduped data and original stats
	result := map[string]any{
		"data":  deduplicated,
		"stats": mergedResult["stats"],
	}

	resultYAML, err := yaml.Marshal(result)
	assert.NoError(t, err)

	err = os.WriteFile("/root/merged-deduplicated.yml", resultYAML, 0o644)
	assert.NoError(t, err, "failed to write merged-deduplicated.yml")

	// Write dedup stats file
	origMergedData, _ := os.ReadFile("/root/merged.yml")

	// Count total duplicate items and unique slice paths
	totalDupItems := 0
	for _, itemStats := range dedupeStats {
		for _, count := range itemStats {
			totalDupItems += count
		}
	}

	dedupeStatsFile := map[string]any{
		"output_file":            "/root/merged-deduplicated.yml",
		"output_size_bytes":      len(resultYAML),
		"output_lines":           countLines(resultYAML),
		"original_size_bytes":    len(origMergedData),
		"size_reduction_pct":     100 - (len(resultYAML) * 100 / len(origMergedData)),
		"total_duplicate_items":  totalDupItems,
		"duplicate_slices_found": len(dedupeStats),
		"duplicates_by_path":     dedupeStats,
	}
	dedupeYAML, _ := yaml.Marshal(dedupeStatsFile)
	os.WriteFile("/root/merged-deduplicated.stats.yml", dedupeYAML, 0o644)

	t.Logf("Deduplicated results written to /root/merged-deduplicated.yml")
	t.Logf("Dedup stats written to /root/merged-deduplicated.stats.yml")
}

// TestDeduplicateSimpleSlice tests deduplication of a simple slice of maps.
func TestDeduplicateSimpleSlice(t *testing.T) {
	data := map[string]any{
		"linters": []any{
			map[string]any{"name": "gofmt", "enabled": true},
			map[string]any{"name": "goimports", "enabled": true},
			map[string]any{"name": "gofmt", "enabled": true}, // duplicate
		},
	}

	result := merge.Deduplicate(data)
	linters := result["linters"].([]any)

	assert.Equal(t, 2, len(linters), "expected 2 unique linters after deduplication")
	assert.Equal(t, "gofmt", linters[0].(map[string]any)["name"])
	assert.Equal(t, "goimports", linters[1].(map[string]any)["name"])
}

// TestDeduplicateNestedMaps tests deduplication with deeply nested structures.
func TestDeduplicateNestedMaps(t *testing.T) {
	data := map[string]any{
		"linters-settings": map[string]any{
			"gofmt": map[string]any{
				"simplify": true,
			},
		},
		"issues": []any{
			map[string]any{"path": "*.pb.go", "text": "generated"},
			map[string]any{"path": "*.pb.go", "text": "generated"}, // duplicate
		},
	}

	result := merge.Deduplicate(data)

	// Verify nested map is preserved
	settings := result["linters-settings"].(map[string]any)
	gofmt := settings["gofmt"].(map[string]any)
	assert.Equal(t, true, gofmt["simplify"])

	// Verify deduplication happened
	issues := result["issues"].([]any)
	assert.Equal(t, 1, len(issues), "expected 1 unique issue after deduplication")
}

// TestDeduplicateWithDifferentKeyOrder tests that key order doesn't affect deduplication.
func TestDeduplicateWithDifferentKeyOrder(t *testing.T) {
	data := map[string]any{
		"items": []any{
			map[string]any{"id": 1, "name": "first"},
			map[string]any{"name": "first", "id": 1}, // same content, different key order
		},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	assert.Equal(t, 1, len(items), "expected 1 unique item (key order shouldn't matter)")
}

// TestDeduplicateMixedTypes tests deduplication with mixed scalar and map slices.
func TestDeduplicateMixedTypes(t *testing.T) {
	data := map[string]any{
		"names": []any{"Alice", "Bob", "Alice"}, // scalars are now deduplicated
		"configs": []any{
			map[string]any{"env": "prod"},
			map[string]any{"env": "prod"}, // duplicate
		},
	}

	result := merge.Deduplicate(data)

	// Scalar slice should be deduplicated
	names := result["names"].([]any)
	assert.Equal(t, 2, len(names), "scalar values in slices are deduplicated")

	// Map slice should be deduplicated
	configs := result["configs"].([]any)
	assert.Equal(t, 1, len(configs), "expected 1 unique config")
}

// TestDeduplicateComplexConfig tests a realistic linter config scenario.
func TestDeduplicateComplexConfig(t *testing.T) {
	data := map[string]any{
		"linters": []any{
			map[string]any{"name": "gofmt"},
			map[string]any{"name": "vet"},
			map[string]any{"name": "gofmt"}, // duplicate
		},
		"linters-settings": map[string]any{
			"gofmt": map[string]any{
				"simplify": true,
			},
			"vet": map[string]any{
				"check-shadowing": true,
			},
		},
		"issues": map[string]any{
			"exclude-rules": []any{
				map[string]any{"path": "test\\.go$", "linters": []any{"gocyclo"}},
				map[string]any{"path": "test\\.go$", "linters": []any{"gocyclo"}}, // duplicate
			},
		},
	}

	result := merge.Deduplicate(data)

	// Check linters deduplication
	linters := result["linters"].([]any)
	assert.Equal(t, 2, len(linters), "expected 2 unique linters")

	// Check nested deduplication
	excludeRules := result["issues"].(map[string]any)["exclude-rules"].([]any)
	assert.Equal(t, 1, len(excludeRules), "expected 1 unique exclude rule")
}

// TestMergeAndDeduplicateWorkflow tests merging multiple configs and deduplicating.
func TestMergeAndDeduplicateWorkflow(t *testing.T) {
	config1 := map[string]any{
		"linters": []any{
			map[string]any{"name": "gofmt"},
			map[string]any{"name": "goimports"},
		},
	}

	config2 := map[string]any{
		"linters": []any{
			map[string]any{"name": "gofmt"}, // duplicate
			map[string]any{"name": "vet"},
		},
	}

	// Merge configs
	merger := merge.NewMergeMap()
	merger.Merge(config1)
	merger.Merge(config2)

	// After merge, linters slice should have all 4 items (with duplicates)
	linters := merger.Data()["linters"].([]any)
	assert.Equal(t, 4, len(linters), "expected 4 linters after merge (with duplicates)")

	// Deduplicate the merged result
	dedup := merge.Deduplicate(merger.Data())
	deduplicatedLinters := dedup["linters"].([]any)
	assert.Equal(t, 3, len(deduplicatedLinters), "expected 3 unique linters after deduplication")
}

// TestDeduplicatePreservesOrder tests that deduplication preserves first occurrence order.
func TestDeduplicatePreservesOrder(t *testing.T) {
	data := map[string]any{
		"items": []any{
			map[string]any{"id": 3, "name": "third"},
			map[string]any{"id": 1, "name": "first"},
			map[string]any{"id": 2, "name": "second"},
			map[string]any{"id": 1, "name": "first"}, // duplicate
		},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	assert.Equal(t, 3, len(items))
	// Order should be preserved from original (first occurrence kept)
	assert.Equal(t, 3, items[0].(map[string]any)["id"])
	assert.Equal(t, 1, items[1].(map[string]any)["id"])
	assert.Equal(t, 2, items[2].(map[string]any)["id"])
}

// TestDeduplicateEmptySlice tests deduplication of empty slices.
func TestDeduplicateEmptySlice(t *testing.T) {
	data := map[string]any{
		"items": []any{},
	}

	result := merge.Deduplicate(data)
	items := result["items"].([]any)

	assert.Equal(t, 0, len(items))
}

// TestDeduplicateNonMapData tests that non-map data returns empty map.
func TestDeduplicateNonMapData(t *testing.T) {
	result := merge.Deduplicate("not a map")
	assert.Equal(t, 0, len(result))

	result = merge.Deduplicate([]string{"a", "b"})
	assert.Equal(t, 0, len(result))

	result = merge.Deduplicate(nil)
	assert.Equal(t, 0, len(result))
}

// Helper: parse file list from find command output
func parseFileList(output string) []string {
	var files []string
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}
	return files
}

// Helper: get top-level keys from a map
func topLevelKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Helper: count total keys in stats map (sum of all counts)
func countTotalKeys(stats map[string]int) int {
	total := 0
	for _, count := range stats {
		total += count
	}
	return total
}

// Helper: count lines in YAML output
func countLines(data []byte) int {
	count := 0
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}
	return count
}
