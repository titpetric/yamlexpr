package yamlexpr

import (
	"fmt"
	"io/fs"

	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/yamlexpr/handlers"
	"github.com/titpetric/yamlexpr/model"
	"github.com/titpetric/yamlexpr/stack"
)

// Expr evaluates YAML documents with variable interpolation, conditionals, and composition.
type Expr struct {
	fs     fs.FS
	config *Config
}

// RegisterHandler registers a directive handler after Expr creation.
func (e *Expr) RegisterHandler(directive string, handler DirectiveHandler) {
	if e.config.handlers == nil {
		e.config.handlers = make(map[string]DirectiveHandler)
	}
	// Track order of registration if this is a new handler
	if _, exists := e.config.handlers[directive]; !exists {
		e.config.handlerOrder = append(e.config.handlerOrder, directive)
	}
	e.config.handlers[directive] = handler
}

// New creates a new Expr evaluator with optional filesystem and configuration options.
// Call with no arguments for a basic evaluator, then use WithFS() and/or WithStandardHandlers() to configure.
// Optional ConfigOption arguments can be passed to customize directive syntax and handlers.
// No handlers are registered by default; use WithStandardHandlers() or WithDirectiveHandler() options.
//
// Example:
//
//	e := yamlexpr.New(myFS, yamlexpr.WithStandardHandlers())
//	e := yamlexpr.New(myFS)  // No handlers
//	e := yamlexpr.New()      // No filesystem, no handlers
func New(opts ...ConfigOption) *Expr {
	config := DefaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	e := &Expr{
		fs:     config.filesystem,
		config: config,
	}

	// Register standard handlers if requested
	if config.registerStandard {
		e.RegisterHandler(e.config.EmbedDirective(), handlers.EmbedHandlerBuiltin(e, e.config.EmbedDirective()))
		e.RegisterHandler(e.config.ForDirective(), handlers.ForHandlerBuiltin(e, e.config.ForDirective()))
		e.RegisterHandler(e.config.IfDirective(), handlers.IfHandlerBuiltin(e.config.IfDirective()))
	}

	return e
}

// NewExtended creates a new Expr evaluator with standard handlers (for, if, embed) already registered.
// This is a convenience function equivalent to New(WithStandardHandlers(), opts...).
// ConfigOption arguments can be passed to customize the evaluator, including WithFS() for filesystem access.
//
// Example:
//
//	e := yamlexpr.NewExtended(yamlexpr.WithFS(myFS))
//	e := yamlexpr.NewExtended(yamlexpr.WithFS(myFS), yamlexpr.WithSyntax(custom))
func NewExtended(opts ...ConfigOption) *Expr {
	allOpts := make([]ConfigOption, 0, len(opts)+1)
	allOpts = append(allOpts, WithStandardHandlers())
	allOpts = append(allOpts, opts...)
	return New(allOpts...)
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

	return e.ProcessWithStack(doc, stack.NewStack(rootVars))
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
		st = stack.New()
	}
	ctx := model.NewContext(&model.ContextOptions{
		Stack: st,
	})
	return e.processWithContext(ctx, doc)
}

// ProcessWithContext implements model.Processor interface.
func (e *Expr) ProcessWithContext(ctx *model.Context, doc any) (any, error) {
	return e.processWithContext(ctx, doc)
}

// ProcessMapWithContext implements model.Processor interface.
func (e *Expr) ProcessMapWithContext(ctx *model.Context, m map[string]any) (any, error) {
	return e.processMapWithContext(ctx, m)
}

// LoadAndMergeFileWithContext implements model.Processor interface.
func (e *Expr) LoadAndMergeFileWithContext(ctx *model.Context, filename string, result map[string]any) error {
	return e.loadAndMergeFileWithContext(ctx, filename, result)
}

// processWithContext is the internal implementation that handles the processing with context.
func (e *Expr) processWithContext(ctx *model.Context, doc any) (any, error) {
	switch d := doc.(type) {
	case map[string]any:
		return e.processMapWithContext(ctx, d)
	case []any:
		return e.processSliceWithContext(ctx, d)
	case string:
		// Interpolate string values with error context
		return handlers.InterpolateStringWithContext(d, ctx.Stack(), ctx.Path())
	default:
		// Return primitives as-is
		return d, nil
	}
}

// processMapWithContext processes a map with ExprContext, handling include, for, if, and custom directives.
func (e *Expr) processMapWithContext(ctx *model.Context, m map[string]any) (any, error) {
	result := make(map[string]any)
	processedKeys := make(map[string]bool) // Track keys handled by handlers

	// Check for custom handlers in registration order (deterministic evaluation)
	for _, directive := range e.config.handlerOrder {
		handler := e.config.handlers[directive]
		if value, ok := m[directive]; ok {
			// Handler found - call it
			handlerResult, consumed, err := handler(ctx, m, value)
			if err != nil {
				return nil, err
			}

			// If handler consumed all processing, return its result directly
			if consumed {
				return handlerResult, nil
			}

			// Handler returned a result but didn't consume all processing
			// Merge the result into our result map and continue
			if handlerResult != nil {
				if resMap, ok := handlerResult.(map[string]any); ok {
					for k, v := range resMap {
						result[k] = v
						// Mark these keys as processed by the handler
						processedKeys[k] = true
					}
				}
			}

			// Remove the directive from processing
			delete(m, directive)
			processedKeys[directive] = true
		}
	}

	// Process remaining keys (skip those handled by custom handlers)
	for k, v := range m {
		// Skip if already processed by a handler
		if processedKeys[k] {
			continue
		}

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

// processSliceWithContext processes a slice with ExprContext.
// Handlers registered for directives (like for, if) will be called when processing maps.
func (e *Expr) processSliceWithContext(ctx *model.Context, s []any) (any, error) {
	result := make([]any, 0, len(s))

	for i, item := range s {
		itemCtx := ctx.AppendPath(fmt.Sprintf("[%d]", i))

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

// loadAndMergeFileWithContext loads a YAML file and merges it into the result with ExprContext.
func (e *Expr) loadAndMergeFileWithContext(ctx *model.Context, filename string, result map[string]any) error {
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

// parseYAML parses YAML data into a map[string]any or []any.
func parseYAML(data []byte) (any, error) {
	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}
	return result, nil
}
