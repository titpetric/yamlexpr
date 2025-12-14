# Discard Handler (`discard:`)

The **`discard:` directive** conditionally omits blocks based on a boolean condition, opposite of `if:`.

## Overview

The `discard:` directive omits a block when set to true (inverse of the `if:` directive). It provides a convenient way to express "exclude when condition is true" logic.

## Syntax

```yaml
key:
  discard: <condition>
  # other keys...
```

## Condition Types

### Boolean Literals

```yaml
active_item:
  discard: false      # Included
  name: "active"

inactive_item:
  discard: true       # Omitted
  name: "inactive"
```

### String Boolean Representations

```yaml
included:
  discard: "false"
  value: 1

excluded:
  discard: "true"
  value: 2

# Also accepts numeric strings
zero_excluded:
  discard: "0"
  value: 3

one_excluded:
  discard: "1"
  value: 4
```

### Numeric Values

```yaml
item1:
  discard: 0          # Included (falsy)
  name: "zero"

item2:
  discard: 1          # Omitted (truthy)
  name: "non-zero"
```

### Null/Nil

```yaml
item:
  discard: null       # Included (falsy)
  name: "present"
```

## Examples

### Basic Discard

```yaml
# Input
testing: true

config:
  discard: false
  debug: true
  
skip_this:
  discard: true
  debug: false

conditional:
  discard: "${testing}"
  test_mode: true

# Output
config:
  debug: true

# skip_this and conditional are omitted
```

### Discard with Conditions

```yaml
# Input
environment: "production"
is_dev: false

services:
  dev_tools:
    discard: "${environment != 'dev'}"
    debug_port: 9000
  
  dev_server:
    discard: "${is_dev == false}"
    port: 3000

# Output (if environment is "production" and is_dev is false)
services:
  # Both are omitted
```

### Filtering Arrays

```yaml
# Input
users:
  - name: "Alice"
    disabled: false
  - name: "Bob"
    disabled: true
  - name: "Charlie"
    disabled: false

active_users:
  - for: "user in users"
    discard: "${user.disabled}"
    name: "${user.name}"

# Output
active_users:
  - name: "Alice"
  - name: "Charlie"
```

### Conditional Feature Disabling

```yaml
# Input
feature_flags:
  legacy_api: true
  analytics: false

endpoints:
  v1:
    discard: "${feature_flags.legacy_api}"
    path: "/api/v1"
  
  v2:
    discard: false
    path: "/api/v2"
  
  analytics:
    discard: "${feature_flags.analytics}"
    path: "/analytics"

# Output
endpoints:
  v2:
    path: "/api/v2"
```

## API Functions

### Discard Handler

**`DiscardHandlerBuiltin() DirectiveHandler`**

Creates a handler for the `discard:` directive.

```go
handler := DiscardHandlerBuiltin()
// Use with: expr.WithDirectiveHandler("discard", handler, 10)
```

**Logic:**
- `true` or truthy: Omits block (returns `nil, true, nil`)
- `false` or falsy: Continues processing (returns `nil, false, nil`)

## Behavior

1. **Evaluation**: Condition is evaluated to boolean
2. **True Result**: Entire block is omitted (not included in output)
3. **False Result**: Block is processed normally
4. **Null Output**: If condition is true, the key produces no output

## Truthiness Rules

Values are evaluated using standard truthiness rules:

### Truthy Values (Block Omitted)
- Non-zero integers: `1`, `42`
- Non-zero floats: `1.5`, `0.001`
- String `"true"`, `"1"`, `"yes"`
- Non-empty strings: `"text"`, `"false"` (string "false" is truthy)
- Non-empty arrays: `[1]`
- Non-empty objects: `{a: 1}`

### Falsy Values (Block Included)
- Zero: `0`, `0.0`
- String `"false"`, `"0"`, `"no"`, `""`
- `false` boolean
- `null` / `nil`
- Empty arrays: `[]`
- Empty objects: `{}`

**Note:** String `"false"` is **truthy** (non-empty string). Use `discard: "${value == false}"` for boolean false.

## Comparison with `if:`

- **`if: condition`**: Includes block if **true**, omits if **false**
- **`discard: condition`**: Omits block if **true**, includes if **false**

```yaml
# Equivalent
block1:
  if: false
  value: 1

block2:
  discard: true
  value: 1

# Both omit their blocks
```

```yaml
# Also equivalent
active1:
  if: true
  enabled: true

active2:
  discard: false
  enabled: true

# Both include their blocks
```

## Error Handling

### Invalid Type

```yaml
discard: []
# Error: discard must be boolean, got []
```

### Invalid String

```yaml
discard: "invalid"
# Error: discard value must be boolean or 'true'/'false', got string 'invalid'
```

### Type Errors

```yaml
discard: {key: "value"}
# Error: discard must be boolean, got map[string]any
```

## Use Cases

### Deprecated Features

```yaml
deprecated:
  discard: true
  warning: "This feature is deprecated"

current:
  discard: false
  enabled: true
```

### Conditional Test Skipping

```yaml
tests:
  - name: "slow_test"
    discard: "${skip_slow_tests}"
    timeout: 300
  
  - name: "fast_test"
    discard: false
    timeout: 5
```

### Environment-Specific Exclusion

```yaml
environment: "production"

development_tools:
  discard: "${environment != 'development'}"
  debug_panel: true
  hot_reload: true

metrics:
  discard: "${environment == 'test'}"
  enabled: true
  endpoint: "https://metrics.example.com"
```

### Permission-Based Feature Exclusion

```yaml
user_role: "guest"

admin_panel:
  discard: "${user_role != 'admin'}"
  path: "/admin"

premium_features:
  discard: "${!user_role.has_premium}"
  features: ["advanced_analytics", "api_access"]
```

## Performance Considerations

- Conditions are evaluated once per block
- Simple boolean values have minimal overhead
- String-to-boolean conversion is cached
- No impact compared to `if:` with inverted condition

## Priority

The `discard:` handler has **priority 10**, running:
- After `include:` (1000) - files are already merged
- After `for:` and `if:` (100) - other directives are processed first
- Before custom handlers (default 0)

## See Also

- [Conditionals (`if:`)](/docs/handlers/if.md) - inverse of discard
- [For Loops (`for:`)](/docs/handlers/for.md) - iterate with discard filtering
- [Interpolation (`${}`)](/docs/handlers/interpolation.md) - variable substitution in discard
- [Syntax Reference](/docs/syntax.md) - complete directive syntax guide
