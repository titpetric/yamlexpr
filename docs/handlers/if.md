# Conditional Handler (`if:`)

The **`if:` directive** conditionally includes or excludes blocks based on evaluated conditions.

## Overview

The `if:` directive evaluates a condition and omits the entire block if the condition is false. This enables dynamic configuration based on variables, expressions, and environment state.

## Syntax

```yaml
key:
  if: <condition>
  # other keys...
```

## Condition Types

### Boolean Literals

```yaml
item:
  if: true         # Always included
  name: "active"

item2:
  if: false        # Always omitted
  name: "inactive"
```

### String Boolean Representations

```yaml
item1:
  if: "true"       # Included
  name: "yes"

item2:
  if: "false"      # Omitted
  name: "no"

item3:
  if: "1"          # Included (truthy)
  name: "one"

item4:
  if: "0"          # Omitted (falsy)
  name: "zero"
```

### Variable References (Interpolated)

```yaml
debug: true
server:
  if: "${debug}"
  log_level: "debug"

status: "active"
service:
  if: "${status == 'active'}"
  enabled: true
```

### Direct Variable Paths

```yaml
user:
  name: "Alice"
  admin: true

config:
  if: "user.admin"         # Simple path
  admin_panel: "/admin"
```

### Expression Evaluation

Full [expr-lang](https://github.com/expr-lang/expr) expressions are supported:

**Comparisons:**

```yaml
count: 10
items:
  - if: "${count > 5}"
    name: "many"
  
  - if: "${count <= 5}"
    name: "few"

status: "active"
service:
  if: "${status == 'active'}"
  running: true
```

**Arithmetic:**

```yaml
price: 100
discount:
  if: "${price > 50}"
  percent: 10
```

**Logical Operators:**

```yaml
user:
  name: "Alice"
  verified: true
  premium: false

access:
  if: "${user.verified && user.premium}"
  features: ["advanced", "priority_support"]
```

**Array/Object Operations:**

```yaml
items: [1, 2, 3, 4, 5]

summary:
  if: "${len(items) > 0}"
  count: 5
  has_data: true
```

## API Functions

### Condition Evaluation

**`EvaluateConditionWithPath(condition any, stack *Stack, path string) (bool, error)`**

Evaluates an `if:` condition and returns the boolean result.

```go
result, err := EvaluateConditionWithPath(condition, stack, "services.database.if")
// Returns: true/false, or error if evaluation fails
```

**Parameters:**
- `condition`: The condition value (bool, string, int, etc.)
- `stack`: Variable stack for resolution
- `path`: Dot-separated path for error context

**Returns:**
- `bool`: Evaluation result
- `error`: Parsing or evaluation errors with context

### Truthiness Evaluation

**`IsTruthy(v any) bool`**

Determines if a value is truthy (non-zero, non-empty, non-nil).

```go
IsTruthy(true)     // true
IsTruthy(1)        // true
IsTruthy(0)        // false
IsTruthy("")       // false
IsTruthy("text")   // true
IsTruthy([]any{})  // false
IsTruthy([]any{1}) // true
```

### String Utilities

**`IsQuoted(s string) bool`**

Checks if a string is already quoted (single or double quotes).

```go
IsQuoted("'text'") // true
IsQuoted(`"text"`) // true
IsQuoted("text")   // false
```

## Examples

### Basic Conditional

```yaml
# Input
debug: true
environment: "production"

logging:
  if: "${debug}"
  level: "debug"
  
cache:
  if: "${environment == 'production'}"
  enabled: true

---

# Output
logging:
  level: "debug"

cache:
  enabled: true
```

### Conditional with For Loop

```yaml
# Input
services:
  - name: "api"
    enabled: true
  - name: "worker"
    enabled: false
  - name: "cache"
    enabled: true

config:
  for: "service in services"
  if: "${service.enabled}"
  service_name: "${service.name}"

---

# Output
config:
  - service_name: "api"
  - service_name: "cache"
```

### Nested Conditions

```yaml
# Input
user:
  name: "Alice"
  role: "admin"
  verified: true

permissions:
  admin_panel:
    if: "${user.role == 'admin' && user.verified}"
    access: true
  
  delete_user:
    if: "${user.role == 'admin'}"
    access: true

---

# Output
permissions:
  admin_panel:
    access: true
  
  delete_user:
    access: true
```

### Conditional List Items

```yaml
# Input
premium: false
enabled_features:
  - "basic"
  - if: "${premium}"
    name: "advanced_analytics"
  - if: "${premium}"
    name: "priority_support"

---

# Output
enabled_features:
  - "basic"
```

### Numeric Conditions

```yaml
# Input
port: 8080
timeout: 30

server_config:
  if: "${port > 0}"
  port: "${port}"
  
  slow_connection:
    if: "${timeout > 60}"
    extra_retries: 3

---

# Output
server_config:
  port: 8080
```

### Array-based Conditions

```yaml
# Input
features: ["auth", "cache", "queue"]

services:
  cache_enabled:
    if: "${contains(features, 'cache')}"
    service: "redis"
  
  queue_enabled:
    if: "${contains(features, 'queue')}"
    service: "rabbitmq"

---

# Output
services:
  cache_enabled:
    service: "redis"
  
  queue_enabled:
    service: "rabbitmq"
```

## Truthiness Rules

Conditions evaluate to `true` or `false` based on these rules:

### Truthy Values
- Non-zero integers: `1`, `42`, `-5`
- Non-zero floats: `1.5`, `0.001`
- Non-empty strings: `"text"`, `"0"` (string "0" is truthy)
- Non-empty arrays: `[1]`, `[null]`
- Non-empty objects: `{a: 1}`
- `true` boolean

### Falsy Values
- Zero: `0`, `0.0`
- Empty string: `""`
- `false` boolean
- `null` / `nil`
- Empty arrays: `[]`
- Empty objects: `{}`

**Note:** String `"0"` is **truthy** (non-empty string). Use `if: "${value == 0}"` for numeric zero.

## Error Handling

Conditions provide detailed error messages:

```yaml
# Undefined variable
if: "${undefined_var}"
# Error: undefined variable 'undefined_var' at services.database.if

# Invalid expression
if: "${item * * 2}"
# Error: error compiling expression 'item * * 2' at services.database.if

# Type mismatch
if: "${user}"
# If user is an object (not convertible to bool), proper error handling applies
```

## Behavior

1. **Evaluation**: Condition is evaluated to boolean
2. **False Result**: Entire block is omitted (not included in output)
3. **True Result**: Block is processed normally (remaining keys are evaluated)
4. **Null Output**: If condition is false, the key produces no output

## Comparison with `discard:`

- **`if: condition`**: Includes block if condition is **true**, omits if **false**
- **`discard: condition`**: Omits block if condition is **true**, includes if **false**

```yaml
# These are equivalent:
block1:
  if: false
  value: 1

block2:
  discard: true
  value: 1

# Both omit their blocks
```

## Performance Considerations

- Conditions are evaluated once per block
- Variables are resolved from the stack (cached lookups)
- Expressions are compiled and executed via expr-lang
- Complex expressions have minimal performance impact

## See Also

- [For Loops (`for:`)](/docs/handlers/for.md) - iterate with conditions
- [Interpolation (`${}`)](/docs/handlers/interpolation.md) - variable substitution in conditions
- [Discard (`discard:`)](/docs/handlers/discard.md) - inverse conditional
- [Syntax Reference](/docs/syntax.md) - complete directive syntax guide
