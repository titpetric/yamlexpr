package model

// DirectiveHandler processes a custom YAML directive.
//
// Parameters:
//   - ctx: Expression context with stack, path, include chain (*model.Context)
//   - block: The containing YAML block (map with directive and template keys)
//   - value: The directive value from YAML
//
// Returns:
//   - result: The processed value
//   - nil: Omit this block
//   - map[string]any: Single item to use
//   - []any: Multiple items to expand
//   - consumed: Whether directive handles all processing
//   - true: Skip normal key processing
//   - false: Continue with normal key processing
//   - error: Processing error with context
type DirectiveHandler func(
	ctx *Context,
	block map[string]any,
	value any,
) (result []any, consumed bool, err error)

// Processor provides document processing capabilities to handlers.
// Handlers use this to recursively process YAML documents.
type Processor interface {
	// ProcessWithContext processes a YAML document (any) with the given context.
	ProcessWithContext(ctx *Context, doc any) ([]any, error)

	// ProcessMapWithContext processes a YAML map with the given context.
	ProcessMapWithContext(ctx *Context, m map[string]any) ([]any, error)

	// LoadAndMergeFileWithContext loads a YAML file and merges it with the given context.
	LoadAndMergeFileWithContext(ctx *Context, filename string, result map[string]any) error
}
