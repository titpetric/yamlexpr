## Simple Value Iteration

Iterate over a list of values with a single loop variable.

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

## With Index and Value

Unpack both the index position and value in a loop.

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

## Nested For Loops

Iterate over nested structures with multiple levels of loops.

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

## For Loop with Filter Condition

Combine for: with if: to filter items during iteration.

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

## For Loop with Expressions

Use expressions and calculations within for loop templates.

**Input:**

```yaml
multiplier: 3
items:
  - for: item in [1, 2, 3, 4]
    num: ${item}
    tripled: ${item * multiplier}
    combined: ${item + multiplier}
```

**Output:**

```yaml
multiplier: 3
items:
  - num: 1
    tripled: 3
    combined: 4
  - num: 2
    tripled: 6
    combined: 5
  - num: 3
    tripled: 9
    combined: 6
  - num: 4
    tripled: 12
    combined: 7
```

## For Loop with Empty Array

For loops over empty arrays produce no output items.

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
