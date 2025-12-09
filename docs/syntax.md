# yamlexpr Syntax Reference

This document describes the special YAML directives supported by yamlexpr: `embed`, `for`, and `if`.

## Variable Interpolation with `${variable}`

Any string value can include variable substitution using `${variable}` syntax:

```yaml
message: "Hello, ${name}!"
path: "/users/${user_id}/profile"
```

This was chosen to enable yaml parsing for `value: ${variable}` ; omitting the `$` sign causes a parsing issue, where the parser starts to expect an object due to `{}`.

This is inspired by GitHub actions.

## Embedding files with `embed`

To enable composition, you can use the `embed` statement at any level of the YAML.

```yaml
embed: _base.yaml

settings:
  embed: _settings.yml
```

Files are resolved relative to the base directory provided to `yamlexpr.New()`.

## Looping with `for`

The `for` directive expands an array by iterating over values and creating multiple items.

Loops allow the following syntax:

- `for: item in items`
- `for: (index, item) in items`

Use the syntax to access both index and value as needed:

```yaml
items:
  - for: (idx, item) in products
    index_str: "${idx}"
    value: "${item}"
```

### Filtering with If

Combine `for` with `if` to filter items:

```yaml
enabled_services:
  - for: item in items
    if: item.active
    service: "${item.name}"
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
embed: "_base.yaml"

services:
   - for: item in service_list
     if: item.enabled
     name: "${item.name}"
     config:
       embed: "_service-defaults.yaml"
       port: ${item.port}
```

This example:

1. Embeds a base config file
2. Loops over a service list
3. Filters services based on the `enabled` flag
4. For each service, embeds default settings and applies service-specific port
