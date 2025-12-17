#!/bin/bash
# Generates conditional (if:) feature documentation
# Usage: ./conditionals.md.sh > conditionals.md

cat << 'EOF'
# Conditionals with `if:`

## Syntax Cheat Sheet

```yaml
# Omit key when condition is false
key:
  if: ${condition}
  other: value

# Include nested structure when true
config:
  if: "true"
  debug: enabled
  level: verbose

# Expression conditions
database:
  if: count > 5 && enabled
  host: localhost
  port: 5432

# Variable references
feature:
  if: ${enable_feature}
  enabled: true
```

## Description

The `if:` directive includes or omits a key and its value based on a boolean condition. When the condition is false, the entire key is removed from the output. When true, the `if:` directive itself is removed and the remaining keys are included.

This enables feature flags, environment-specific configurations, and conditional composition without leaving null/false values behind.

## Core Concepts

- **Key omission**: When `if: false`, the parent key is completely removed (not set to null)
- **Directive removal**: When `if: true`, only the `if:` key is removed; other keys remain
- **Boolean evaluation**: Conditions are evaluated as Go boolean expressions
- **Works at any level**: Can be used on top-level keys, nested maps, or array items

## Expression Types

### Boolean Literals

**Input:**
```yaml
config:
  enabled:
    if: true
    value: "yes"
  disabled:
    if: false
    value: "no"
  other: "always_present"
```

**Output:**
```yaml
config:
  enabled:
    value: "yes"
  other: "always_present"
```

### Variable References

Variables can be used directly in conditions:

**Input:**
```yaml
debug_mode: true
server:
  debug:
    if: ${debug_mode}
    level: verbose
```

**Output:**
```yaml
debug_mode: true
server:
  debug:
    level: verbose
```

### Expression Conditions

Complex conditions are evaluated using [expr-lang](https://github.com/expr-lang/expr):

**Input:**
```yaml
replicas: 3
enable_autoscaling: true
scaling_config:
  if: replicas > 1 && enable_autoscaling
  min: 2
  max: 10
```

**Output:**
```yaml
replicas: 3
enable_autoscaling: true
scaling_config:
  min: 2
  max: 10
```

### Negation

Use `!` to negate conditions:

**Input:**
```yaml
production: false
development_tools:
  if: !production
  debug: true
  hotreload: enabled
```

**Output:**
```yaml
production: false
development_tools:
  debug: true
  hotreload: enabled
```

## With For Loops

The `if:` directive works alongside `for:` to filter items:

**Input:**
```yaml
services:
  - for: service in all_services
    if: ${service.enabled}
    name: "${service.name}"
    port: ${service.port}
all_services:
  - name: api
    enabled: true
    port: 8080
  - name: worker
    enabled: false
    port: 9000
  - name: cache
    enabled: true
    port: 6379
```

**Output:**
```yaml
services:
  - name: api
    port: 8080
  - name: cache
    port: 6379
all_services:
  - name: api
    enabled: true
    port: 8080
  - name: worker
    enabled: false
    port: 9000
  - name: cache
    enabled: true
    port: 6379
```

## Array Item Filtering

When `if:` appears on an array item, the entire item is omitted if false:

**Input:**
```yaml
items:
  - name: "item1"
    active: true
  - name: "item2"
    if: false
    active: false
  - name: "item3"
    active: true
```

**Output:**
```yaml
items:
  - name: "item1"
    active: true
  - name: "item3"
    active: true
```

## Nested Conditions

Multiple `if:` directives can be nested:

**Input:**
```yaml
app:
  cache:
    if: ${use_cache}
    redis:
      if: ${use_redis}
      host: redis.local
      port: 6379
    memcached:
      if: ${use_memcached}
      host: memcached.local
use_cache: true
use_redis: true
use_memcached: false
```

**Output:**
```yaml
app:
  cache:
    redis:
      host: redis.local
      port: 6379
use_cache: true
use_redis: true
use_memcached: false
```

## Common Use Cases

- **Feature flags**: Conditionally include features based on configuration
- **Environment-specific settings**: Different config for dev/staging/prod
- **Dependency-based configuration**: Include services only if prerequisites are enabled
- **Optional components**: Cloud providers, logging backends, monitoring systems
- **Scaling policies**: Include autoscaling config only for production

## Supported Condition Types

| Type | Example | Notes |
|------|---------|-------|
| Boolean | `if: true` | Literal true/false |
| Variable | `if: ${flag}` | References top-level variable |
| Field access | `if: item.enabled` | Dot notation for nested fields |
| Comparison | `if: count > 5` | `<`, `>`, `<=`, `>=`, `==`, `!=` |
| Logical AND | `if: a && b` | Both conditions must be true |
| Logical OR | `if: a \|\| b` | Either condition can be true |
| Negation | `if: !flag` | Inverts the condition |
| Complex | `if: (a > 5) && (b == "test")` | Parentheses for grouping |

EOF
