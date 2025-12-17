#!/bin/bash
# Generates for loop feature documentation
# Usage: ./for-loops.md.sh > for-loops.md

cat << 'EOF'
# For Loops with `for:`

## Syntax Cheat Sheet

```yaml
# Simple iteration over array variable
items:
  - for: item in items_list
    name: "${item}"

# With index and value
services:
  - for: (index, service) in all_services
    number: ${index}
    name: "${service.name}"

# Ignore index with underscore
configs:
  - for: (_, config) in configurations
    setting: "${config.value}"

# Ignore value with underscore
indexes:
  - for: (idx, _) in items
    position: ${idx}

# Direct array literal
statuses:
  - for: status in ["active", "pending", "failed"]
    current: "${status}"
```

## Description

The `for:` directive expands an array by iterating over values and creating multiple items. This is essential for generating repetitive configuration structures and templates.

Each iteration creates a copy of the template with loop variables available for interpolation and expressions. The `for:` directive itself is removed from each output item.

## Core Concepts

- **Iteration variable**: The variable name(s) introduced by the loop
- **Source array**: The array being iterated (can be a variable reference or literal array)
- **Template**: All other keys in the block become the template for each iteration
- **Scope**: Loop variables are available in the iteration scope via interpolation and expressions

## Iteration Patterns

### Simple Value Iteration

**Input:**
```yaml
servers:
  - for: server in server_list
    name: "${server}"
server_list:
  - "api"
  - "worker"
  - "cache"
```

**Output:**
```yaml
servers:
  - name: "api"
  - name: "worker"
  - name: "cache"
server_list:
  - "api"
  - "worker"
  - "cache"
```

### With Index and Value

**Input:**
```yaml
indexed_items:
  - for: (idx, item) in items
    index: ${idx}
    value: "${item}"
items:
  - "first"
  - "second"
  - "third"
```

**Output:**
```yaml
indexed_items:
  - index: 0
    value: "first"
  - index: 1
    value: "second"
  - index: 2
    value: "third"
items:
  - "first"
  - "second"
  - "third"
```

### Iterating Over Objects

**Input:**
```yaml
services:
  - for: svc in service_configs
    name: "${svc.name}"
    port: ${svc.port}
    replicas: ${svc.replicas}
service_configs:
  - name: "api"
    port: 8080
    replicas: 3
  - name: "worker"
    port: 9000
    replicas: 1
```

**Output:**
```yaml
services:
  - name: "api"
    port: 8080
    replicas: 3
  - name: "worker"
    port: 9000
    replicas: 1
service_configs:
  - name: "api"
    port: 8080
    replicas: 3
  - name: "worker"
    port: 9000
    replicas: 1
```

### With Filter Conditions

Combine `for:` with `if:` to filter items:

**Input:**
```yaml
enabled_services:
  - for: svc in all_services
    if: ${svc.enabled}
    name: "${svc.name}"
    port: ${svc.port}
all_services:
  - name: "api"
    enabled: true
    port: 8080
  - name: "disabled-worker"
    enabled: false
    port: 9000
  - name: "cache"
    enabled: true
    port: 6379
```

**Output:**
```yaml
enabled_services:
  - name: "api"
    port: 8080
  - name: "cache"
    port: 6379
all_services:
  - name: "api"
    enabled: true
    port: 8080
  - name: "disabled-worker"
    enabled: false
    port: 9000
  - name: "cache"
    enabled: true
    port: 6379
```

### Nested For Loops

For loops can be nested for multi-level expansion:

**Input:**
```yaml
matrix:
  - for: os in operating_systems
    os: "${os}"
    versions:
      - for: version in versions_list
        version: "${version}"
operating_systems:
  - "ubuntu"
  - "windows"
versions_list:
  - "18.04"
  - "20.04"
```

**Output:**
```yaml
matrix:
  - os: "ubuntu"
    versions:
      - version: "18.04"
      - version: "20.04"
  - os: "windows"
    versions:
      - version: "18.04"
      - version: "20.04"
operating_systems:
  - "ubuntu"
  - "windows"
versions_list:
  - "18.04"
  - "20.04"
```

### Ignoring Index or Value with Underscore

Use `_` to ignore the index or value:

**Input:**
```yaml
# Ignore index, only use value
names:
  - for: (_, name) in ["alice", "bob", "charlie"]
    person: "${name}"

# Ignore value, only use index
positions:
  - for: (i, _) in ["a", "b", "c"]
    index: ${i}
```

**Output:**
```yaml
names:
  - person: "alice"
  - person: "bob"
  - person: "charlie"
positions:
  - index: 0
  - index: 1
  - index: 2
```

### Empty Array Handling

Iterating over an empty array produces no output items:

**Input:**
```yaml
items:
  - for: item in []
    value: "${item}"
```

**Output:**
```yaml
items: []
```

## With Interpolation and Expressions

Loop variables can be used in interpolations and expressions:

**Input:**
```yaml
build_config:
  - for: (idx, platform) in platforms
    if: idx > 0 || !skip_first
    platform: "${platform}"
    artifact: "build-${platform}-v1.0.0.tar.gz"
    size_mb: "${platform == 'arm64' ? 150 : 250}"
platforms:
  - "amd64"
  - "arm64"
  - "arm"
skip_first: false
```

**Output:**
```yaml
build_config:
  - platform: "amd64"
    artifact: "build-amd64-v1.0.0.tar.gz"
    size_mb: "250"
  - platform: "arm64"
    artifact: "build-arm64-v1.0.0.tar.gz"
    size_mb: "150"
  - platform: "arm"
    artifact: "build-arm-v1.0.0.tar.gz"
    size_mb: "250"
platforms:
  - "amd64"
  - "arm64"
  - "arm"
skip_first: false
```

## Nested Structures with Complex Templates

**Input:**
```yaml
environments:
  - for: env in environment_list
    name: "${env.name}"
    replicas: ${env.replicas}
    database:
      host: "db-${env.region}.internal"
      port: 5432
      credentials:
        include: "_db-credentials.yaml"
environment_list:
  - name: "staging"
    region: "us-west"
    replicas: 2
  - name: "production"
    region: "us-east"
    replicas: 5
```

**Output:**
```yaml
environments:
  - name: "staging"
    replicas: 2
    database:
      host: "db-us-west.internal"
      port: 5432
      credentials:
        username: "dbuser"
        password: "secret"
  - name: "production"
    replicas: 5
    database:
      host: "db-us-east.internal"
      port: 5432
      credentials:
        username: "dbuser"
        password: "secret"
environment_list:
  - name: "staging"
    region: "us-west"
    replicas: 2
  - name: "production"
    region: "us-east"
    replicas: 5
```

## Common Use Cases

- **Service enumeration**: Generate config for each service/microservice
- **Platform builds**: Create build configurations for multiple targets
- **Environment scaling**: Define replicas or resources for different environments
- **Batch operations**: Generate similar structures with variations
- **Test matrices**: Create combinations of parameters

## Edge Cases and Special Behavior

### Quoted vs Unquoted Syntax

Both quoted and unquoted syntax are supported:

```yaml
# Quoted (recommended)
- for: "item in items_list"

# Unquoted (also valid)
- for: item in items_list
```

### Variable Source

The source array can be:
- A variable reference: `for: item in items`
- A literal array: `for: item in ["a", "b", "c"]`
- From nested paths: `for: item in config.services` (if available in scope)

### Order of Evaluation

Loop variables are available:
1. In field values via interpolation: `"${item}"`
2. In conditional expressions: `if: item.enabled`
3. In nested structures created during that iteration

EOF
