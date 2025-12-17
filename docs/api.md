# Package yamlexpr

```go
import (
	"github.com/titpetric/yamlexpr"
}
```

## Types

```go
// Document represents a single YAML document after processing.
type Document map[string]any
```

```go
// DocumentContent represents a YAML document with optional frontmatter and sections.
type DocumentContent struct {
	// Frontmatter is the first yaml code block at the start of the file.
	Frontmatter map[string]any
	// Sections are the remaining yaml blocks in the file.
	// With fixtures, the first section is the yamlexpr "template",
	// while the remaining sections are document outputs from Parse().
	Sections []string
}
```

```go
// Expr evaluates YAML documents with variable interpolation, conditionals, and composition.
type Expr struct {
	fs     fs.FS
	config *Config
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

```go
// MatrixDirective represents the parsed matrix configuration
// Fields are exported for testing purposes.
type MatrixDirective struct {
	// Dimensions contains array values that form the cartesian product.
	Dimensions map[string][]any
	// Variables contains non-array values to add to each combination.
	Variables map[string]any
	// Include specifies additional custom combinations to add.
	Include []map[string]any
	// Exclude specifies combinations to filter out from the product.
	Exclude []map[string]any
}
```

```go
// Model type aliases.
type (
	// Config aliases model.Config.
	Config = model.Config
	// ConfigOption aliases model.ConfigOption.
	ConfigOption = model.ConfigOption
	// Context aliases model.Context.
	Context = model.Context
	// ContextOptions aliases model.ContextOptions.
	ContextOptions = model.ContextOptions
	// DirectiveHandler aliases model.DirectiveHandler.
	DirectiveHandler = model.DirectiveHandler
	// Syntax aliases model.SyntaxHandler.
	Syntax = model.Syntax
)
```

## Vars

```go
// Model function/value aliases.
var (
	// DefaultConfig aliases model.DefaultConfig.
	DefaultConfig = model.DefaultConfig
	// NewContext aliases model.NewContext.
	NewContext = model.NewContext
	// WithFS aliases model.WithFS.
	WithFS = model.WithFS
	// WithSyntax aliases model.WithFS.
	WithSyntax = model.WithSyntax
)
```

## Function symbols

- `func MapMatchesSpec (m map[string]any, spec map[string]any) bool`
- `func New (rootFS fs.FS, opts ...ConfigOption) *Expr`
- `func ParseDocument (content string) (*DocumentContent, error)`
- `func TrimLeadingDocumentMarker (content string) string`
- `func ValuesEqual (a,b any) bool`
- `func (*DocumentContent) GetFrontmatterField (field string) (string, bool)`
- `func (*DocumentContent) GetFrontmatterFieldWithDefault (field,defaultVal string) string`
- `func (*Expr) Load (filename string) ([]Document, error)`
- `func (*Expr) Parse (doc Document) ([]Document, error)`

### MapMatchesSpec

MapMatchesSpec checks if a map contains all key:value pairs from a specification map. Used for matrix include/exclude matching and other spec-based filtering. Returns true only if every key in spec exists in the map with an equal value.

```go
func MapMatchesSpec(m map[string]any, spec map[string]any) bool
```

### New

New creates a new Expr evaluator with the given filesystem for includes. Optional ConfigOption arguments can be passed to customize directive syntax.

```go
func New(rootFS fs.FS, opts ...ConfigOption) *Expr
```

### ParseDocument

ParseDocument parses a YAML document into frontmatter and sections separated by ---. The document format is:

```
---
title: "Example"
description: "..."
---
first section (usually input)
---
second section (usually output)
---
additional sections...
```

If the document starts with ---, it is trimmed before parsing.

```go
func ParseDocument(content string) (*DocumentContent, error)
```

### TrimLeadingDocumentMarker

TrimLeadingDocumentMarker trims the leading --- from a YAML document if present.

```go
func TrimLeadingDocumentMarker(content string) string
```

### ValuesEqual

ValuesEqual checks if two values are equal, handling primitives and type coercion. Used for comparing values in matrix specs where YAML may parse numbers as float64 or int.

```go
func ValuesEqual(a, b any) bool
```

### GetFrontmatterField

GetFrontmatterField extracts a specific field from frontmatter.

```go
func (*DocumentContent) GetFrontmatterField(field string) (string, bool)
```

### GetFrontmatterFieldWithDefault

GetFrontmatterFieldWithDefault extracts a field or returns a default.

```go
func (*DocumentContent) GetFrontmatterFieldWithDefault(field, defaultVal string) string
```

### Load

Load loads a YAML file and processes it with expression evaluation. Returns a slice of Documents. For root-level for: or similar directives, may return multiple documents. For regular documents, returns a single-item slice. The filename is resolved relative to the filesystem provided to New().

```go
func (*Expr) Load(filename string) ([]Document, error)
```

### Parse

Parse processes a Document (map[string]any) with expression evaluation. Returns a slice of Documents. For root-level for: directives, may return multiple documents. For regular documents, returns a single-item slice.

```go
func (*Expr) Parse(doc Document) ([]Document, error)
```
