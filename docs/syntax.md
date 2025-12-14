# yamlexpr Syntax Reference

yamlexpr provides several directives and features for YAML composition, conditional evaluation, and variable interpolation.

## Variable Interpolation with `${...}`

Any string value can include variable substitution and expression evaluation using `${...}` syntax:

### Simple Variable References

```yaml
message: "Hello, ${name}!"
path: "/users/${user_id}/profile"
config_file: ${config_path}
```

### Expression Evaluation

The `${}` syntax supports full expr-lang expressions, enabling arithmetic, string functions, and comparisons:

**Arithmetic operations:**

```yaml
price: 100
items:
  - for: [1, 2, 3]
    quantity: ${item}
    total: ${item * price}
    discounted: ${item * price * 0.9}
```

**String functions:**

```yaml
items:
  - for: ["hello", "world"]
    word: ${item}
    uppercase: ${upper(item)}
    length: ${len(item)}
```

**Comparisons:**

```yaml
items:
  - for: [1, 5, 10, 15]
    num: ${item}
    is_large: ${item > 10}
    is_even: ${item % 2 == 0}
```

Variables are resolved from the stack (root-level keys or passed variables). Undefined variables produce an error.

**Note:** The `${}` syntax was chosen to enable proper YAML parsing. Without the `$` sign, `{variable}` causes parsing issues where the parser expects an object.

**Inspiration:** This syntax is inspired by GitHub Actions.

## Directives (Handlers)

Directives are special YAML keys that control document processing. They can be combined and nested for powerful composition patterns.

### `embed` - File Composition

The `embed` directive loads external YAML files and merges their content into the current structure.

**Basic usage:**

```yaml
embed: "base.yaml"
```

**Multiple files:**

```yaml
embed:
  - "base.yaml"
  - "overrides.yaml"
```

Files are resolved relative to the base directory provided to `yamlexpr.New()`. Files are processed in order and merged together.

**Example:**

```yaml
# config.yaml
embed: "_defaults.yaml"

settings:
  debug: true
  port: 8080
```

The resulting document will have all keys from `_defaults.yaml` merged with the keys from `config.yaml`, with `config.yaml` taking precedence.

### `for` - Loop Expansion

The `for` directive expands an array by iterating over values and creating multiple items.

**Basic syntax:**

- `for: item in items` - Iterate over items with single variable
- `for: (idx, item) in items` - Iterate with both index and item
- `for: [value1, value2, ...]` - Direct array (uses default `item` variable)

**Single variable:**

```yaml
steps:
  - for: step in workflow_steps
    name: "${step.name}"
    run: "${step.command}"
```

**Direct array syntax:**

```yaml
users:
  - for: ["alice", "bob", "charlie"]
    name: "${item}"
    active: true
```

When using direct array syntax, the loop variable is automatically named `item`.

**With index and value:**

```yaml
items:
  - for: (idx, item) in products
    position: "${idx}"
    name: "${item.name}"
    price: "${item.price}"
```

**Omitting variables with `_`:**

```yaml
items:
  - for: (_, item) in products
    name: "${item.name}"
```

### `if` - Conditional Inclusion

The `if` directive includes or omits a block based on a boolean condition.

**Boolean values:**

```yaml
config:
  debug:
    if: true
    level: "verbose"
  production:
    if: false
    optimized: true
```

When `if: false`, the key is removed from the output entirely.

**Variable references:**

```yaml
config:
  experimental_feature:
    if: ${enable_experimental}
    description: "This is an experimental feature"
```

**Expression-based conditions:**

Using [expr-lang library](https://github.com/expr-lang/expr) for complex expressions:

```yaml
services:
  - for: item in services
    if: item.enabled
    name: "${item.name}"
  
  - for: item in services
    if: item.port > 1024
    name: "${item.name}"
    port: "${item.port}"

  - for: item in services
    if: item.status == "active" && item.verified
    name: "${item.name}"
```

**Supported expressions:**

- Field access: `item.active`
- Comparisons: `count > 5`, `status == "active"`, `port != 8080`
- Logical operators: `&&`, `||`, `!`
- Arithmetic: `count + 1`, `price * 0.9`
- Function calls: `len(items) > 0`

### `discard` - Inverse Conditional

The `discard` directive is the opposite of `if` - it removes a block when true.

```yaml
config:
  deprecated_key:
    discard: true
    value: "this will be removed"
  
  active_key:
    discard: false
    value: "this will be kept"
```

This is useful when you want to express logic as "remove if X" rather than "keep if not X".

### `matrix` - Cartesian Product Expansion

The `matrix` directive expands jobs or configurations across multiple dimensions (inspired by GitHub Actions).

**Basic usage:**

```yaml
jobs:
  test:
    matrix:
      os: [linux, macos, windows]
      version: [1.20, 1.21]
    runs_on: "${matrix.os}"
    go_version: "${matrix.version}"
```

Expands to 6 job combinations (3 OS × 2 versions).

**With exclude:**

```yaml
matrix:
  os: [linux, macos, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
```

Removes specific combinations.

**With include:**

```yaml
matrix:
  os: [linux, windows]
  include:
    - os: macos
      arch: arm64
      xcode: "14"
```

Adds extra combinations that don't fit the cartesian product.

## Combined Features

Directives can be combined in a single YAML document for powerful composition patterns:

```yaml
embed: "_base.yaml"

services:
  - for: item in service_list
    if: item.enabled
    name: "${item.name}"
    image: "${item.image}"
    ports:
      - for: port in item.ports
        port: "${port}"
    config:
      embed: "_service-defaults.yaml"
      environment:
        SERVICE_PORT: "${item.port}"
        SERVICE_NAME: "${item.name}"
    healthcheck:
      if: ${item.enable_healthcheck}
      path: "${item.healthcheck_path}"
```

This example:

1. Embeds a base configuration
2. Loops over a service list
3. Filters services based on the `enabled` flag
4. For each service, embeds default settings
5. Expands ports with another nested loop
6. Conditionally includes health check configuration

## Directive Evaluation Order

When multiple directives are present in a single map, they are evaluated in this order:

1. `embed` (priority 1000) - File merging happens first
2. `for` (priority 100) - Loop expansion before conditionals
3. `if` (priority 100) - Conditional evaluation
4. Custom handlers - User-defined handlers
5. Regular key processing - Normal YAML keys

This order ensures that:

- Embedded files are merged before other processing
- Loop variables are in scope when conditionals are evaluated
- Interpolation is applied to all values automatically

## Type Support

All directives work with:

- **Strings:** Variable interpolation
- **Numbers:** Condition evaluation, loop indices
- **Booleans:** Direct condition values
- **Arrays:** Loop expansion, multiple embeds
- **Objects:** Nested directive application
- **Null:** Treated as false in conditions

## Error Handling

- **Undefined variables:** Strict mode errors; lenient mode leaves placeholder unchanged
- **Invalid syntax:** Parse errors with context (path, line info)
- **Missing files:** Embed errors with filename
- **Type mismatches:** Clear error messages explaining expected vs actual type
- **Cyclic includes:** Detected and reported to prevent infinite loops
