package yamlexpr

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/expr-lang/expr"
	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/yamlexpr/stack"
)

// Expr evaluates YAML documents with variable interpolation, conditionals, and composition.
type Expr struct {
	fs     fs.FS
	config *Config
}

// New creates a new Expr evaluator with the given filesystem for includes.
// Optional ConfigOption arguments can be passed to customize directive syntax.
func New(rootFS fs.FS, opts ...ConfigOption) *Expr {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &Expr{
		fs:     rootFS,
		config: config,
	}
}

// Process processes a YAML document (any) with expression evaluation.
// Handles for loops, if conditions, includes, and variable interpolation.
// Root-level keys in the document are available as variables.
func (e *Expr) Process(doc any, rootVars map[string]any) (any, error) {
	if rootVars == nil {
		rootVars = make(map[string]any)
	}
	if rootMap, ok := doc.(map[string]any); ok {
		for k, v := range rootMap {
			rootVars[k] = v
		}
	}

	return e.ProcessWithStack(doc, stack.New(rootVars))
}

// Load loads a YAML file and processes it with expression evaluation.
// Returns a map[string]any containing the processed YAML data.
// The filename is resolved relative to the filesystem provided to New().
func (e *Expr) Load(filename string) (map[string]any, error) {
	data, err := fs.ReadFile(e.fs, filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	// Parse YAML
	parsed, err := parseYAML(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file %s: %w", filename, err)
	}

	// Process the document
	result, err := e.Process(parsed, nil)
	if err != nil {
		return nil, fmt.Errorf("error processing file %s: %w", filename, err)
	}

	// Convert to map[string]any
	resultMap, ok := result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected map[string]any after processing %s, got %T", filename, result)
	}

	return resultMap, nil
}

// ProcessWithStack processes a YAML document with a given variable stack.
func (e *Expr) ProcessWithStack(doc any, st *stack.Stack) (any, error) {
	if st == nil {
		st = stack.New(nil)
	}
	ctx := NewExprContext(&ExprContextOptions{
		Stack: st,
	})
	return e.processWithContext(ctx, doc)
}

// processWithContext is the internal implementation that handles the processing with context.
func (e *Expr) processWithContext(ctx *ExprContext, doc any) (any, error) {
	switch d := doc.(type) {
	case map[string]any:
		return e.processMapWithContext(ctx, d)
	case []any:
		return e.processSliceWithContext(ctx, d)
	case string:
		// Interpolate string values with error context
		return interpolateStringWithContext(d, ctx.Stack(), ctx.Path())
	default:
		// Return primitives as-is
		return d, nil
	}
}

// processMapWithContext processes a map with ExprContext, handling include, for, and if directives.
func (e *Expr) processMapWithContext(ctx *ExprContext, m map[string]any) (any, error) {
	result := make(map[string]any)

	// Check for include directive
	if incl, ok := m[e.config.IncludeDirective()]; ok {
		if err := e.handleIncludeWithContext(ctx, incl, result); err != nil {
			return nil, err
		}
		// Remove include from processing
		delete(m, e.config.IncludeDirective())
	}

	// Check for for directive
	if forExpr, ok := m[e.config.ForDirective()]; ok {
		return e.handleForWithContext(ctx, forExpr, m)
	}

	// Check for if directive
	if ifExpr, ok := m[e.config.IfDirective()]; ok {
		// Evaluate condition with path context
		ok, err := evaluateConditionWithPath(ifExpr, ctx.Stack(), ctx.Path()+"."+e.config.IfDirective())
		if err != nil {
			return nil, err
		}
		if !ok {
			// Return empty map if condition is false (omit the entire block)
			return nil, nil
		}
		// Remove if from processing
		delete(m, e.config.IfDirective())
	}

	// Process remaining keys
	for k, v := range m {
		childCtx := ctx.AppendPath(k)
		processed, err := e.processWithContext(childCtx, v)
		if err != nil {
			return nil, err
		}
		// Only include non-nil results (if: false returns nil)
		if processed != nil {
			result[k] = processed
		}
	}

	return result, nil
}

// processSliceWithContext processes a slice with ExprContext, handling for and if directives.
func (e *Expr) processSliceWithContext(ctx *ExprContext, s []any) (any, error) {
	result := make([]any, 0, len(s))

	for i, item := range s {
		itemCtx := ctx.AppendPath(fmt.Sprintf("[%d]", i))

		// Check if item is a map with for or if directives
		if m, ok := item.(map[string]any); ok {
			// Check for for directive first (should be evaluated before if)
			if forExpr, ok := m[e.config.ForDirective()]; ok {
				processed, err := e.handleForWithContext(itemCtx, forExpr, m)
				if err != nil {
					return nil, err
				}
				// handleFor returns a slice, extend result
				if slice, ok := processed.([]any); ok {
					result = append(result, slice...)
				}
				continue
			}

			// If no for directive, check if directive
			if ifExpr, ok := m[e.config.IfDirective()]; ok {
				// Evaluate condition
				ok, err := evaluateConditionWithPath(ifExpr, ctx.Stack(), itemCtx.Path()+"."+e.config.IfDirective())
				if err != nil {
					return nil, err
				}
				if !ok {
					// Skip this item
					continue
				}
				// Remove if from processing
				delete(m, e.config.IfDirective())
			}
		}

		processed, err := e.processWithContext(itemCtx, item)
		if err != nil {
			return nil, err
		}
		if processed != nil {
			result = append(result, processed)
		}
	}

	return result, nil
}

// handleIncludeWithContext processes an include directive with ExprContext.
func (e *Expr) handleIncludeWithContext(ctx *ExprContext, incl any, result map[string]any) error {
	// Handle single file
	if filename, ok := incl.(string); ok {
		return e.loadAndMergeFileWithContext(ctx, filename, result)
	}

	// Handle list of files
	if files, ok := incl.([]any); ok {
		for _, f := range files {
			if filename, ok := f.(string); ok {
				if err := e.loadAndMergeFileWithContext(ctx, filename, result); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return fmt.Errorf("include must be a string or list of strings, got %T", incl)
}

// loadAndMergeFileWithContext loads a YAML file and merges it into the result with ExprContext.
func (e *Expr) loadAndMergeFileWithContext(ctx *ExprContext, filename string, result map[string]any) error {
	data, err := fs.ReadFile(e.fs, filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filename, err)
	}

	// Parse YAML
	included, err := parseYAML(data)
	if err != nil {
		return fmt.Errorf("error parsing YAML file %s: %w", filename, err)
	}

	// Create new context for included file
	includedCtx := ctx.WithInclude(filename)

	// Process the included document
	processed, err := e.processWithContext(includedCtx, included)
	if err != nil {
		return fmt.Errorf("error processing included file %s: %w", filename, err)
	}

	// Recursively merge into result
	mergeRecursive(result, processed)
	return nil
}

// mergeRecursive recursively merges src into dst.
// For maps: recursively merges nested maps
// For slices: appends items instead of replacing
// For other types: overwrites the value
func mergeRecursive(dst, src any) {
	switch srcVal := src.(type) {
	case map[string]any:
		// Handle map merging
		if dstMap, ok := dst.(map[string]any); ok {
			for k, v := range srcVal {
				if existingVal, exists := dstMap[k]; exists {
					// Key exists, merge recursively
					switch v.(type) {
					case map[string]any:
						// Both are maps, merge recursively
						if _, isMap := existingVal.(map[string]any); isMap {
							mergeRecursive(existingVal, v)
						} else {
							// Destination is not a map, overwrite
							dstMap[k] = v
						}
					case []any:
						// Source is slice, append to existing
						if dstSlice, isSlice := existingVal.([]any); isSlice {
							dstMap[k] = append(dstSlice, v.([]any)...)
						} else {
							// Destination is not a slice, overwrite
							dstMap[k] = v
						}
					default:
						// Scalar value, overwrite
						dstMap[k] = v
					}
				} else {
					// Key doesn't exist, add it
					dstMap[k] = v
				}
			}
		}
	case []any:
		// Slices can't be merged at top level (handled in map case)
	default:
		// Primitive type, can't merge
	}
}

// handleForWithContext processes a for directive with ExprContext.
// The for directive expands a map template for each item in a collection.
// Supports both simple and complex for expressions:
//   - "item in items" - binds each item to 'item'
//   - "(idx, item) in items" - binds index to 'idx' and item to 'item'
//   - Variables can be "_" to omit from the stack
//
// m should contain "for" key and template keys.
func (e *Expr) handleForWithContext(ctx *ExprContext, forExpr any, m map[string]any) (any, error) {
	// Get the collection to iterate over and parse the for expression
	var items []any
	var loopVars *ForLoopExpr

	switch v := forExpr.(type) {
	case []any:
		// Direct array literal: for: [1, 2, 3]
		items = v
		// Default to single "item" variable for direct arrays
		loopVars = &ForLoopExpr{
			Variables: []string{"item"},
			Source:    "",
		}
	case string:
		// Parse as new for expression (e.g., "item in items" or "(idx, item) in items")
		var err error
		loopVars, err = parseForExpr(v)
		if err != nil {
			pathCtx := ""
			if ctx.Path() != "" {
				pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
			}
			return nil, fmt.Errorf("invalid for expression '%s'%s: %w", v, pathCtx, err)
		}

		// Resolve the source variable from the stack
		sourceVal, ok := ctx.Stack().Resolve(loopVars.Source)
		if !ok {
			pathCtx := ""
			if ctx.Path() != "" {
				pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
			}
			return nil, fmt.Errorf("undefined variable '%s'%s", loopVars.Source, pathCtx)
		}

		// Convert source to slice
		if slice, ok := sourceVal.([]any); ok {
			items = slice
		} else {
			pathCtx := ""
			if ctx.Path() != "" {
				pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
			}
			return nil, fmt.Errorf("for: variable '%s' must be an array, got %T%s", loopVars.Source, sourceVal, pathCtx)
		}
	default:
		pathCtx := ""
		if ctx.Path() != "" {
			pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
		}
		return nil, fmt.Errorf("for: expected array or string expression, got %T%s", forExpr, pathCtx)
	}

	// If empty collection, return empty slice
	if len(items) == 0 {
		return []any{}, nil
	}

	// Iterate over items and expand template
	result := make([]any, 0, len(items))
	for idx, item := range items {
		// Build the scope map with only non-underscore variables
		scope := make(map[string]any)
		for i, varName := range loopVars.Variables {
			if varName == "_" {
				// Skip underscore variables (intentional omission)
				continue
			}

			// Bind the appropriate value based on position
			switch i {
			case 0:
				// First variable is usually the item (or index if 2 variables)
				if len(loopVars.Variables) == 2 {
					scope[varName] = idx
				} else {
					scope[varName] = item
				}
			case 1:
				// Second variable is the item
				scope[varName] = item
			default:
				// Additional variables (for potential future use)
				scope[varName] = item
			}
		}

		// Create new stack scope with loop variables
		ctx.PushStackScope(scope)

		// Create a fresh copy of the template for each iteration (all keys except for directive)
		template := make(map[string]any)
		for k, v := range m {
			if k != e.config.ForDirective() {
				template[k] = v
			}
		}

		// Create context for this iteration
		itemCtx := ctx.AppendPath(fmt.Sprintf("[%d]", idx))

		// Process template with current item in scope
		expanded, err := e.processMapWithContext(itemCtx, template)
		if err != nil {
			ctx.PopStackScope()
			return nil, err
		}
		if expanded != nil {
			result = append(result, expanded)
		}

		// Pop the scope for this iteration
		ctx.PopStackScope()
	}

	return result, nil
}

// parseYAML parses YAML data into a map[string]any or []any.
func parseYAML(data []byte) (any, error) {
	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}
	return result, nil
}

// evaluateConditionWithPath evaluates an if condition with path context for error messages.
// Supports:
// - Boolean values: true/false
// - Interpolated expressions: "${item.active}"
// - Direct variable paths: item.active (converted to expressions via go-expr)
// - Complex expressions: item.status == 'active', item.count > 5, etc.
// Returns errors with variable context and path if referenced variables don't exist.
func evaluateConditionWithPath(condition any, st *stack.Stack, path string) (bool, error) {
	switch v := condition.(type) {
	case bool:
		return v, nil
	case string:
		// Check for literal true/false first
		switch v {
		case "true", "1", "yes":
			return true, nil
		case "false", "0", "no", "":
			return false, nil
		}

		// Handle interpolated expressions like "${item.active}"
		if strings.Contains(v, "${") {
			str, err := interpolateStringWithContext(v, st, path)
			if err != nil {
				return false, err
			}
			// After interpolation, try to parse as boolean
			switch str {
			case "true", "1", "yes":
				return true, nil
			case "false", "0", "no", "":
				return false, nil
			default:
				// The interpolated result couldn't be parsed as boolean
				// If it's a string comparison (like "active == 'active'"), quote the left side
				v = quoteUnquotedComparisons(str)
			}
		}

		// Use go-expr to evaluate the expression
		env := st.All()
		program, err := expr.Compile(v, expr.Env(env))
		if err != nil {
			pathCtx := ""
			if path != "" {
				pathCtx = fmt.Sprintf(" at %s", path)
			}
			return false, fmt.Errorf("error compiling expression '%s'%s: %w", v, pathCtx, err)
		}

		result, err := expr.Run(program, env)
		if err != nil {
			pathCtx := ""
			if path != "" {
				pathCtx = fmt.Sprintf(" at %s", path)
			}
			return false, fmt.Errorf("error evaluating expression '%s'%s: %w", v, pathCtx, err)
		}

		// Convert result to boolean
		return isTruthy(result), nil

	case int, int8, int16, int32, int64:
		// Non-zero is true
		return v != 0, nil
	case float32, float64:
		// Non-zero is true
		return v != 0.0, nil
	default:
		pathCtx := ""
		if path != "" {
			pathCtx = fmt.Sprintf(" at %s", path)
		}
		return false, fmt.Errorf("unsupported condition type: %T%s", condition, pathCtx)
	}
}

// quoteUnquotedComparisons adds quotes around unquoted string literals in comparisons.
// For example: "active == 'active'" stays the same, but "active == test" becomes "'active' == 'test'"
func quoteUnquotedComparisons(expr string) string {
	// Check for comparison operators
	operators := []string{"==", "!=", "<", ">", "<=", ">="}

	for _, op := range operators {
		if strings.Contains(expr, op) {
			parts := strings.Split(expr, op)
			if len(parts) == 2 {
				left := strings.TrimSpace(parts[0])
				right := strings.TrimSpace(parts[1])

				// Quote unquoted parts (but not if they're already quoted)
				if !isQuoted(left) && !strings.Contains(left, ".") && !strings.Contains(left, "(") {
					left = "'" + left + "'"
				}
				if !isQuoted(right) && !strings.Contains(right, ".") && !strings.Contains(right, "(") {
					right = "'" + right + "'"
				}

				return left + " " + op + " " + right
			}
		}
	}

	return expr
}

// isQuoted checks if a string is already quoted
func isQuoted(s string) bool {
	s = strings.TrimSpace(s)
	return (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}

// isTruthy returns true for non-empty/non-zero values.
func isTruthy(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int, int8, int16, int32, int64:
		return val != 0
	case uint, uint8, uint16, uint32, uint64:
		return val != 0
	case float32, float64:
		return val != 0.0
	case string:
		return val != ""
	case []any:
		return len(val) > 0
	case map[string]any:
		return len(val) > 0
	default:
		return v != nil
	}
}
