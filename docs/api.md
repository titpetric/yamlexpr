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
	// registerStandard indicates if standard handlers should be registered
	registerStandard bool
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
// ExprContext is a backwards-compatible alias for model.Context.
type ExprContext = model.Context
```

```go
// ExprContextOptions is a backwards-compatible alias for model.ContextOptions.
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
// Syntax defines the directive keywords used in YAML documents.
// Empty fields retain their default values when merged with defaults.
type Syntax struct {
	// If is the directive keyword for conditional blocks (default: "if").
	If string `json:"if" yaml:"if"`
	// For is the directive keyword for iteration blocks (default: "for").
	For string `json:"for" yaml:"for"`
	// Embed is the directive keyword for file embedding/inclusion (default: "embed").
	Embed string `json:"embed" yaml:"embed"`
}
```

## Vars

```go
// DefaultSyntax is the default syntax configuration with standard directive names.
var DefaultSyntax = Syntax{
	If:    "if",
	For:   "for",
	Embed: "embed",
}
```

## Function symbols

- `func DefaultConfig () *Config`
- `func New (opts ...ConfigOption) *Expr`
- `func NewExprContext (options *ExprContextOptions) *ExprContext`
- `func NewExtended (opts ...ConfigOption) *Expr`
- `func WithDirectiveHandler (directive string, handler DirectiveHandler) ConfigOption`
- `func WithFS (filesystem fs.FS) ConfigOption`
- `func WithStandardHandlers () ConfigOption`
- `func WithSyntax (syntax Syntax) ConfigOption`
- `func (*Config) EmbedDirective () string`
- `func (*Config) ForDirective () string`
- `func (*Config) GetHandler (directive string) DirectiveHandler`
- `func (*Config) IfDirective () string`
- `func (*Expr) Load (filename string) (map[string]any, error)`
- `func (*Expr) LoadAndMergeFileWithContext (ctx *model.Context, filename string, result map[string]any) error`
- `func (*Expr) Process (doc any, rootVars map[string]any) (any, error)`
- `func (*Expr) ProcessMapWithContext (ctx *model.Context, m map[string]any) (any, error)`
- `func (*Expr) ProcessWithContext (ctx *model.Context, doc any) (any, error)`
- `func (*Expr) ProcessWithStack (doc any, st *stack.Stack) (any, error)`
- `func (*Expr) RegisterHandler (directive string, handler DirectiveHandler)`

### DefaultConfig

DefaultConfig returns the default configuration with standard directive names.

```go
func DefaultConfig() *Config
```

### New

New creates a new Expr evaluator with optional filesystem and configuration options. Call with no arguments for a basic evaluator, then use WithFS() and/or WithStandardHandlers() to configure. Optional ConfigOption arguments can be passed to customize directive syntax and handlers. No handlers are registered by default; use WithStandardHandlers() or WithDirectiveHandler() options.

Example:

```
e := yamlexpr.New(myFS, yamlexpr.WithStandardHandlers())
e := yamlexpr.New(myFS)  // No handlers
e := yamlexpr.New()      // No filesystem, no handlers
```

```go
func New(opts ...ConfigOption) *Expr
```

### NewExprContext

NewExprContext is a backwards-compatible alias for model.NewContext

```go
func NewExprContext(options *ExprContextOptions) *ExprContext
```

### NewExtended

NewExtended creates a new Expr evaluator with standard handlers (for, if, embed) already registered. This is a convenience function equivalent to New(WithStandardHandlers(), opts...). ConfigOption arguments can be passed to customize the evaluator, including WithFS() for filesystem access.

Example:

```
e := yamlexpr.NewExtended(yamlexpr.WithFS(myFS))
e := yamlexpr.NewExtended(yamlexpr.WithFS(myFS), yamlexpr.WithSyntax(custom))
```

```go
func NewExtended(opts ...ConfigOption) *Expr
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

If a handler is registered for a built-in directive (if, for, embed), it overrides the default implementation for that directive.

```go
func WithDirectiveHandler(directive string, handler DirectiveHandler) ConfigOption
```

### WithFS

WithFS sets the filesystem for resource loading (embed directive). If not set, only in-memory processing is available.

Example:

```
e := yamlexpr.New(yamlexpr.WithFS(myFS), yamlexpr.WithStandardHandlers())
```

```go
func WithFS(filesystem fs.FS) ConfigOption
```

### WithStandardHandlers

WithStandardHandlers registers the standard handlers (for, if, embed). This is a convenience option to enable the built-in directives.

Example:

```
e := yamlexpr.New(yamlexpr.WithFS(fs), yamlexpr.WithStandardHandlers())
```

This is equivalent to manually registering each handler:

```
e := yamlexpr.New(yamlexpr.WithFS(fs),
	yamlexpr.WithDirectiveHandler("for", handlers.ForHandlerBuiltin(e, "for")),
	yamlexpr.WithDirectiveHandler("if", handlers.IfHandlerBuiltin("if")),
	yamlexpr.WithDirectiveHandler("embed", handlers.EmbedHandlerBuiltin(e, "embed")),
)
```

```go
func WithStandardHandlers() ConfigOption
```

### WithSyntax

WithSyntax sets custom directive syntax, preserving defaults for empty fields. Empty string values in the Syntax struct will use the default keywords.

Example:

```
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:    "v-if",
	For:   "v-for",
	Embed: "v-embed",
}))
```

Or partially customize (empty fields keep defaults):

```
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:  "v-if",
	For: "v-for",
	// Embed remains "embed"
}))
```

```go
func WithSyntax(syntax Syntax) ConfigOption
```

### EmbedDirective

EmbedDirective returns the current embed directive keyword.

```go
func (*Config) EmbedDirective() string
```

### ForDirective

ForDirective returns the current for directive keyword.

```go
func (*Config) ForDirective() string
```

### GetHandler

GetHandler returns the handler for a directive, or nil if not registered.

```go
func (*Config) GetHandler(directive string) DirectiveHandler
```

### IfDirective

IfDirective returns the current if directive keyword.

```go
func (*Config) IfDirective() string
```

### Load

Load loads a YAML file and processes it with expression evaluation. Returns a map[string]any containing the processed YAML data. The filename is resolved relative to the filesystem provided to New().

```go
func (*Expr) Load(filename string) (map[string]any, error)
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
func (*Expr) ProcessWithStack(doc any, st *stack.Stack) (any, error)
```

### RegisterHandler

RegisterHandler registers a directive handler after Expr creation.

```go
func (*Expr) RegisterHandler(directive string, handler DirectiveHandler)
```
