# Interpolation Handler

The **interpolation** handler processes `${...}` placeholder syntax in string values, replacing them with resolved values from the variable stack.

## Overview

Interpolation is automatically applied to all string values during document processing. This handler provides utilities and functions to support both simple variable references and complex expressions.

## Features

- **Simple variable references**: `${variable}` → resolved from stack
- **Nested paths**: `${user.name}` → navigates nested objects
- **Expression evaluation**: `${item * 2}` → full expr-lang expressions
- **Type preservation**: `${count}` returns native type (number, boolean, etc.), not string
- **Multiple modes**: Strict (error on undefined) and permissive (returns nil)

## Syntax

### Variable References

```yaml
message: "Hello ${name}"
greeting: "Welcome back ${user.first_name}"
port: ${server.port}
```

### Expression Evaluation

The `${}` syntax supports full [expr-lang](https://github.com/expr-lang/expr) expressions:

**Arithmetic:**

```yaml
doubled: ${count * 2}
discounted: ${price * 0.9}
remaining: ${total - used}
```

**String Functions:**

```yaml
uppercase: ${upper(name)}
length: ${len(items)}
contains: ${contains(text, "search")}
```

**Logical Comparisons:**

```yaml
is_active: ${status == "active"}
is_large: ${count > 10}
is_even: ${count % 2 == 0}
```

**Array/Object Operations:**

```yaml
item_count: ${len(items)}
has_items: ${len(items) > 0}
first_item: ${items[0]}
```

## Single vs. Multiple Interpolations

### Single Interpolation (Type Preserved)

When a value contains **only one** `${}` expression, the native type is preserved:

```yaml
count: ${items.length}      # Returns: 5 (number, not "5")
enabled: ${is_active}       # Returns: true (boolean, not "true")
config: ${database}         # Returns: map, list, or any type
nothing: ${null_variable}   # Returns: null (not "null")
```

### Multiple Interpolations (String Result)

When a value contains **multiple** expressions or mixed text, result is always a string:

```yaml
message: "Count: ${count}, Active: ${is_active}"  # Returns: "Count: 5, Active: true"
path: "/home/${user}/files/${dir}"                # Returns: string
```

## Modes

### Strict Mode (Default)

Errors on undefined variables. Used for document processing:

```go
result, err := InterpolateValue(value, stack, path)
// If variable is undefined: error
```

**Usage:**
- Document processing (catches configuration errors early)
- Conditional expressions (strict type checking)

### Permissive Mode

Returns `nil` if any variable is undefined. Used for optional values:

```go
result, err := InterpolateValuePermissive(value, stack)
// If variable is undefined: returns nil
```

**Usage:**
- Matrix dimensions (optional configuration)
- Fallback values
- Optional interpolation

## API Functions

### String Interpolation

**`InterpolateStringWithContext(s string, stack *Stack, path string) (string, error)`**

Replaces `${...}` placeholders with values from the stack. Always returns a string.

```go
result, err := InterpolateStringWithContext("Hello ${name}", stack, "greeting")
// Returns: "Hello Alice", nil
```

**Parameters:**
- `s`: String containing `${...}` patterns
- `stack`: Variable stack for resolution
- `path`: Dot-separated path for error context (e.g., "config.database.host")

**Errors:**
- Undefined variable: `"undefined variable 'name' at config.database.host"`
- Type conversion failure: `"cannot convert variable 'config' to string at message"`
- Expression evaluation: `"error evaluating expression 'item * 2' at results[0]"`

### Value Interpolation (Type Preserving)

**`InterpolateValue(value any, stack *Stack, path string) (any, error)`**

Interpolates a value while preserving native types for single expressions.

```go
result, err := InterpolateValue("${count}", stack, "item.quantity")
// Returns: 5 (number), nil

result, err := InterpolateValue("${name}", stack, "greeting")
// Returns: "Alice" (string), nil

result, err := InterpolateValue(false, stack, "setting")
// Returns: false (unchanged for non-strings), nil
```

### Value Interpolation with Context

**`InterpolateValueWithContext(s string, stack *Stack, path string) (any, error)`**

Similar to `InterpolateValue`, works specifically with strings.

```go
result, err := InterpolateValueWithContext("${items}", stack, "data")
// Returns: []any{...} (preserved type), nil
```

### Permissive Interpolation

**`InterpolateStringPermissive(s string, stack *Stack) (any, error)`**

Returns `nil` if any variable is undefined (vs. error in strict mode).

```go
result, err := InterpolateStringPermissive("${optional_var}", stack)
// If undefined: returns nil, nil
// If defined: returns "value", nil
```

**`InterpolateValuePermissive(value any, stack *Stack) (any, error)`**

Permissive version for any value type.

```go
result, err := InterpolateValuePermissive(value, stack)
// Returns nil if any variable is undefined
```

### Utility Functions

**`ContainsInterpolation(s string) bool`**

Checks if a string contains `${...}` patterns.

```go
if ContainsInterpolation(value) {
	// Value needs interpolation
}
```

## Examples

### Basic Variable Interpolation

```yaml
# Input
name: "Alice"
greeting: "Hello ${name}"

---

# Output
name: "Alice"
greeting: "Hello Alice"
```

### Expression Evaluation

```yaml
# Input
count: 5
items:
  - for: "item in values"
    name: "item_${item}"
    doubled: ${item * 2}

---

# Output
count: 5
items:
  - name: "item_1"
    doubled: 2
  - name: "item_2"
    doubled: 4
  - name: "item_3"
    doubled: 6
  - name: "item_4"
    doubled: 8
  - name: "item_5"
    doubled: 10
```

### Type Preservation

```yaml
# Input
active: true
items:
  - 1
  - 2
database:
  host: "localhost"

config:
  is_active: ${active}         # Boolean preserved
  item_list: ${items}          # Array preserved
  db_config: ${database}       # Object preserved
  message: "DB: ${database}"   # String (mixed content)

---

# Output
config:
  is_active: true              # Still boolean
  item_list:                   # Still array
    - 1
    - 2
  db_config:                   # Still object
    host: "localhost"
  message: "DB: map[host:localhost]"  # String representation
```

### Conditional Interpolation

```yaml
# Input
user:
  name: "Alice"
  role: "admin"

permissions:
  - action: "delete"
    if: "${user.role == 'admin'}"

---

# Output
permissions:
  - action: "delete"
```

### Using Expressions for Calculations

```yaml
# Input
price: 100
quantity: 5
tax_rate: 0.08

order:
  subtotal: ${price * quantity}
  tax: ${price * quantity * tax_rate}
  total: ${price * quantity * (1 + tax_rate)}

---

# Output
order:
  subtotal: 500
  tax: 40
  total: 540
```

## Error Handling

Interpolation provides detailed error messages with context:

```yaml
# Input (with undefined variable)
message: "Hello ${undefined_name}"

# Error output:
# undefined variable 'undefined_name' at greeting
```

```yaml
# Input (with invalid expression)
count: ${invalid * * 2}

# Error output:
# error compiling expression 'invalid * * 2' at counter
```

## Implementation Details

- **Pattern**: `\$\{([^}]+)\}` - matches `${...}` with any content inside
- **Engine**: Uses [expr-lang](https://github.com/expr-lang/expr) for expression evaluation
- **Scope**: Accesses all stack variables (current scope + parent scopes)
- **Order**: Processes all `${}` patterns left-to-right in a string

## Configuration

Interpolation behavior can be controlled via the main Expr configuration:

```go
expr := yamlexpr.New(fs)

// Default: strict mode (errors on undefined)
result, err := expr.Process(data)

// Custom stack: provides additional variables
stack := stack.New(map[string]any{
	"env":     "production",
	"version": "1.0.0",
})
result, err := expr.ProcessWithStack(data, stack)
```

## See Also

- [Conditionals (`if:`)](/docs/handlers/if.md) - use interpolation in conditions
- [For Loops (`for:`)](/docs/handlers/for.md) - use interpolation in loop templates
- [Syntax Reference](/docs/syntax.md) - complete directive syntax guide
- [API Reference](/docs/api.md) - complete API documentation
