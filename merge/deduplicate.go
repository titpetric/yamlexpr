package merge

import (
	"fmt"
	"sort"
	"strings"
)

// Deduplicate recursively deduplicates a map[string]any, removing duplicate values
// from slices. For slices of maps, duplicates are identified by comparing keys and values.
// Returns the deduplicated map.
func Deduplicate(data any) map[string]any {
	m, ok := data.(map[string]any)
	if !ok {
		return make(map[string]any)
	}
	result, _ := deduplicateWithStats(m, "")
	return result
}

// DeduplicateWithStats recursively deduplicates and returns statistics about removed duplicates.
// Returns: deduplicated map and stats map[path]map[item]count_of_duplicates.
func DeduplicateWithStats(data any) (map[string]any, map[string]map[string]int) {
	m, ok := data.(map[string]any)
	if !ok {
		return make(map[string]any), make(map[string]map[string]int)
	}
	return deduplicateWithStats(m, "")
}

func deduplicateWithStats(m map[string]any, pathPrefix string) (map[string]any, map[string]map[string]int) {
	result := make(map[string]any)
	stats := make(map[string]map[string]int)

	for key, value := range m {
		fullPath := buildPath(pathPrefix, key)
		deduped, substats := deduplicateValueWithStats(value, fullPath)
		result[key] = deduped

		// Merge substats
		for path, itemStats := range substats {
			if _, exists := stats[path]; !exists {
				stats[path] = make(map[string]int)
			}
			for item, count := range itemStats {
				stats[path][item] = count
			}
		}
	}
	return result, stats
}

func deduplicateValueWithStats(value any, pathPrefix string) (any, map[string]map[string]int) {
	stats := make(map[string]map[string]int)
	switch v := value.(type) {
	case map[string]any:
		result, substats := deduplicateWithStats(v, pathPrefix)
		return result, substats
	case []any:
		result, substats := deduplicateSliceWithStats(v, pathPrefix)
		return result, substats
	default:
		return value, stats
	}
}

// deduplicateSliceWithStats removes duplicates and tracks them.
// Returns deduplicated slice and map[path]map[value]count showing actual values.
func deduplicateSliceWithStats(slice []any, pathPrefix string) ([]any, map[string]map[string]int) {
	seen := make(map[string]bool)
	var result []any
	stats := make(map[string]map[string]int)

	if _, exists := stats[pathPrefix]; !exists {
		stats[pathPrefix] = make(map[string]int)
	}

	for _, item := range slice {
		if m, ok := item.(map[string]any); ok {
			// Deduplicate maps by comparison key
			key := generateComparisonKey(m)
			if !seen[key] {
				seen[key] = true
				result = append(result, deduplicateMapNoStats(m))
			} else {
				// Track duplicate - show readable value representation
				valueStr := formatMapValue(m)
				stats[pathPrefix][valueStr]++
			}
		} else {
			// Deduplicate scalar values (strings, ints, etc.) by stringified comparison
			key := fmt.Sprintf("%v", item)
			if !seen[key] {
				seen[key] = true
				result = append(result, item)
			} else {
				// Track duplicate - show actual value
				stats[pathPrefix][key]++
			}
		}
	}

	return result, stats
}

// formatMapValue creates a readable representation of a map for display in stats
func formatMapValue(m map[string]any) string {
	if len(m) == 0 {
		return "{}"
	}

	// For simple maps with single field, show it clearly
	if len(m) == 1 {
		for k, v := range m {
			return fmt.Sprintf("%s=%v", k, v)
		}
	}

	// For complex maps, use comparison key format
	return generateComparisonKey(m)
}

// deduplicateMapNoStats is a simple recursive deduplication without stats tracking
func deduplicateMapNoStats(m map[string]any) map[string]any {
	result := make(map[string]any)
	for key, value := range m {
		result[key] = deduplicateValueNoStats(value)
	}
	return result
}

func deduplicateValueNoStats(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return deduplicateMapNoStats(v)
	case []any:
		result, _ := deduplicateSliceWithStats(v, "")
		return result
	default:
		return value
	}
}

func buildPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// generateComparisonKey creates a comparison key from a map's sorted keys and values.
// Format: "key1-key2-value1-value2" for consistent deduplication.
func generateComparisonKey(m map[string]any) string {
	// Collect and sort keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder

	buf.WriteString("{")
	for k, v := range keys {
		if k > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf(`"%s": "%v"`, v, m[v]))
	}
	buf.WriteString("}")

	return buf.String()
}
