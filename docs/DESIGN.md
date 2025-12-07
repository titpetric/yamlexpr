# yamlexpr Design Document

## Purpose

YAML is a data format. `yamlexpr` adds evaluation capabilities to YAML: variable interpolation, conditionals, file composition, and loop expansion. This allows YAML to be used as a simple configuration or template format where values can be computed from variables.

## Example

Input YAML with variables:
```yaml
include: "_base.yaml"
database:
  host: ${db_host}
  port: ${db_port}

servers:
  - for: ${server_list}
    if: item.enabled
    name: ${item.name}
    ip: ${item.ip}

debug:
  if: ${enable_debug}
  level: "verbose"
```

Variables supplied:
```go
stack := stack.New(map[string]any{
  "db_host": "localhost",
  "db_port": 5432,
  "server_list": []map[string]any{
    {"name": "web-1", "ip": "10.0.1.1", "enabled": true},
    {"name": "web-2", "ip": "10.0.1.2", "enabled": false},
  },
  "enable_debug": true,
})

expr := yamlexpr.New(fs, stack)
output, err := expr.Evaluate(input)
```

Output: YAML with includes merged, interpolations resolved, loops expanded, and conditions evaluated. For syntax details, see [Syntax Reference](syntax.md).

## Code Organization

```
yamlexpr/
├── stack/           # Variable scoping and resolution
│   ├── stack.go     # Stack type for variable lookup
│   └── stack_test.go
├── expr.go          # Main Expr type for processing YAML documents
├── context.go       # ExprContext carries evaluation state through processing
├── interpolate.go   # Handles ${var} substitution in strings
├── for_loop.go      # Expands for: directives
├── util.go          # Helper functions
└── testdata/fixtures/ # Test fixtures with input.yaml and input.yaml.expected
```

## Architecture

**Stack package**: Variable scoping with path resolution (e.g., `user.name`, `items[0]`). Importable by any project needing variable lookup.

**Root package**: YAML document processing. Depends on stack for variable resolution. Handles:
- Parsing YAML into maps/slices
- Processing include, for, if directives
- Interpolating ${} syntax in strings
- Maintaining ExprContext through the document tree for error reporting and path tracking

**ExprContext**: Carries evaluation state (stack, current path, include chain) through recursive processing functions.

## Testing

Tests use black box style with exported APIs only. Each module has corresponding _test.go file:
- `stack_test.go`: Stack API
- `expr_test.go`: Expr type and main Evaluate() function
- `context_test.go`: ExprContext API
- `interpolate_test.go`: Interpolation behavior
- `expr_fixtures_test.go`: Fixture-based integration tests

Fixtures are the source of truth. Format:
```yaml
# input.yaml
items:
  - name: ${item_name}
---
# input.yaml.expected
items:
  - name: "resolved_value"
```

Run: `go test ./...`

## Dependencies

- `gopkg.in/yaml.v3`: YAML parsing and serialization
- `github.com/expr-lang/expr`: Expression evaluation for if conditions
- `github.com/stretchr/testify`: Test assertions
