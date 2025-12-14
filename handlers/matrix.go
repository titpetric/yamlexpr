package handlers

import (
	"fmt"
	"sort"

	"github.com/titpetric/yamlexpr/model"
)

// MatrixDirective represents the parsed matrix configuration
type MatrixDirective struct {
	Dimensions map[string][]any
	Include    []map[string]any
	Exclude    []map[string]any
}

// parseMatrixDirective converts the matrix map into structured form
func parseMatrixDirective(m map[string]any) (*MatrixDirective, error) {
	md := &MatrixDirective{
		Dimensions: make(map[string][]any),
	}

	// Extract dimensions (everything except include/exclude)
	for k, v := range m {
		if k == "include" || k == "exclude" {
			continue
		}

		// Convert value to array
		switch val := v.(type) {
		case []any:
			md.Dimensions[k] = val
		default:
			return nil, fmt.Errorf("matrix dimension '%s' must be an array, got %T", k, v)
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

	// Collect dimension names in consistent order for reproducibility
	keys := make([]string, 0, len(md.Dimensions))
	for k := range md.Dimensions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

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
			if jobMatchesSpec(job, excl) {
				continue jobLoop // Skip this job
			}
		}
		result = append(result, job)
	}

	return result
}

// applyEmbeds adds or merges embed specs into the job matrix
// For each embed spec:
//   - If all keys match existing jobs, merge into those jobs
//   - If no match found, create new job with embed keys
func applyEmbeds(jobs []map[string]any, embeds []map[string]any) ([]map[string]any, error) {
	result := make([]map[string]any, len(jobs))
	copy(result, jobs)

	for _, emb := range embeds {
		matched := false

		// Find all jobs that match this embed spec
		for i, job := range result {
			if jobMatchesSpec(job, emb) {
				// Merge embed into job (embed can override)
				matched = true
				merged := make(map[string]any)
				for k, v := range job {
					merged[k] = v
				}
				for k, v := range emb {
					merged[k] = v
				}
				result[i] = merged
			}
		}

		// If no match found, create new job
		if !matched {
			newJob := make(map[string]any)
			for k, v := range emb {
				newJob[k] = v
			}
			result = append(result, newJob)
		}
	}

	return result, nil
}

// jobMatchesSpec returns true if job contains all key:value pairs from spec
// Used for both include matching (merge into matching jobs)
// and exclude matching (remove matching jobs)
func jobMatchesSpec(job map[string]any, spec map[string]any) bool {
	for specKey, specVal := range spec {
		jobVal, exists := job[specKey]
		if !exists {
			return false
		}
		if !valuesEqual(jobVal, specVal) {
			return false
		}
	}
	return true
}

// valuesEqual checks if two values are equal (handles primitives)
func valuesEqual(a, b any) bool {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case int:
		// Try both int and float64 (YAML may parse as float)
		switch bv := b.(type) {
		case int:
			return av == bv
		case float64:
			return float64(av) == bv
		}
		return false
	case float64:
		// Handle both float64 and int
		switch bv := b.(type) {
		case float64:
			return av == bv
		case int:
			return av == float64(bv)
		}
		return false
	case bool:
		bv, ok := b.(bool)
		return ok && av == bv
	default:
		// Fallback: compare string representations
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}
}

// ExpandMatrixAtRoot expands a root-level matrix directive into multiple job configurations.
// Returns a slice of processed maps, one for each matrix combination.
// This is used by Expr.Process to handle root-level matrix expansion.
func Matrix(processor Processor, matrixDirective string) DirectiveHandler {
	return func(ctx *model.Context, template map[string]any, _ any) ([]any, bool, error) {
		returnErr := func(err error) ([]any, bool, error) {
			return nil, false, err
		}

		matrixValue, ok := template[matrixDirective]
		if !ok {
			return returnErr(fmt.Errorf("matrix key not found"))
		}

		matrixMap, ok := matrixValue.(map[string]any)
		if !ok {
			return returnErr(fmt.Errorf("matrix must be a map, got %T", matrixValue))
		}

		// Parse matrix directive
		matrixDir, err := parseMatrixDirective(matrixMap)
		if err != nil {
			return returnErr(fmt.Errorf("error parsing matrix: %w", err))
		}

		// Expand base matrix (cartesian product)
		jobs := expandMatrixBase(matrixDir)

		// Apply exclude rules
		jobs = applyExcludes(jobs, matrixDir.Exclude)

		// Apply include rules
		jobs, err = applyEmbeds(jobs, matrixDir.Include)
		if err != nil {
			return returnErr(fmt.Errorf("error applying include rules: %w", err))
		}

		// Collect all dimension keys for null-filling
		allKeys := make(map[string]bool)
		for _, job := range jobs {
			for k := range job {
				allKeys[k] = true
			}
		}

		// Process each job with matrix variables in scope
		result := make([]any, 0, len(jobs))
		for _, jobVars := range jobs {
			// Ensure all dimension keys are present (fill missing with null)
			for k := range allKeys {
				if _, exists := jobVars[k]; !exists {
					jobVars[k] = nil
				}
			}

			// Create template copy without matrix directive
			jobTemplate := make(map[string]any)
			for k, v := range template {
				if k != "matrix" {
					jobTemplate[k] = v
				}
			}

			// Create context with job variables and process template
			st := ctx.Stack()
			st.Push(jobVars)

			jobCtx := ctx.WithPath(ctx.Path())

			processed, err := processor.ProcessMapWithContext(jobCtx, jobTemplate)
			st.Pop()

			if err != nil {
				return returnErr(err)
			}

			if processed != nil {
				result = append(result, processed)
			}
		}

		return result, false, nil
	}
}
