package model

import "github.com/titpetric/yamlexpr/stack"

// ContextOptions holds configurable options for a new Context.
type ContextOptions struct {
	// Stack is the resolver stack for variable lookups.
	Stack *stack.Stack

	// Path is the initial path in the document (defaults to "").
	Path string

	// IncludeChain is the initial chain of included files.
	IncludeChain []string
}
