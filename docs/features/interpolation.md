# Interpolation

## Syntax Cheat Sheet

```yaml
${variable}                   # Simple variable substitution
${object.nested.field}        # Nested field access
${array[0]}                   # Array index access
${variable | filter}          # Apply filters (if supported)
```

## Description

Variable interpolation allows you to embed dynamic values into string fields using `${variable}` syntax. This enables configuration files to reference variables from the document root, nested objects, or expressions.

The `${}` syntax was chosen to remain valid YAML while avoiding parser ambiguity (bare `{variable}` causes YAML to expect an object structure).

## Core Concepts

- **Variables come from document root**: Any top-level key becomes a variable
- **Nested access**: Use dot notation to access nested fields (`${config.database.host}`)
- **Array access**: Use bracket notation for array indices (`${servers[0]}`)
- **Type coercion**: Non-string values are converted to their string representation

## Examples

### Simple Substitution

**Input:**

```yaml
multiplier: 3
items:
  - for: [1, 2, 3, 4]
    num: ${item}
    tripled: ${item * multiplier}
    combined: ${item + multiplier}
---
items:
  - combined: 4
    num: 1
```

**Output:**

```yaml
apiVersion: v1
app: myapp
metadata:
    name: myapp
    namespace: production
namespace: production
registry: docker.io
spec:
    image: docker.io/myapp:1.0.0
version: 1.0.0
```

### Nested Field Access

**Input:**

```yaml
config:
    database:
        host: localhost
        port: 5432
connection_string: jdbc:postgresql://localhost:5432/mydb
```

**Output:**

```yaml
config:
    database:
        host: localhost
        port: 5432
connection_string: jdbc:postgresql://localhost:5432/mydb
```

### In For Loops

Variables can be interpolated within for loop iterations. Each loop variable becomes available for substitution.

**Input:**

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
     service: ${item.name}
---
enabled_services:
  - service: "api"
  - service: "scheduler"
```

**Output:**

```yaml
enabled_services:
  - service: "api"
  - service: "scheduler"
```

## Common Use Cases

- **Configuration templates**: Reference environment names, domains, or service endpoints
- **Build configurations**: Reference artifact versions, platforms, or regions
- **Dynamic naming**: Generate names based on variables
- **Connection strings**: Build database URLs, API endpoints from components

## Supported Variable Types

All variable types are supported for interpolation:
- **Strings**: `"value"` → string
- **Numbers**: `123` → `"123"`
- **Booleans**: `true` → `"true"`
- **Objects**: Converted to YAML representation or error if used directly in string
- **Arrays**: Converted to YAML representation or error if used directly in string
