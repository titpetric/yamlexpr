#!/bin/bash
# Generates interpolation feature documentation
# Usage: ./interpolation.md.sh > interpolation.md

cat << 'EOF'
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
EOF

cat /root/github/yamlexpr-main/testdata/fixtures/046-expression-with-variables.yaml | head -10
cat << 'EOF'
```

**Output:**

```yaml
EOF

cd /root/github/yamlexpr-main && go run cmd/yamlexpr/main.go << 'YAMLEOF' 2>/dev/null
apiVersion: v1
metadata:
  name: ${app}
  namespace: ${namespace}
spec:
  image: "${registry}/${app}:${version}"
app: "myapp"
namespace: "production"
registry: "docker.io"
version: "1.0.0"
YAMLEOF

cat << 'EOF'
```

### Nested Field Access

**Input:**

```yaml
EOF

cat > /tmp/nested.yaml << 'YAMLEOF'
config:
  database:
    host: "localhost"
    port: 5432
connection_string: "jdbc:postgresql://${config.database.host}:${config.database.port}/mydb"
YAMLEOF

cd /root/github/yamlexpr-main && go run cmd/yamlexpr/main.go /tmp/nested.yaml 2>/dev/null || cat /tmp/nested.yaml

cat << 'EOF'
```

**Output:**

```yaml
EOF

cd /root/github/yamlexpr-main && go run cmd/yamlexpr/main.go << 'YAMLEOF' 2>/dev/null
config:
  database:
    host: "localhost"
    port: 5432
connection_string: "jdbc:postgresql://${config.database.host}:${config.database.port}/mydb"
YAMLEOF

cat << 'EOF'
```

### In For Loops

Variables can be interpolated within for loop iterations. Each loop variable becomes available for substitution.

**Input:**

```yaml
EOF

cat /root/github/yamlexpr-main/testdata/fixtures/050-for-with-if-filtering.yaml | head -20

cat << 'EOF'
```

**Output:**

```yaml
EOF

cd /root/github/yamlexpr-main && cat testdata/fixtures/050-for-with-if-filtering.yaml | sed -n '/^---$/,$ {/^---$/d; p}'

cat << 'EOF'
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

EOF
