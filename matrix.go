package yamlexpr

import (
	"fmt"
	"sort"

	"github.com/titpetric/yamlexpr/model"
)

// MatrixDirective represents the parsed matrix configuration
// Fields are exported for testing purposes.
type MatrixDirective struct {
	// Dimensions contains array values that form the cartesian product.
	Dimensions map[string][]any
	// Variables contains non-array values to add to each combination.
	Variables map[string]any
	// Include specifies additional custom combinations to add.
	Include []map[string]any
	// Exclude specifies combinations to filter out from the product.
	Exclude []map[string]any
}

// parseMatrixDirective converts the matrix map into structured form
func parseMatrixDirective(m map[string]any) (*MatrixDirective, error) {
	md := &MatrixDirective{
		Dimensions: make(map[string][]any),
		Variables:  make(map[string]any),
	}

	// Extract dimensions and variables (everything except include/exclude)
	for k, v := range m {
		if k == "include" || k == "exclude" {
			continue
		}

		// Check if value is an array (dimension) or scalar (variable)
		switch val := v.(type) {
		case []any:
			md.Dimensions[k] = val
		default:
			// Non-array values are treated as variables to add to each job
			md.Variables[k] = v
		}
	}

	// Parse include (optional)
	if incl, ok := m["include"]; ok {
		if inclList, ok := incl.([]any); ok {
			for i, item := range inclList {
				if itemMap, ok := item.(map[string]any); ok {
					md.Include = append(md.Include, itemMap)
				} else {
					return nil, fmt.Errorf("include[%d] must be a map, got %T", i, item)
				}
			}
		} else {
			return nil, fmt.Errorf("include must be an array, got %T", incl)
		}
	}

	// Parse exclude (optional)
	if excl, ok := m["exclude"]; ok {
		if exclList, ok := excl.([]any); ok {
			for i, item := range exclList {
				if itemMap, ok := item.(map[string]any); ok {
					md.Exclude = append(md.Exclude, itemMap)
				} else {
					return nil, fmt.Errorf("exclude[%d] must be a map, got %T", i, item)
				}
			}
		} else {
			return nil, fmt.Errorf("exclude must be an array, got %T", excl)
		}
	}

	return md, nil
}

// expandMatrixBase generates the cartesian product of all dimensions
func expandMatrixBase(md *MatrixDirective) []map[string]any {
	if len(md.Dimensions) == 0 {
		return []map[string]any{}
	}

	// Collect dimension names in sorted order for consistent reproducible ordering
	keys := make([]string, 0, len(md.Dimensions))
	for k := range md.Dimensions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Check if any dimension is empty - if so, return empty result
	for _, k := range keys {
		if len(md.Dimensions[k]) == 0 {
			return []map[string]any{}
		}
	}

	// Generate cartesian product using indices
	result := []map[string]any{}
	indices := make([]int, len(keys))

	for {
		job := make(map[string]any)
		for i, k := range keys {
			job[k] = md.Dimensions[k][indices[i]]
		}
		result = append(result, job)

		// Increment indices (next combination)
		i := len(indices) - 1
		for ; i >= 0; i-- {
			indices[i]++
			if indices[i] < len(md.Dimensions[keys[i]]) {
				break
			}
			indices[i] = 0
		}
		if i < 0 {
			break
		}
	}

	return result
}

// applyExcludes removes jobs matching exclude specs
// A job matches if ALL keys in the exclude spec match the job
func applyExcludes(jobs []map[string]any, excludes []map[string]any) []map[string]any {
	result := make([]map[string]any, 0, len(jobs))

jobLoop:
	for _, job := range jobs {
		// Check if this job matches any exclude spec
		for _, excl := range excludes {
			if MapMatchesSpec(job, excl) {
				continue jobLoop // Skip this job
			}
		}
		result = append(result, job)
	}

	return result
}

// applyIncludes adds or merges include specs into the job matrix
// For each include spec:
//   - If all keys match existing jobs, merge into those jobs
//   - If no match found, create new job with include keys
func applyIncludes(jobs []map[string]any, includes []map[string]any) ([]map[string]any, error) {
	result := make([]map[string]any, len(jobs))
	copy(result, jobs)

	for _, incl := range includes {
		matched := false

		// Find all jobs that match this include spec
		for i, job := range result {
			if MapMatchesSpec(job, incl) {
				// Merge include into job (include can override)
				matched = true
				merged := make(map[string]any)
				for k, v := range job {
					merged[k] = v
				}
				for k, v := range incl {
					merged[k] = v
				}
				result[i] = merged
			}
		}

		// If no match found, create new job
		if !matched {
			newJob := make(map[string]any)
			for k, v := range incl {
				newJob[k] = v
			}
			result = append(result, newJob)
		}
	}

	return result, nil
}

// handleMatrixWithContext processes a matrix directive with Context.
// The matrix directive expands a template for each combination of dimension values.
// It works similar to for loops but creates a cartesian product instead of iterating a single collection.
//
// At root level, matrix returns multiple documents.
// Within a list item, matrix expands to multiple items.
func (e *Expr) handleMatrixWithContext(ctx *model.Context, matrixValue any, m map[string]any) (any, error) {
	// Parse matrix map
	matrixMap, ok := matrixValue.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("matrix must be a map, got %T", matrixValue)
	}

	// Parse matrix directive
	matrixDir, err := parseMatrixDirective(matrixMap)
	if err != nil {
		return nil, fmt.Errorf("error parsing matrix: %w", err)
	}

	// Expand base matrix (cartesian product)
	jobs := expandMatrixBase(matrixDir)

	// Apply exclude rules
	jobs = applyExcludes(jobs, matrixDir.Exclude)

	// Apply include rules
	jobs, err = applyIncludes(jobs, matrixDir.Include)
	if err != nil {
		return nil, fmt.Errorf("error applying include rules: %w", err)
	}

	// Collect all dimension keys for null-filling
	allDimensionKeys := make(map[string]bool)
	for k := range matrixDir.Dimensions {
		allDimensionKeys[k] = true
	}

	// Collect template keys that might need null values
	// These are keys that appear in the template as potential interpolations
	templateKeys := make(map[string]bool)
	for k := range m {
		if k != e.config.MatrixDirective() {
			templateKeys[k] = true
		}
	}

	// Process each job with matrix variables in scope
	result := make([]any, 0, len(jobs))
	for idx, jobVars := range jobs {
		// Merge non-dimension variables (like run: steps) into each job
		for k, v := range matrixDir.Variables {
			jobVars[k] = v
		}

		// Ensure all dimension keys are present (fill missing with null)
		for k := range allDimensionKeys {
			if _, exists := jobVars[k]; !exists {
				jobVars[k] = nil
			}
		}

		// For template keys that aren't dimensions, initialize to null if not set
		// This allows template values like "xcode: ${xcode}" to interpolate to null
		// when the variable isn't in the job
		for k := range templateKeys {
			if _, exists := jobVars[k]; !exists {
				if _, isDim := allDimensionKeys[k]; !isDim {
					// Non-dimension template key, initialize to null
					jobVars[k] = nil
				}
			}
		}

		// Create template copy without matrix directive
		template := make(map[string]any)
		for k, v := range m {
			if k != e.config.MatrixDirective() {
				template[k] = v
			}
		}

		// Create new stack scope with matrix variables
		ctx.Push(jobVars)

		// Create context for this iteration
		itemCtx := ctx.AppendPath(fmt.Sprintf("[%d]", idx))

		// Process template with current job in scope
		expanded, err := e.processMapWithContext(itemCtx, template)
		if err != nil {
			ctx.Pop()
			return nil, err
		}

		// After processing, merge in keys that should always be present
		if expandedMap, ok := expanded.(map[string]any); ok {
			// Add dimension keys with their values
			for dimKey := range matrixDir.Dimensions {
				expandedMap[dimKey] = jobVars[dimKey]
			}

			// Add back non-dimension template keys that resulted in null
			// These are keys like "xcode: ${xcode}" that interpolated but weren't included because value was nil
			for k := range templateKeys {
				if _, isDim := allDimensionKeys[k]; !isDim {
					// This is a non-dimension template key
					// If it's not in expanded but is in jobVars, add it
					if _, exists := expandedMap[k]; !exists {
						if val, hasVal := jobVars[k]; hasVal {
							expandedMap[k] = val
						}
					}
				}
			}
		}

		if expanded != nil {
			result = append(result, expanded)
		}

		// Pop the scope for this iteration
		ctx.Pop()
	}

	return result, nil
}
