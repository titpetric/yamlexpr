# yamlexpr Documentation

Complete documentation for yamlexpr - YAML composition, interpolation, and conditional evaluation.

## Getting Started

- **[Tutorial](tutorial.md)** - Comprehensive guide with real-world examples
- **[Quick Start](../README.md#quick-start)** - Get up and running in minutes

## Reference Documentation

### Directives & Handlers

yamlexpr provides powerful directives for YAML composition and transformation:

- **[Include (`include:`)](/docs/handlers/include.md)** - Load and merge external YAML files
- **[For Loops (`for:`)](/docs/handlers/for.md)** - Iterate over collections with loop variables
- **[Conditionals (`if:`)](/docs/handlers/if.md)** - Include/exclude blocks based on conditions
- **[Discard (`discard:`)](/docs/handlers/discard.md)** - Omit blocks conditionally (inverse of `if:`)
- **[Matrix (`matrix:`)](/docs/handlers/matrix.md)** - Generate combinations using cartesian product
- **[Interpolation (`${}`)](/docs/handlers/interpolation.md)** - Substitute variables and expressions

### Core Topics

- **[Syntax Reference](syntax.md)** - Complete guide to all directives and syntax
- **[API Reference](api.md)** - Complete API documentation and examples
- **[Custom Syntax Configuration](custom-syntax.md)** - Configure directive keywords (Vue, Angular, or custom)
- **[Design Document](DESIGN.md)** - Architecture and design decisions

## Development

- **[Development Guide](DEVELOPMENT.md)** - Development workflow and feature implementation
- **[Test Coverage](testing-coverage.md)** - Test coverage analysis and metrics

## Directory Structure

```
docs/
├── README.md (this file)
├── tutorial.md
├── syntax.md
├── api.md
├── custom-syntax.md
├── DESIGN.md
├── DEVELOPMENT.md
├── testing-coverage.md
├── handlers/
│   ├── include.md
│   ├── for.md
│   ├── if.md
│   ├── discard.md
│   ├── matrix.md
│   └── interpolation.md
```

## Handler Priority

Handlers are executed in order of priority (highest first):

| Priority | Handler         | Purpose                             |
|----------|-----------------|-------------------------------------|
| 1000     | `embed:`        | Load and merge files (always first) |
| 100      | `for:`          | Loop expansion                      |
| 100      | `if:`           | Conditional inclusion               |
| 50       | `interpolation` | Variable substitution (automatic)   |
| 10       | `discard:`      | Conditional omission                |
| 5        | `matrix:`       | Cartesian product expansion         |

## Quick Reference

### Variable Interpolation

Use `${}` syntax in string values to substitute variables and expressions:

```yaml
message: "Hello ${name}"
port: ${server.port}
doubled: ${count * 2}
```

See: [Interpolation Handler](/docs/handlers/interpolation.md)

### File Composition

Load and merge external YAML files:

```yaml
include: "base.yaml"
config:
  key: "value"
```

See: [Include Handler](/docs/handlers/include.md)

### Loop Expansion

Iterate over collections:

```yaml
items:
  - for: "item in items"
    name: "${item}"
```

See: [For Loop Handler](/docs/handlers/for.md)

### Conditional Blocks

Include/exclude blocks based on conditions:

```yaml
block:
  if: "${condition}"
  value: "only included if condition is true"
```

See: [If Handler](/docs/handlers/if.md)

### Cartesian Product

Generate combinations from dimensions:

```yaml
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
```

See: [Matrix Handler](/docs/handlers/matrix.md)

## Installation & Setup

```bash
go get github.com/titpetric/yamlexpr
```

## Examples

### Basic YAML Processing

```go
expr := yamlexpr.New(yamlexpr.WithFS(os.DirFS(".")))
docs, err := expr.Load("config.yaml")
```

### With Variables

```go
st := stack.New(map[string]any{
	"env":   "production",
	"debug": false,
})
result, err := expr.ProcessWithStack(data, st)
```

See: [API Reference](api.md) for complete API documentation

## Common Patterns

### Configuration with Environment Overrides

```yaml
# base.yaml
database:
  host: "localhost"
  port: 5432

# production.yaml
include: "base.yaml"
database:
  host: "db.prod.example.com"
```

### Dynamic Service Configuration

```yaml
services:
  - for: "service in service_list"
    if: "${service.enabled}"
    name: "${service.name}"
    port: "${service.port}"
```

### Test Matrix

```yaml
matrix:
  os: [ubuntu, windows, macos]
  node: [18, 20, 22]
  exclude:
    - os: windows
      node: 18

tests:
  run: "npm test"
  os: "${os}"
  node: "${node}"
```

## Troubleshooting

### Undefined Variable Error

```
undefined variable 'varname' at path.to.key
```

**Solution:** Ensure the variable is defined in the root document or stack.

### Invalid Expression

```
error evaluating expression 'item * * 2'
```

**Solution:** Check expression syntax. Use valid operators and parentheses.

### File Not Found

```
error including "missing.yaml": file does not exist
```

**Solution:** Verify file path is relative to the filesystem root provided to `yamlexpr.New()`.

## Contributing

Documentation is kept up-to-date with the implementation. When adding features:

1. Update relevant handler documentation in `docs/handlers/`
2. Update syntax reference in `docs/syntax.md`
3. Add examples to `docs/tutorial.md`
4. Update `docs/api.md` with new API methods

## See Also

- [GitHub Repository](https://github.com/titpetric/yamlexpr)
- [Issue Tracker](https://github.com/titpetric/yamlexpr/issues)
