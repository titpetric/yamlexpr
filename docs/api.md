# Package yamlexpr

```go
import (
	"github.com/titpetric/yamlexpr"
}
```

## Types

```go
// Expr evaluates YAML documents with variable interpolation, conditionals, and composition.
type Expr struct {
	fs fs.FS
}
```

```go
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
```

```go
// ExprContextOptions holds configurable options for a new ExprContext.
type ExprContextOptions struct {
	// Stack is the resolver stack for variable lookups.
	Stack *stack.Stack

	// Path is the initial path in the document (defaults to "").
	Path string

	// IncludeChain is the initial chain of included files.
	IncludeChain []string
}
```

```go
// ForLoopExpr represents a parsed for loop expression.
type ForLoopExpr struct {
	// Variables is a list of variable names to bind. Can include "_" to omit.
	Variables []string

	// Source is the name of the variable to iterate over.
	Source string
}
```

## Function symbols

- `func New (rootFS fs.FS) *Expr`
- `func NewExprContext (options *ExprContextOptions) *ExprContext`
- `func (*Expr) Load (filename string) (map[string]any, error)`
- `func (*Expr) Process (doc any, rootVars map[string]any) (any, error)`
- `func (*Expr) ProcessWithStack (doc any, st *stack.Stack) (any, error)`
- `func (*ExprContext) AppendPath (segment string) *ExprContext`
- `func (*ExprContext) FormatIncludeChain () string`
- `func (*ExprContext) Path () string`
- `func (*ExprContext) PopStackScope ()`
- `func (*ExprContext) PushStackScope (m map[string]any)`
- `func (*ExprContext) Stack () *stack.Stack`
- `func (*ExprContext) WithInclude (filename string) *ExprContext`
- `func (*ExprContext) WithPath (newPath string) *ExprContext`

### New

New creates a new Expr evaluator with the given filesystem for includes.

```go
func New(rootFS fs.FS) *Expr
```

### NewExprContext

NewExprContext returns an ExprContext initialized for the given options.

```go
func NewExprContext(options *ExprContextOptions) *ExprContext
```

### Load

Load loads a YAML file and processes it with expression evaluation. Returns a map[string]any containing the processed YAML data. The filename is resolved relative to the filesystem provided to New().

```go
func (*Expr) Load(filename string) (map[string]any, error)
```

### Process

Process processes a YAML document (any) with expression evaluation. Handles for loops, if conditions, includes, and variable interpolation. Root-level keys in the document are available as variables.

```go
func (*Expr) Process(doc any, rootVars map[string]any) (any, error)
```

### ProcessWithStack

ProcessWithStack processes a YAML document with a given variable stack.

```go
func (*Expr) ProcessWithStack(doc any, st *stack.Stack) (any, error)
```

### AppendPath

AppendPath appends a segment to the current path. For example, AppendPath("key") on a context with path "config" results in "config.key". AppendPath("[0]") results in "config[0]".

```go
func (*ExprContext) AppendPath(segment string) *ExprContext
```

### FormatIncludeChain

FormatIncludeChain returns the include chain formatted for error messages. Example: "config.yaml -> database.yaml -> secrets.yaml"

```go
func (*ExprContext) FormatIncludeChain() string
```

### Path

Path returns the current path in the document.

```go
func (*ExprContext) Path() string
```

### PopStackScope

PopStackScope pops the top-most variable scope from the stack.

```go
func (*ExprContext) PopStackScope()
```

### PushStackScope

PushStackScope pushes a new variable scope onto the stack. Used when entering a for loop iteration or other scoped contexts.

```go
func (*ExprContext) PushStackScope(m map[string]any)
```

### Stack

Stack returns the variable resolution stack.

```go
func (*ExprContext) Stack() *stack.Stack
```

### WithInclude

WithInclude returns a new context extended with a filename in the include chain. Used when processing included files to track the chain of includes.

```go
func (*ExprContext) WithInclude(filename string) *ExprContext
```

### WithPath

WithPath returns a new context with the path updated. Useful for tracking location while descending into nested structures.

```go
func (*ExprContext) WithPath(newPath string) *ExprContext
```
