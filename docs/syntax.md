# yamlexpr Syntax Reference

This document describes the special YAML directives supported by yamlexpr: `include`, `for`, and `if`.

## Variable Interpolation

Any string value can include variable substitution using `${variable}` syntax:

```yaml
message: "Hello, ${name}!"
path: "/users/${user_id}/profile"
```

Variables are resolved from the stack context. You can access nested values using dot notation:

```yaml
greeting: "Welcome ${user.name} to ${application.title}"
```

## Include Directive

The `include` directive pulls in external YAML files and merges them into the current document.

### Single File

```yaml
include: "_base-config.yaml"
debug: true
```

The base config is merged, then the `debug: true` key is added.

### Multiple Files

```yaml
include:
  - "_base-config.yaml"
  - "_database.yaml"
logging:
  level: "info"
```

Files are resolved relative to the base directory provided to `Expr.New()`.

## For Loop

The `for` directive expands an array by iterating over values and creating multiple items.

### Simple Loop

Iterate over array values, available as `${item}`:

```yaml
users:
  - for: ["alice", "bob", "charlie"]
    name: "${item}"
    active: true
```

Result:
```yaml
users:
  - name: "alice"
    active: true
  - name: "bob"
    active: true
  - name: "charlie"
    active: true
```

### Looping Over Inline Arrays

You can also reference a variable or inline array literal:

```yaml
items:
  - for: ${product_list}
    sku: "${item.id}"
    price: ${item.cost}
```

### Loop with Index

Use the `(idx, item) in` syntax to access both index and value:

```yaml
products: ["apple", "banana", "cherry"]
items:
  - for: (idx, item) in products
    index_str: "${idx}"
    value: "${item}"
```

Result:
```yaml
products: ["apple", "banana", "cherry"]
items:
  - index_str: "0"
    value: "apple"
  - index_str: "1"
    value: "banana"
  - index_str: "2"
    value: "cherry"
```

### Ignoring Index or Item

Use underscore `_` to skip a variable:

```yaml
names: ["alice", "bob", "charlie"]
output:
  - for: (_, name) in names
    person: "${name}"
```

### Filtering with If

Combine `for` with `if` to filter items:

```yaml
enabled_services:
  - for:
      - name: "api"
        active: true
      - name: "worker"
        active: false
      - name: "scheduler"
        active: true
    if: item.active
    service: "${item.name}"
```

Result:
```yaml
enabled_services:
  - service: "api"
  - service: "scheduler"
```

Only items where `item.active` is true are included.

### Empty Array

Looping over an empty array produces no output:

```yaml
items:
  - for: []
```

Result:
```yaml
items: []
```

## If Conditional

The `if` directive includes or omits a key based on a boolean condition.

### Omit Key on False

```yaml
config:
  name: "production"
  debug:
    if: false
    enabled: true
  version: "1.0"
```

Result:
```yaml
config:
  name: "production"
  version: "1.0"
```

The `debug` key is removed entirely because `if: false`.

### Include Key on True

```yaml
config:
  name: "production"
  debug:
    if: true
    enabled: true
    level: "verbose"
  version: "1.0"
```

Result:
```yaml
config:
  name: "production"
  debug:
    enabled: true
    level: "verbose"
  version: "1.0"
```

When `if: true`, the `if` directive itself is removed and the remaining keys are included.

### Condition Expressions

`if` conditions are evaluated as boolean expressions:

- Boolean values: `if: true`, `if: false`
- Variable references: `if: ${enable_feature}`
- Expressions using the [expr-lang library](https://github.com/expr-lang/expr):
  - `if: item.active` (field access)
  - `if: count > 5` (comparisons)
  - `if: name != "admin"` (equality)
  - `if: !disabled` (negation)
  - `if: status == "active" && verified` (logic operators)

### If with For

When `if` appears alongside `for`, the condition filters each item:

```yaml
enabled_services:
  - for:
      - name: "api"
        active: true
      - name: "worker"
        active: false
    if: item.active
    service: "${item.name}"
```

Only items matching the condition are expanded.

### If with Nested Keys

If the `if` key is on a map, that entire map is included/omitted:

```yaml
config:
  database:
    if: ${use_postgres}
    host: "localhost"
    port: 5432
```

If `use_postgres` is false, the entire `database` key is removed.

## Combined Features

Features can be combined in a single YAML document:

```yaml
include: "_base.yaml"

services:
  - for: ${service_list}
    if: item.enabled
    name: "${item.name}"
    config:
      include: "_service-defaults.yaml"
      port: ${item.port}
```

This example:
1. Includes a base config file
2. Loops over a service list
3. Filters services based on the `enabled` flag
4. For each service, includes default settings and applies service-specific port
