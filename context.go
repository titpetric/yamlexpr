package yamlexpr

import (
	"strings"

	"github.com/titpetric/yamlexpr/stack"
)

// ExprContext carries evaluation context and state used during YAML document processing.
// Each Process operation gets its own ExprContext, making concurrent processing safe.
type ExprContext struct {
	// stack holds variable scope and data resolution
	stack *stack.Stack

	// path tracks the current location in the document (e.g., "config.database", "[0].items[1]")
	// used for error messages and context
	path string

	// includeChain tracks the chain of included files for error context
	includeChain []string
}

// ExprContextOptions holds configurable options for a new ExprContext.
type ExprContextOptions struct {
	// Stack is the resolver stack for variable lookups.
	Stack *stack.Stack

	// Path is the initial path in the document (defaults to "").
	Path string

	// IncludeChain is the initial chain of included files.
	IncludeChain []string
}

// NewExprContext returns an ExprContext initialized for the given options.
func NewExprContext(options *ExprContextOptions) *ExprContext {
	if options == nil {
		options = &ExprContextOptions{}
	}

	ctx := &ExprContext{
		stack:        options.Stack,
		path:         options.Path,
		includeChain: options.IncludeChain,
	}

	if ctx.stack == nil {
		ctx.stack = stack.New(nil)
	}
	if ctx.includeChain == nil {
		ctx.includeChain = []string{}
	}

	return ctx
}

// Stack returns the variable resolution stack.
func (ctx *ExprContext) Stack() *stack.Stack {
	return ctx.stack
}

// Path returns the current path in the document.
func (ctx *ExprContext) Path() string {
	return ctx.path
}

// WithPath returns a new context with the path updated.
// Useful for tracking location while descending into nested structures.
func (ctx *ExprContext) WithPath(newPath string) *ExprContext {
	return &ExprContext{
		stack:        ctx.stack,
		path:         newPath,
		includeChain: ctx.includeChain,
	}
}

// AppendPath appends a segment to the current path.
// For example, AppendPath("key") on a context with path "config" results in "config.key".
// AppendPath("[0]") results in "config[0]".
func (ctx *ExprContext) AppendPath(segment string) *ExprContext {
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
func (ctx *ExprContext) WithInclude(filename string) *ExprContext {
	newChain := make([]string, len(ctx.includeChain)+1)
	copy(newChain, ctx.includeChain)
	newChain[len(ctx.includeChain)] = filename
	return &ExprContext{
		stack:        ctx.stack,
		path:         ctx.path,
		includeChain: newChain,
	}
}

// FormatIncludeChain returns the include chain formatted for error messages.
// Example: "config.yaml -> database.yaml -> secrets.yaml"
func (ctx *ExprContext) FormatIncludeChain() string {
	if len(ctx.includeChain) == 0 {
		return ""
	}
	return strings.Join(ctx.includeChain, " -> ")
}

// PushStackScope pushes a new variable scope onto the stack.
// Used when entering a for loop iteration or other scoped contexts.
func (ctx *ExprContext) PushStackScope(m map[string]any) {
	ctx.stack.Push(m)
}

// PopStackScope pops the top-most variable scope from the stack.
func (ctx *ExprContext) PopStackScope() {
	ctx.stack.Pop()
}
