package model

import (
	"strings"

	"github.com/titpetric/yamlexpr/stack"
)

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

// Push pushes a new variable scope onto the stack.
// Used when entering a for loop iteration or other scoped contexts.
func (ctx *Context) Push(m map[string]any) {
	ctx.stack.Push(m)
}

// Pop pops the top-most variable scope from the stack.
func (ctx *Context) Pop() {
	ctx.stack.Pop()
}

// Count returns the number of stack frames.
func (ctx *Context) Count() int {
	return ctx.stack.Count()
}
