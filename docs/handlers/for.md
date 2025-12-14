# For Loop Handler (`for:`)

The **`for:` directive** expands templates by iterating over collections, with loop variables available for interpolation.

## Overview

The `for:` directive repeats a template for each item in a collection (array), automatically binding loop variables that can be referenced in interpolations and conditions.

## Syntax

The `for:` directive accepts a **single-line expression** in string format.

### Basic Iteration

```yaml
items:
  - for: "item in items"
    name: "${item.name}"
    price: "${item.price}"
```

### Index and Item

```yaml
items:
  - for: "(idx, item) in items"
    index: "${idx}"
    name: "${item.name}"
```

### Direct Array Literal

```yaml
items:
  - for: [1, 2, 3]
    number: "${item}"
    squared: "${item * item}"
```

### Nested Path

```yaml
users:
  - for: "account in user.accounts"
    account_id: "${account.id}"
```

## Parsing

### Expression Format

The `for:` directive uses the format: `<variables> in <source>`

**Variables:**
- Single variable: `item` → binds each item
- Index and item: `(idx, item)` → binds index (0-based) and item
- Omit variable: `(_, item)` or `(idx, _)` → skip binding specific variables
- Whitespace handling: spaces around commas are trimmed

**Source:**
- Variable name: `items` → resolves from stack
- Nested path: `user.accounts` → navigates object path
- Must resolve to an array at runtime

### Parsing Error Cases

```go
// Valid
"item in items"
"(idx, item) in items"
"(_, item) in items"

// Invalid - caught at parse time
"item in"               // No source
"in items"              // No variables
"(idx,) in items"       // Trailing comma
"(idx, item,) in items" // Trailing comma in parens
"(0idx, item) in items" // Variable starts with digit
"() in items"           // Empty variable list
```

## Loop Variables

### Simple Loop

```yaml
# Input
items: ["apple", "banana", "cherry"]

shopping:
  - for: "fruit in items"
    name: "${fruit}"

# Output
shopping:
  - name: "apple"
  - name: "banana"
  - name: "cherry"
```

### Index and Item

```yaml
# Input
colors: ["red", "green", "blue"]

palette:
  - for: "(idx, color) in colors"
    position: "${idx}"
    name: "${color}"

# Output
palette:
  - position: 0
    name: "red"
  - position: 1
    name: "green"
  - position: 2
    name: "blue"
```

### Omitting Variables

```yaml
# Input
items: ["a", "b", "c"]

# Only use item, ignore index
results1:
  - for: "(_, item) in items"
    value: "${item}"

# Only use index, ignore item
results2:
  - for: "(idx, _) in items"
    number: "${idx}"

# Output
results1:
  - value: "a"
  - value: "b"
  - value: "c"

results2:
  - number: 0
  - number: 1
  - number: 2
```

## Root-Level Loops

When `for:` is used at the document root level, it produces multiple documents:

```yaml
# config.yaml
for: "env in environments"
name: "${env.name}"
settings: "${env.config}"

# Input
environments:
  - name: "dev"
    config: {debug: true}
  - name: "prod"
    config: {debug: false}

# Output (two documents)
---
name: "dev"
settings:
  debug: true
---
name: "prod"
settings:
  debug: false
```

## Combining with Other Directives

### For with If

```yaml
# Input
services:
  - name: "api"
    enabled: true
  - name: "worker"
    enabled: false

config:
  - for: "service in services"
    if: "${service.enabled}"
    service_name: "${service.name}"

# Output
config:
  - service_name: "api"
```

### Nested Loops

```yaml
# Input
projects:
  - name: "project1"
    tasks: ["task1", "task2"]
  - name: "project2"
    tasks: ["task3"]

results:
  - for: "project in projects"
    project_name: "${project.name}"
    items:
      - for: "task in project.tasks"
        task_name: "${task}"

# Output
results:
  - project_name: "project1"
    items:
      - task_name: "task1"
      - task_name: "task2"
  - project_name: "project2"
    items:
      - task_name: "task3"
```

### For with Expressions

```yaml
# Input
numbers: [1, 2, 3, 4, 5]

results:
  - for: "(idx, num) in numbers"
    original: "${num}"
    doubled: "${num * 2}"
    is_even: "${num % 2 == 0}"
    position: "${idx + 1}"

# Output
results:
  - original: 1
    doubled: 2
    is_even: false
    position: 1
  - original: 2
    doubled: 4
    is_even: true
    position: 2
  - original: 3
    doubled: 6
    is_even: false
    position: 3
  - original: 4
    doubled: 8
    is_even: true
    position: 4
  - original: 5
    doubled: 10
    is_even: false
    position: 5
```

## API Functions

### Parse For Expression

**`ParseForExpr(expr string) (*ForLoopExpr, error)`**

Parses a `for:` directive expression string.

```go
loopExpr, err := ParseForExpr("item in items")
// Returns: ForLoopExpr{Variables: ["item"], Source: "items"}, nil

loopExpr, err := ParseForExpr("(idx, item) in items")
// Returns: ForLoopExpr{Variables: ["idx", "item"], Source: "items"}, nil
```

**Returns:**
- `ForLoopExpr.Variables`: List of variable names
- `ForLoopExpr.Source`: Source expression to iterate over
- `error`: Parsing error with context (missing " in ", invalid syntax, etc.)

### Build Loop Scope

**`BuildScope(varNames []string, idx int, item any) map[string]any`**

Creates a variable scope for a loop iteration.

```go
scope := BuildScope([]string{"idx", "item"}, 0, "apple")
// Returns: map[string]any{"idx": 0, "item": "apple"}
```

**Parameters:**
- `varNames`: Variable names from parsed for expression
- `idx`: Zero-based iteration index
- `item`: Current item from collection

**Logic:**
- First variable with one var: bound to item
- First variable with two vars: bound to index
- Second variable (if exists): bound to item
- Underscore variables: omitted from scope

### Create For Handler

**`ForHandler(proc Processor, forDirective string) DirectiveHandler`**

Creates a handler for the `for:` directive.

```go
handler := ForHandler(processor, "for")
// Use with: expr.WithDirectiveHandler("for", handler, 100)
```

### Expand For at Root

**`ExpandForAtRoot(ctx *Context, template map[string]any, processor Processor, forDirective string) ([]any, error)`**

Expands a root-level `for:` directive into multiple documents.

```go
docs, err := ExpandForAtRoot(ctx, template, processor, "for")
// Returns: []any with one document per loop iteration
```

## Scope Management

### Variable Scope

Loop variables are added to the stack temporarily:

```go
// Push: before processing iteration
ctx.PushStackScope(scope)

// Process template with loop variables in scope
result := processor.ProcessMapWithContext(ctx, template)

// Pop: after processing iteration
ctx.PopStackScope()
```

### Scope Rules

1. Loop variables shadow parent scope variables with same names
2. Parent scope variables are still accessible
3. Scope is popped after each iteration (variables don't leak)

```yaml
# Input
item: "parent_item"
items: ["child1", "child2"]

config:
  - for: "item in items"
    current: "${item}"         # "child1", "child2" (loop var)

after_loop: "${item}"          # "parent_item" (parent var)

# Output
config:
  - current: "child1"
  - current: "child2"

after_loop: "parent_item"
```

## Error Handling

Detailed error messages with context:

```yaml
# Undefined source variable
for: "item in undefined_items"
# Error: undefined variable 'undefined_items' at items.for

# Invalid expression format
for: "item items"  # Missing " in "
# Error: invalid for expression 'item items': no ' in ' found in for expression

# Source not array
for: "item in config"  # If config is object/string
items: [1, 2, 3]
# Error: for: variable 'config' must be an array, got map[string]any

# Invalid variable name
for: "(0idx, item) in items"
# Error: invalid variable name '0idx': must not start with a digit

# Empty variable
for: "(idx, , item) in items"
# Error: empty variable name in list
```

## Performance Considerations

- **Template Copying**: Template is copied for each iteration (shallow copy)
- **Stack Management**: Stack scope is pushed/popped for each iteration (efficient)
- **Interpolation**: Variables are resolved on-demand during interpolation
- **Nested Loops**: Each nesting level adds minimal overhead

## Comparison with Root-Level For

### Document-Level For Loop

```yaml
# config.yaml
services:
  - for: "service in configs"
    name: "${service.name}"
```

**Produces:** Single document with expanded array

### Root-Level For Loop

```yaml
# config.yaml
for: "env in environments"
name: "${env}"
config: "${env}"
```

**Produces:** Multiple documents (one per iteration)

## Examples

### Configuration Matrix

```yaml
# Input
environments: ["dev", "staging", "prod"]
services:
  - name: "api"
  - name: "worker"

configs:
  - for: "(env_idx, env) in environments"
    environment: "${env}"
    services:
      - for: "service in services"
        service_name: "${service.name}"
        env_level: "${env_idx}"

# Output
configs:
  - environment: "dev"
    services:
      - service_name: "api"
        env_level: 0
      - service_name: "worker"
        env_level: 0
  - environment: "staging"
    services:
      - service_name: "api"
        env_level: 1
      - service_name: "worker"
        env_level: 1
  - environment: "prod"
    services:
      - service_name: "api"
        env_level: 2
      - service_name: "worker"
        env_level: 2
```

### Filtered List

```yaml
# Input
users:
  - name: "Alice"
    role: "admin"
  - name: "Bob"
    role: "user"
  - name: "Charlie"
    role: "admin"

admins:
  - for: "user in users"
    if: "${user.role == 'admin'}"
    name: "${user.name}"

# Output
admins:
  - name: "Alice"
  - name: "Charlie"
```

### Index-Based Operations

```yaml
# Input
items: ["a", "b", "c"]

list:
  - for: "(idx, item) in items"
    item: "${item}"
    position_label: "${item}${idx + 1}"  # a1, b2, c3

# Output
list:
  - item: "a"
    position_label: "a1"
  - item: "b"
    position_label: "b2"
  - item: "c"
    position_label: "c3"
```

## See Also

- [Conditionals (`if:`)](/docs/handlers/if.md) - filter loop items
- [Interpolation (`${}`)](/docs/handlers/interpolation.md) - use variables in loop
- [Matrix (`matrix:`)](/docs/handlers/matrix.md) - cartesian product expansion
- [Syntax Reference](/docs/syntax.md) - complete directive syntax guide
