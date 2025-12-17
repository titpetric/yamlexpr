# Document Expansion

## Syntax Cheat Sheet

```yaml
# For loop at root level - expands to multiple documents
for: item in items
name: "${item}"
items:
  - "alice"
  - "bob"

# Matrix at root level - expands to multiple documents
matrix:
  os: [linux, windows]
  version: [1, 2]
job_name: "${os}-v${version}"

# Regular document - single output
config:
  name: "app"
  version: "1.0"
```

## Description

When `for:` or `matrix:` directives appear at the root level of a YAML document, they expand the document into multiple output documents. This is essential for CI/CD workflows and multi-document generation.

The `Parse()` and `Load()` methods return a slice of Documents, accommodating both single and multi-document results.

## Core Concepts

- **Root-level directives**: `for:` or `matrix:` at document root level expand to multiple documents
- **Document expansion**: One input document can produce multiple output documents
- **Variable propagation**: All keys in the input document become available as variables
- **Sequential output**: Multiple documents are output in order, separated by `---`

## For Loop at Root Level

A root-level `for:` creates one document per iteration:

**Input:**

```yaml
for: item in items
name: "${item}"
items:
  - "alice"
  - "bob"
  - "charlie"
```

**Output (3 documents):**

```yaml
name: "alice"
items:
  - "alice"
  - "bob"
  - "charlie"
---
name: "bob"
items:
  - "alice"
  - "bob"
  - "charlie"
---
name: "charlie"
items:
  - "alice"
  - "bob"
  - "charlie"
```

Each document includes all original variables. The `for:` directive itself is removed from each output.

## For Loop with Index at Root Level

**Input:**

```yaml
for: (idx, item) in items
index: ${idx}
name: "${item}"
items:
  - "first"
  - "second"
  - "third"
```

**Output (3 documents):**

```yaml
index: 0
name: "first"
items:
  - "first"
  - "second"
  - "third"
---
index: 1
name: "second"
items:
  - "first"
  - "second"
  - "third"
---
index: 2
name: "third"
items:
  - "first"
  - "second"
  - "third"
```

## Matrix at Root Level

A root-level `matrix:` creates one document per combination:

**Input:**

```yaml
matrix:
  os: [linux, windows]
  version: [12, 14]
job_name: "Test ${os} v${version}"
os: ${os}
version: ${version}
```

**Output (4 documents):**

```yaml
job_name: "Test linux v12"
os: linux
version: 12
---
job_name: "Test linux v14"
os: linux
version: 14
---
job_name: "Test windows v12"
os: windows
version: 12
---
job_name: "Test windows v14"
os: windows
version: 14
```

## Matrix with Exclude

**Input:**

```yaml
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
name: "${os}/${arch}"
os: ${os}
arch: ${arch}
```

**Output (3 documents):**

```yaml
name: "linux/x86_64"
os: linux
arch: x86_64
---
name: "linux/arm64"
os: linux
arch: arm64
---
name: "windows/x86_64"
os: windows
arch: x86_64
```

## Matrix with Include

**Input:**

```yaml
matrix:
  os: [linux]
  arch: [x86_64]
  include:
    - os: macos
      arch: arm64
      xcode: "14"
job: "${os}/${arch}"
os: ${os}
arch: ${arch}
xcode: ${xcode}
```

**Output (2 documents):**

```yaml
job: "linux/x86_64"
os: linux
arch: x86_64
xcode: null
---
job: "macos/arm64"
os: macos
arch: arm64
xcode: "14"
```

## Complex Example with Multiple Features

**Input:**

```yaml
for: env in environments
environment: "${env.name}"
region: "${env.region}"
services:
  - for: svc in env.services
    if: ${svc.enabled}
    name: "${svc.name}"
    port: ${svc.port}
    include: "_service-defaults.yaml"
environments:
  - name: "staging"
    region: "us-west"
    services:
      - name: "api"
        port: 8080
        enabled: true
      - name: "debug"
        port: 9000
        enabled: false
  - name: "production"
    region: "us-east"
    services:
      - name: "api"
        port: 8080
        enabled: true
      - name: "cache"
        port: 6379
        enabled: true
```

**Output (2 documents):**

Document 1 (staging):

```yaml
environment: "staging"
region: "us-west"
services:
  - name: "api"
    port: 8080
    timeout: 30
    retries: 3
environments:
  # ... original data ...
```

Document 2 (production):

```yaml
environment: "production"
region: "us-east"
services:
  - name: "api"
    port: 8080
    timeout: 30
    retries: 3
  - name: "cache"
    port: 6379
    timeout: 30
    retries: 3
environments:
  # ... original data ...
```

## Processing Multiple Documents

When using the API to process documents that may expand:

```go
expr := yamlexpr.New(os.DirFS("."))

// Load may return multiple documents
docs, err := expr.Load("config.yaml")
if err != nil {
	log.Fatal(err)
}

// Iterate over all resulting documents
for i, doc := range docs {
	fmt.Printf("Document %d:\n", i)
	output, _ := yaml.Marshal(doc)
	fmt.Println(string(output))
}
```

Or with Parse:

```go
inputDoc := map[string]any{
	"for":   "item in items",
	"name":  "${item}",
	"items": []any{"a", "b", "c"},
}

docs, err := expr.Parse(yamlexpr.Document(inputDoc))
// docs will contain 3 Document items
```

## Common Use Cases

- **CI/CD job matrices**: GitHub Actions-style test matrices
- **Multi-environment deployment**: Generate one document per environment
- **Test scenario generation**: Create multiple test configurations
- **Bulk resource creation**: Generate multiple Kubernetes manifests
- **Configuration enumeration**: Create configs for all service variants
- **Build pipelines**: Generate build jobs for multiple platforms/versions

## Comparison with Nested Expansion

### Nested For (Single Document)

**Input:**

```yaml
services:
  - for: svc in service_list
    name: "${svc}"
service_list:
  - "api"
  - "worker"
```

**Output (1 document):**

```yaml
services:
  - name: "api"
  - name: "worker"
service_list:
  - "api"
  - "worker"
```

### Root-Level For (Multiple Documents)

**Input:**

```yaml
for: svc in service_list
name: "${svc}"
service_list:
  - "api"
  - "worker"
```

**Output (2 documents):**

```yaml
name: "api"
service_list:
  - "api"
  - "worker"
---
name: "worker"
service_list:
  - "api"
  - "worker"
```

## Edge Cases

### Empty Arrays

For loops with empty arrays produce no documents:

```yaml
for: item in []
name: "${item}"
```

Result: Empty document list (0 documents)

### Single Item Expansion

Even with one item, expansion happens:

```yaml
for: item in ["single"]
name: "${item}"
```

Result: 1 document with `name: "single"`

### Combining with Nested Expansion

Root-level and nested expansions work together:

```yaml
matrix:
  env: [staging, prod]
environment: "${env}"
services:
  - for: svc in service_list
    name: "${svc}"
service_list:
  - "api"
  - "worker"
```

Result: 2 documents (from matrix), each with 2 services (from nested for)
