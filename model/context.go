package model

import (
	"strings"

	"github.com/titpetric/yamlexpr/stack"
)

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
) (result any, consumed bool, err error)

// Processor provides document processing capabilities to handlers.
// Handlers use this to recursively process YAML documents.
type Processor interface {
	// ProcessWithContext processes a YAML document (any) with the given context.
	ProcessWithContext(ctx *Context, doc any) (any, error)

	// ProcessMapWithContext processes a YAML map with the given context.
	ProcessMapWithContext(ctx *Context, m map[string]any) (any, error)

	// LoadAndMergeFileWithContext loads a YAML file and merges it with the given context.
	LoadAndMergeFileWithContext(ctx *Context, filename string, result map[string]any) error
}

// Context carries evaluation context and state used during YAML document processing.
// Each Process operation gets its own Context, making concurrent processing safe.
type Context struct {
	// stack holds variable scope and data resolution
	stack *stack.Stack

	// path tracks the current location in the document (e.g., "config.database", "[0].items[1]")
	// used for error messages and context
	path string

	// includeChain tracks the chain of included files for error context
	includeChain []string
}

// ContextOptions holds configurable options for a new Context.
type ContextOptions struct {
	// Stack is the resolver stack for variable lookups.
	Stack *stack.Stack

	// Path is the initial path in the document (defaults to "").
	Path string

	// IncludeChain is the initial chain of included files.
	IncludeChain []string
}

// NewContext returns a Context initialized for the given options.
func NewContext(options *ContextOptions) *Context {
	if options == nil {
		options = &ContextOptions{}
	}

	ctx := &Context{
		stack:        options.Stack,
		path:         options.Path,
		includeChain: options.IncludeChain,
	}

	if ctx.stack == nil {
		ctx.stack = stack.New()
	}
	if ctx.includeChain == nil {
		ctx.includeChain = []string{}
	}

	return ctx
}

// Stack returns the variable resolution stack.
func (ctx *Context) Stack() *stack.Stack {
	return ctx.stack
}

// Path returns the current path in the document.
func (ctx *Context) Path() string {
	return ctx.path
}

// WithPath returns a new context with the path updated.
// Useful for tracking location while descending into nested structures.
func (ctx *Context) WithPath(newPath string) *Context {
	return &Context{
		stack:        ctx.stack,
		path:         newPath,
		includeChain: ctx.includeChain,
	}
}

// AppendPath appends a segment to the current path.
// For example, AppendPath("key") on a context with path "config" results in "config.key".
// AppendPath("[0]") results in "config[0]".
func (ctx *Context) AppendPath(segment string) *Context {
	newPath := ctx.path
	if segment == "" {
		return ctx.WithPath(newPath)
	}

	// Handle array index notation
	if strings.HasPrefix(segment, "[") {
		if newPath != "" {
			newPath = newPath + segment
		} else {
			newPath = segment
		}
	} else {
		// Handle regular keys
		if newPath != "" {
			newPath = newPath + "." + segment
		} else {
			newPath = segment
		}
	}

	return ctx.WithPath(newPath)
}

// WithInclude returns a new context extended with a filename in the include chain.
// Used when processing included files to track the chain of includes.
func (ctx *Context) WithInclude(filename string) *Context {
	newChain := make([]string, len(ctx.includeChain)+1)
	copy(newChain, ctx.includeChain)
	newChain[len(ctx.includeChain)] = filename
	return &Context{
		stack:        ctx.stack,
		path:         ctx.path,
		includeChain: newChain,
	}
}

// FormatIncludeChain returns the include chain formatted for error messages.
// Example: "config.yaml -> database.yaml -> secrets.yaml"
func (ctx *Context) FormatIncludeChain() string {
	if len(ctx.includeChain) == 0 {
		return ""
	}
	return strings.Join(ctx.includeChain, " -> ")
}

// PushStackScope pushes a new variable scope onto the stack.
// Used when entering a for loop iteration or other scoped contexts.
func (ctx *Context) PushStackScope(m map[string]any) {
	ctx.stack.Push(m)
}

// PopStackScope pops the top-most variable scope from the stack.
func (ctx *Context) PopStackScope() {
	ctx.stack.Pop()
}
