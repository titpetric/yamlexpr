# yamlexpr

[![Go Reference](https://pkg.go.dev/badge/github.com/titpetric/yamlexpr.svg)](https://pkg.go.dev/github.com/titpetric/yamlexpr) [![Go Report Card](https://goreportcard.com/badge/github.com/titpetric/yamlexpr)](https://goreportcard.com/report/github.com/titpetric/yamlexpr) [![PkgGoDev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/titpetric/yamlexpr) [![Test Coverage](https://img.shields.io/badge/coverage-82.26%25-green)](docs/testing-coverage.md)

YAML composition, interpolation, and conditional evaluation for Go.

## Documentation

- **[Tutorial](docs/tutorial.md)** - Comprehensive guide with real-world examples
- **[Syntax Reference](docs/syntax.md)** - Complete guide to `include`, `for`, and `if` directives
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
	result, err := expr.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Result is map[string]any containing the processed YAML
	fmt.Printf("%+v\n", result)
}
```

### Process YAML Data with Variables

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/titpetric/yamlexpr"
	"github.com/titpetric/yamlexpr/stack"
)

func main() {
	// Create an Expr evaluator
	expr := yamlexpr.New(os.DirFS("."))

	// Prepare your YAML data
	data := map[string]any{
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

	// Process with variable interpolation, conditionals, and composition
	result, err := expr.Process(data)
	if err != nil {
		log.Fatal(err)
	}

	// Result is map[string]any with all expressions evaluated
	fmt.Printf("%+v\n", result)
}
```

### Process with Custom Variables

```go
// Create a custom variable stack
st := stack.New(map[string]any{
	"version": "1.0.0",
	"region":  "us-west-2",
})

expr := yamlexpr.New(os.DirFS("."))
result, err := expr.ProcessWithStack(data, st)
if err != nil {
	log.Fatal(err)
}
```

### Use Custom Directive Syntax

```go
// Use Vue.js-style directives
expr := yamlexpr.New(os.DirFS("."), yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:      "v-if",
	For:     "v-for",
	Include: "v-include",
}))

// Now your YAML can use v-if, v-for, and v-include
result, err := expr.Load("config.yaml")
```

See [Custom Syntax Configuration](docs/custom-syntax.md) for more examples.

## Features

- [X] **Variable Interpolation**: Use `${variable.path}` syntax in string values
- [X] **Conditionals**: Include/exclude blocks with `if:` directive
- [X] **For Loops**: Iterate and expand templates with `for:` directive
- [X] **Composition**: Include external YAML files with `include:` directive

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
