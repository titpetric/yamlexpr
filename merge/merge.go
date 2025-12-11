// The merge package provides a map merge utility which considers how often the
// keys are merged. The key is basically a computed json path expression that's
// typical for interpolating values, e.g. `item.name`, `items[0].name`.
//
// In addition, it also provides deduplication functionality. The deduplication
// considers an encoding key which is still valid json/yaml.

package merge

// MergeMap accumulates data from multiple map sources with statistics.
// It performs deep merging of nested structures and tracks key usage and value diversity.
type MergeMap struct {
	data     map[string]any   // merged data
	stats    map[string]int   // key usage count for each JSON path
	distinct map[string][]any // distinct values per JSON path
}

// NewMergeMap constructs an empty MergeMap.
func NewMergeMap() *MergeMap {
	return &MergeMap{
		data:     make(map[string]any),
		stats:    make(map[string]int),
		distinct: make(map[string][]any),
	}
}

// Merge performs a deep merge of source into m.data, recursively handling nested structures.
// Slices are appended. Each key path (in JSON notation like "user.name") gets a +1 in stats,
// and each distinct value is tracked in distinct.
func (m *MergeMap) Merge(source map[string]any) {
	m.mergeRecursive(source, m.data, "")
}

func (m *MergeMap) mergeRecursive(source, target map[string]any, pathPrefix string) {
	for key, sourceValue := range source {
		fullPath := m.buildPath(pathPrefix, key)

		existingValue, exists := target[key]

		// Both are maps: recurse
		if sourceMap, isSourceMap := sourceValue.(map[string]any); isSourceMap {
			if existingMap, isExistingMap := existingValue.(map[string]any); exists && isExistingMap {
				m.recordKey(fullPath, sourceValue)
				m.mergeRecursive(sourceMap, existingMap, fullPath)
				continue
			}
			// Source is map, existing is not: deep copy the map and record recursively
			m.recordKey(fullPath, sourceValue)
			m.recordKeysRecursive(sourceMap, fullPath)
			target[key] = m.deepCopyMap(sourceMap)
			continue
		}

		// Both are slices: append
		if sourceSlice, isSourceSlice := sourceValue.([]any); isSourceSlice {
			m.recordKey(fullPath, sourceValue)
			if existingSlice, isExistingSlice := existingValue.([]any); exists && isExistingSlice {
				target[key] = append(existingSlice, sourceSlice...)
				continue
			}
			// Source is slice, existing is not: copy the slice
			target[key] = append([]any{}, sourceSlice...)
			continue
		}

		// For scalar values or mismatched types: record and overwrite
		m.recordKey(fullPath, sourceValue)
		target[key] = sourceValue
	}
}

// recordKeysRecursive recursively records all keys in a map structure.
func (m *MergeMap) recordKeysRecursive(source map[string]any, pathPrefix string) {
	for key, value := range source {
		fullPath := m.buildPath(pathPrefix, key)
		m.recordKey(fullPath, value)

		if nestedMap, isMap := value.(map[string]any); isMap {
			m.recordKeysRecursive(nestedMap, fullPath)
		}
	}
}

// buildPath constructs a JSON path notation (e.g., "user.name").
func (m *MergeMap) buildPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// recordKey increments the stat counter and appends to distinct values for the path.
func (m *MergeMap) recordKey(path string, value any) {
	m.stats[path]++
	m.distinct[path] = appendIfDistinct(m.distinct[path], value)
}

// appendIfDistinct adds value to the slice if it's not already present.
func appendIfDistinct(slice []any, value any) []any {
	for _, v := range slice {
		if m := areEqual(v, value); m {
			return slice
		}
	}
	return append(slice, value)
}

// areEqual checks if two values are equal (handles basic types and maps/slices).
func areEqual(a, b any) bool {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case int:
		bv, ok := b.(int)
		return ok && av == bv
	case float64:
		bv, ok := b.(float64)
		return ok && av == bv
	case bool:
		bv, ok := b.(bool)
		return ok && av == bv
	case nil:
		return b == nil
	}
	return false
}

// deepCopyMap recursively copies a map structure.
func (m *MergeMap) deepCopyMap(src map[string]any) map[string]any {
	copy := make(map[string]any)
	for key, value := range src {
		if nestedMap, isMap := value.(map[string]any); isMap {
			copy[key] = m.deepCopyMap(nestedMap)
		} else if nestedSlice, isSlice := value.([]any); isSlice {
			copy[key] = append([]any{}, nestedSlice...)
		} else {
			copy[key] = value
		}
	}
	return copy
}

// Data returns the merged data.
func (m *MergeMap) Data() map[string]any {
	return m.data
}

// Stats returns the key usage counts.
func (m *MergeMap) Stats() map[string]int {
	return m.stats
}

// Distinct returns the distinct values per key path.
func (m *MergeMap) Distinct() map[string][]any {
	return m.distinct
}
