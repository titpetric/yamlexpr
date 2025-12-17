# yamlexpr Quick Reference

Quick syntax guide for all yamlexpr features.

## Interpolation

```yaml
# Variable substitution
name: "Hello, ${user_name}!"
path: "/users/${user_id}/profile"

# Nested field access
connection: "postgres://${db.host}:${db.port}/mydb"

# Array access
first_item: ${items[0]}
```

**More:** [Interpolation docs](interpolation.md)

## Conditionals (`if:`)

```yaml
# Omit key when false
config:
  debug:
    if: ${enable_debug}
    level: verbose

# Works with expressions
server:
  autoscaling:
    if: replicas > 1 && enable_scaling
    min: 2
    max: 10

# Negation
development_only:
  if: !production
  enabled: true
```

**More:** [Conditionals docs](conditionals.md)

## For Loops

```yaml
# Simple iteration
servers:
  - for: server in server_list
    name: "${server}"

# With index
indexed:
  - for: (idx, item) in items
    position: ${idx}
    value: "${item}"

# Ignore value with underscore
positions:
  - for: (i, _) in items
    index: ${i}

# Filter with if
enabled:
  - for: svc in services
    if: ${svc.enabled}
    name: "${svc.name}"
```

**More:** [For loops docs](for-loops.md)

## Matrix

```yaml
# Generate combinations
- matrix:
    os: [linux, windows, macos]
    version: [18, 20, 22]
  name: "${os}-v${version}"

# Exclude combinations
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
  name: "${os}/${arch}"

# Add custom combinations
matrix:
  os: [linux]
  include:
    - os: macos
      arch: arm64
      xcode: "14"
  name: "${os}"
```

**More:** [Matrix docs](matrix.md)

## Include

```yaml
# Single file
include: "_base.yaml"
app:
  name: "myapp"

# In nested structure
database:
  include: "_db-config.yaml"
  pool_size: 20

# Multiple files
services:
  include:
    - "_monitoring.yaml"
    - "_logging.yaml"
  environment: production
```

**More:** [Include docs](include.md)

## Document Expansion

```yaml
# For loop at root level - creates multiple documents
for: env in [staging, production]
environment: "${env}"
services: []
---
# Results in 2 separate documents

# Matrix at root level - creates multiple documents  
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
job_name: "${os}-${arch}"
```

**More:** [Document expansion docs](document-expansion.md)

## Combined Features

```yaml
# For + If + Include
for: svc in services
if: ${svc.enabled}
name: "${svc.name}"
port: ${svc.port}
include: "_service-defaults.yaml"

# Matrix + Exclude + Include
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
  include:
    - os: macos
      arch: arm64
job: "${os}/${arch}"
```

## Common Patterns

### Reusable Defaults

```yaml
# _defaults.yaml
timeout: 30
retries: 3

# config.yaml
for: item in items
name: "${item}"
include: "_defaults.yaml"
```

### Layered Configuration

```yaml
include: "_base.yaml"              # Common settings
database:
  include: "_db-${env}.yaml"       # Environment-specific
  pool_size: ${db_pools[env]}      # Variable override
```

### Environment-Specific Services

```yaml
for: env in environments
environment: "${env}"
services:
  - for: svc in env.services
    if: ${svc.enabled}
    name: "${svc.name}"
    include: "_service-base.yaml"
```

## Syntax Rules

| Item       | Syntax                         | Notes                                   |
|------------|--------------------------------|-----------------------------------------|
| Variables  | `${var}`                       | Required `${}` for YAML parsing         |
| Nested     | `${a.b.c}`                     | Dot notation for nested access          |
| Arrays     | `${items[0]}`                  | Bracket notation for indices            |
| If         | `if: condition`                | Can be boolean, variable, or expression |
| For        | `for: var in array`            | Creates array iterations                |
| For Index  | `for: (idx, var) in array`     | Tuple unpacking                         |
| Underscore | `(_, val)` or `(idx, _)`       | Ignore index or value                   |
| Matrix     | `matrix: {a: [...], b: [...]}` | Cartesian product                       |
| Exclude    | `exclude: [{a: x, b: y}]`      | Filter matrix combinations              |
| Include    | `include: "file.yaml"`         | Single or array of files                |

## Expression Operators

```yaml
# Comparison
if: count > 5
if: name == "admin"
if: status != "failed"
if: value >= 100

# Logical
if: a && b          # AND
if: a || b          # OR
if: !flag           # NOT

# Field access
if: item.enabled
if: config.server.debug
if: users[0].active
```

## Tips

- **Quotes**: Use quotes for `if:` and `for:` directives to ensure YAML parsing
- **Variables**: All root-level keys become available as variables
- **Order**: Includes are processed first, then for/matrix, then if
- **Nesting**: Features can be deeply nested; use whitespace for readability
- **Files**: Include paths are relative to the base directory passed to `Expr.New()`

EOF
