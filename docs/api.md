# Package yamlexpr

```go
import (
	"github.com/titpetric/yamlexpr"
}
```

## Types

```go
// Config holds configuration options for the Expr evaluator.
type Config struct {
	// syntax defines the directive keywords used in YAML documents
	syntax Syntax
	// handlers maps directive names to their handler functions
	handlers map[string]DirectiveHandler
	// handlerOrder tracks the order handlers were registered (for deterministic evaluation)
	handlerOrder []string
	// filesystem is the FS used for loading resources (can be nil)
	filesystem fs.FS
}
```

```go
// ConfigOption is a functional option for configuring an Expr instance.
type ConfigOption func(*Config)
```

```go
// DirectiveHandler is a backwards-compatible alias for model.DirectiveHandler.
type DirectiveHandler = model.DirectiveHandler
```

```go
// Expr evaluates YAML documents with variable interpolation, conditionals, and composition.
type Expr struct {
	fs     fs.FS
	config *Config
}
```

```go
// ExprContext is an alias for model.Context.
type ExprContext = model.Context
```

```go
// ExprContextOptions is an alias for model.ContextOptions.
type ExprContextOptions = model.ContextOptions
```

```go
// ForLoopExpr represents a parsed for loop expression.
// It holds the variable names to bind and the source collection to iterate over.
//
// Variables can include "_" to omit a specific binding position (e.g., ignoring index).
//
// Example: for the expression "item in items", Variables is ["item"] and Source is "items".
// For the expression "(idx, item) in items", Variables is ["idx", "item"] and Source is "items".
type ForLoopExpr struct {
	// Variables is a list of variable names to bind. Can include "_" to omit.
	Variables []string

	// Source is the name of the variable to iterate over.
	Source string
}
```

```go
// Matrix is an alias for handlers.MatrixDirective.
type Matrix = handlers.MatrixDirective
```

```go
// Syntax defines the directive keywords used in YAML documents.
// Empty fields retain their default values when merged with defaults.
type Syntax struct {
	// If is the directive keyword for conditional blocks (default: "if").
	If string `json:"if" yaml:"if"`
	// For is the directive keyword for iteration blocks (default: "for").
	For string `json:"for" yaml:"for"`
	// Include is the directive keyword for file inclusion/composition (default: "include").
	Include string `json:"include" yaml:"include"`
}
```

## Vars

```go
// DefaultSyntax is the default syntax configuration with standard directive names.
var DefaultSyntax = Syntax{
	If:      "if",
	For:     "for",
	Include: "include",
}
```

## Function symbols

- `func DefaultConfig () *Config`
- `func New (opts ...ConfigOption) *Expr`
- `func NewExprContext (options *ExprContextOptions) *ExprContext`
- `func WithDirectiveHandler (directive string, handler DirectiveHandler) ConfigOption`
- `func WithFS (filesystem fs.FS) ConfigOption`
- `func WithSyntax (syntax Syntax) ConfigOption`
- `func (*Config) ForDirective () string`
- `func (*Config) IfDirective () string`
- `func (*Config) IncludeDirective () string`
- `func (*Expr) Load (filename string) ([]any, error)`
- `func (*Expr) LoadAndMergeFileWithContext (ctx *model.Context, filename string, result map[string]any) error`
- `func (*Expr) Process (doc any, rootVars map[string]any) (any, error)`
- `func (*Expr) ProcessMapWithContext (ctx *model.Context, m map[string]any) (any, error)`
- `func (*Expr) ProcessWithContext (ctx *model.Context, doc any) (any, error)`
- `func (*Expr) ProcessWithStack (st *stack.Stack, doc any) (any, error)`
- `func (*Expr) RegisterHandler (directive string, handler DirectiveHandler)`

### DefaultConfig

DefaultConfig returns the default configuration with standard directive names.

```go
func DefaultConfig() *Config
```

### New

New creates a new Expr evaluator with standard handlers registered. Standard handlers include: for, if, include, and matrix directives. Optional ConfigOption arguments can be passed to customize directive syntax or filesystem.

Example:

```
e := yamlexpr.New(yamlexpr.WithFS(myFS))
e := yamlexpr.New(yamlexpr.WithFS(myFS), yamlexpr.WithSyntax(custom))
e := yamlexpr.New()  // No filesystem, handlers registered
```

```go
func New(opts ...ConfigOption) *Expr
```

### NewExprContext

NewExprContext creates a new evaluation context with the given options.

```go
func NewExprContext(options *ExprContextOptions) *ExprContext
```

### WithDirectiveHandler

WithDirectiveHandler registers a custom handler for a directive name. The handler will be called for any block containing the specified directive.

Example:

```
e := yamlexpr.New(fs,
	yamlexpr.WithDirectiveHandler("matrix", myMatrixHandler),
	yamlexpr.WithDirectiveHandler("repeat", myRepeatHandler),
)
```

If a handler is registered for a built-in directive (if, for, include), it overrides the default implementation for that directive.

```go
func WithDirectiveHandler(directive string, handler DirectiveHandler) ConfigOption
```

### WithFS

WithFS sets the filesystem for resource loading (include directive). If not set, only in-memory processing is available.

Example:

```
e := yamlexpr.New(yamlexpr.WithFS(myFS))
```

```go
func WithFS(filesystem fs.FS) ConfigOption
```

### WithSyntax

WithSyntax sets custom directive syntax, preserving defaults for empty fields. Empty string values in the Syntax struct will use the default keywords.

Example:

```
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:      "v-if",
	For:     "v-for",
	Include: "v-include",
}))
```

Or partially customize (empty fields keep defaults):

```
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:  "v-if",
	For: "v-for",
	// Include remains "include"
}))
```

```go
func WithSyntax(syntax Syntax) ConfigOption
```

### ForDirective

ForDirective returns the current for directive keyword.

```go
func (*Config) ForDirective() string
```

### IfDirective

IfDirective returns the current if directive keyword.

```go
func (*Config) IfDirective() string
```

### IncludeDirective

IncludeDirective returns the current include directive keyword.

```go
func (*Config) IncludeDirective() string
```

### Load

Load loads a YAML file and processes it with expression evaluation. Returns a slice of documents. For root-level for: or matrix: directives, returns multiple documents (one per iteration/combination). For regular documents, returns a single-item slice. The filename is resolved relative to the filesystem provided to New().

```go
func (*Expr) Load(filename string) ([]any, error)
```

### LoadAndMergeFileWithContext

LoadAndMergeFileWithContext implements model.Processor interface.

```go
func (*Expr) LoadAndMergeFileWithContext(ctx *model.Context, filename string, result map[string]any) error
```

### Process

Process processes a YAML document (any) with expression evaluation. Handles for loops, if conditions, includes, and variable interpolation. Root-level keys in the document are available as variables.

```go
func (*Expr) Process(doc any, rootVars map[string]any) (any, error)
```

### ProcessMapWithContext

ProcessMapWithContext implements model.Processor interface.

```go
func (*Expr) ProcessMapWithContext(ctx *model.Context, m map[string]any) (any, error)
```

### ProcessWithContext

ProcessWithContext implements model.Processor interface.

```go
func (*Expr) ProcessWithContext(ctx *model.Context, doc any) (any, error)
```

### ProcessWithStack

ProcessWithStack processes a YAML document with a given variable stack.

```go
func (*Expr) ProcessWithStack(st *stack.Stack, doc any) (any, error)
```

### RegisterHandler

RegisterHandler registers a directive handler after Expr creation.

```go
func (*Expr) RegisterHandler(directive string, handler DirectiveHandler)
```
