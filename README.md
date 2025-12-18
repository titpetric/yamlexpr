# yamlexpr

[![Go Reference](https://pkg.go.dev/badge/github.com/titpetric/yamlexpr.svg)](https://pkg.go.dev/github.com/titpetric/yamlexpr) [![Go Report Card](https://goreportcard.com/badge/github.com/titpetric/yamlexpr)](https://goreportcard.com/report/github.com/titpetric/yamlexpr) [![PkgGoDev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/titpetric/yamlexpr) [![Test Coverage](https://img.shields.io/badge/coverage-40.1%25-yellowgreen)](docs/testing-coverage.md)

YAML composition, interpolation, and conditional evaluation for Go.

## Example YAML

```yaml
matrix:                                     # Document expansion on root level
  environment: [production, development]    # Iteration dimensions.

environment: ${environment}                 # Interpolation from matrix vars.
port: 8080

services:
  - for: service in available_services      # Loops
    if: "${service.enabled}"                # Conditions
    name: "${service.name}"                 # Interpolation
    replicas: "${service.replicas}"

database:
  include: "database-${environment}.yaml"   # Composition
```

```yaml
# database-development.yaml
host: localhost
port: 5432
```

The main goals of yamlexpr is to be a lightweight data traversal engine,
upon which new functionality can be built. It evaluates data and
produces a deterministic plain yaml output. That output can be used
further to provide execution pipelines, documentation and other uses.

## Features

- [X] **Variable Interpolation**: Use `${variable.path}` syntax in string values
- [X] **Conditionals**: Include/exclude blocks with `if:` directive
- [X] **For Loops**: Iterate and expand templates with `for:` directive
- [X] **Matrix Expansion**: Generate combinations with `matrix:` directive (with `exclude:` and `include:`)
- [X] **Composition**: Include external YAML files with `include:` directive
- [X] **Document Expansion**: Root-level directives create multiple output documents

### Getting Started
- **[Tutorial](docs/tutorial.md)** - Comprehensive guide with real-world examples
- **[Quick Reference](docs/features/)** - Syntax cheat sheet and common patterns

### Feature Documentation
- **[Interpolation](docs/features/interpolation.md)** - Variable substitution with `${variable}` syntax
- **[Conditionals](docs/features/conditionals.md)** - Include/exclude with `if:` directive
- **[For Loops](docs/features/for-loops.md)** - Iterate and expand with `for:` directive
- **[Matrix](docs/features/matrix.md)** - Generate combinations with `matrix:` directive
- **[Include](docs/features/include.md)** - Compose files with `include:` directive
- **[Document Expansion](docs/features/document-expansion.md)** - Root-level directives creating multiple documents

### Reference
- **[Syntax Reference](docs/features/)** - Complete guide to all directives
- **[API Reference](docs/api.md)** - Complete API documentation
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
