# yamlexpr

[![Go Reference](https://pkg.go.dev/badge/github.com/titpetric/yamlexpr.svg)](https://pkg.go.dev/github.com/titpetric/yamlexpr) [![Go Report Card](https://goreportcard.com/badge/github.com/titpetric/yamlexpr)](https://goreportcard.com/report/github.com/titpetric/yamlexpr) [![PkgGoDev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/titpetric/yamlexpr) [![Test Coverage](https://img.shields.io/badge/coverage-40.1%25-yellowgreen)](docs/testing-coverage.md)

YAML composition, interpolation, and conditional evaluation for Go.

## Documentation

### Getting Started
- **[Tutorial](docs/tutorial.md)** - Comprehensive guide with real-world examples
- **[Quick Reference](docs/features/QUICK_REFERENCE.md)** - Syntax cheat sheet and common patterns

### Feature Documentation
- **[Interpolation](docs/features/interpolation.md)** - Variable substitution with `${variable}` syntax
- **[Conditionals](docs/features/conditionals.md)** - Include/exclude with `if:` directive
- **[For Loops](docs/features/for-loops.md)** - Iterate and expand with `for:` directive
- **[Matrix](docs/features/matrix.md)** - Generate combinations with `matrix:` directive
- **[Include](docs/features/include.md)** - Compose files with `include:` directive
- **[Document Expansion](docs/features/document-expansion.md)** - Root-level directives creating multiple documents

### Reference
- **[Syntax Reference](docs/syntax.md)** - Complete guide to all directives
- **[Custom Syntax Configuration](docs/custom-syntax.md)** - Configure directive keywords (Vue, Angular, or custom style)
- **[API Reference](docs/api.md)** - Complete API documentation
- **[Design Document](docs/DESIGN.md)** - Architecture and design decisions
- **[Development Guide](docs/DEVELOPMENT.md)** - Development workflow and feature implementation
- **[Test Coverage](docs/testing-coverage.md)** - Test coverage analysis

## Installation

```bash
go get github.com/titpetric/yamlexpr
```

## Quick Start

### Load and Evaluate YAML from a File

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/titpetric/yamlexpr"
)

func main() {
	// Create an Expr evaluator with the current directory as the filesystem
	expr := yamlexpr.New(os.DirFS("."))

	// Load and evaluate a YAML file
	// Returns a slice of Documents (may be multiple if using root-level for: or matrix:)
	docs, err := expr.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Process each resulting document
	for i, doc := range docs {
		fmt.Printf("Document %d: %+v\n", i, doc)
	}
}
```

### Parse and Evaluate YAML Data

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/titpetric/yamlexpr"
)

func main() {
	// Create an Expr evaluator
	expr := yamlexpr.New(os.DirFS("."))

	// Prepare your YAML data as Document (map[string]any)
	data := yamlexpr.Document{
		"name": "production",
		"env": map[string]any{
			"debug": false,
			"port":  8080,
		},
		"services": []any{
			map[string]any{"name": "api", "active": true},
			map[string]any{"name": "worker", "active": false},
		},
	}

	// Parse with variable interpolation, conditionals, and composition
	docs, err := expr.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	// Process each resulting document
	for i, doc := range docs {
		fmt.Printf("Document %d: %+v\n", i, doc)
	}
}
```

### Use Custom Directive Syntax

```go
// Use Vue.js-style directives
expr := yamlexpr.New(os.DirFS("."), yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:      "v-if",
	For:     "v-for",
	Include: "v-include",
	Matrix:  "v-matrix",
}))

// Now your YAML can use v-if, v-for, v-include, and v-matrix
docs, err := expr.Load("config.yaml")
```

See [Custom Syntax Configuration](docs/custom-syntax.md) for more examples.

## Features

- [X] **Variable Interpolation**: Use `${variable.path}` syntax in string values
- [X] **Conditionals**: Include/exclude blocks with `if:` directive
- [X] **For Loops**: Iterate and expand templates with `for:` directive
- [X] **Matrix Expansion**: Generate combinations with `matrix:` directive (with `exclude:` and `include:`)
- [X] **Composition**: Include external YAML files with `include:` directive
- [X] **Document Expansion**: Root-level directives create multiple output documents

## API

### Expr.Load(filename string) (map[string]any, error)

Loads a YAML file and evaluates all expressions in it. Files are resolved relative to the filesystem provided to `New()`.

```go
expr := yamlexpr.New(os.DirFS("."))
config, err := expr.Load("config.yaml")
```

### Expr.Process(doc any) (any, error)

Processes a YAML document (parsed as `map[string]any` or `[]any`) with expression evaluation.

```go
result, err := expr.Process(yamlData)
```

### Expr.ProcessWithStack(doc any, st *stack.Stack) (any, error)

Processes a YAML document with a custom variable stack for access to additional variables beyond the document's root keys.

```go
result, err := expr.ProcessWithStack(yamlData, customStack)
```

## Example YAML

```yaml
# config.yaml
environment: production
port: 8080

services:
  - for: "service in available_services"
    if: "${service.enabled}"
    name: "${service.name}"
    replicas: "${service.replicas}"

database:
  include: "database.yaml"
```

```yaml
# database.yaml
host: localhost
port: 5432
```
