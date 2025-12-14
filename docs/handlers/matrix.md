# Matrix Handler (`matrix:`)

The **`matrix:` directive** expands templates using cartesian product combinations of dimensions, similar to GitHub Actions matrix strategy.

## Overview

The `matrix:` directive generates multiple configurations from a set of dimensions and their values. It automatically combines all dimension values, applies exclusion and inclusion rules, and makes matrix variables available for interpolation.

## Syntax

```yaml
jobs:
  test:
    matrix:
      os: [linux, windows, macos]
      arch: [x86_64, arm64]
      exclude:
        - os: windows
          arch: arm64
      include:
        - os: darwin
          arch: arm64
          silicon: true
```

## Dimensions

Matrix dimensions are key-value pairs where values are arrays:

```yaml
matrix:
  os: [linux, windows, macos]
  arch: [x86_64, arm64]
  node_version: [18, 20, 22]
```

This creates a **cartesian product**: 3 × 2 × 3 = **18 combinations**.

## Cartesian Product

All dimension values are combined:

```yaml
# Input
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]

# Generates
os: linux,   arch: x86_64
os: linux,   arch: arm64
os: windows, arch: x86_64
os: windows, arch: arm64
```

## Exclude Rules

Exclude rules remove specific combinations:

```yaml
matrix:
  os: [linux, windows, macos]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64        # Removes: windows + arm64
    - os: macos
      arch: x86_64       # Removes: macos + x86_64

# Remaining combinations: 6 - 2 = 4 combinations
```

**Matching Logic**: A job matches an exclude rule if **ALL** keys in the rule match the job.

```yaml
# This exclude rule:
- os: windows
  arch: arm64

# Matches jobs with BOTH:
# - os == "windows"
# - arch == "arm64"

# Does NOT match:
# - os: windows, arch: x86_64  (arch doesn't match)
# - os: linux, arch: arm64     (os doesn't match)
```

## Include Rules

Include rules add new combinations or override existing ones:

```yaml
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  include:
    - os: darwin
      arch: arm64
      silicon: true
```

**Matching Logic:**
- If all keys in include rule match an existing job, **merge** into that job
- If no match found, **add** as new job

```yaml
# Input
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  include:
    - os: windows
      arch: x86_64
      mingw: true         # Matches windows + x86_64, adds mingw
    
    - os: macos
      arch: arm64         # No match, creates new job

# Result: 5 jobs (4 base - exclude nothing + 1 new + 1 merge)
```

## Single-Line Expression

Matrix dimensions must be specified as arrays in YAML. The `matrix:` key itself is a **single-line structure**.

```yaml
# Valid - all on single logical line
matrix:
  dimension1: [value1, value2]
  dimension2: [value3, value4]
```

## Variables in Scope

Each matrix combination makes its dimension values available as variables:

```yaml
# Input
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]

config:
  name: "Build for ${os}/${arch}"

# Output (4 documents)
---
config:
  name: "Build for linux/x86_64"
---
config:
  name: "Build for linux/arm64"
---
config:
  name: "Build for windows/x86_64"
---
config:
  name: "Build for windows/arm64"
```

## Examples

### Basic Matrix

```yaml
# Input
matrix:
  os: [ubuntu, macos, windows]
  python: [3.9, 3.10, 3.11]

jobs:
  test:
    os: "${os}"
    python: "${python}"
    name: "Test on ${os} / Python ${python}"

# Output (9 jobs)
jobs:
  test:
    os: "ubuntu"
    python: 3.9
    name: "Test on ubuntu / Python 3.9"
  
  test:
    os: "ubuntu"
    python: 3.10
    name: "Test on ubuntu / Python 3.10"
  
  # ... 7 more combinations
```

### With Exclusions

```yaml
# Input
matrix:
  os: [ubuntu, windows, macos]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
    - os: macos
      arch: x86_64

jobs:
  build:
    os: "${os}"
    arch: "${arch}"

# Output (4 jobs)
# ubuntu + x86_64
# ubuntu + arm64
# windows + x86_64
# macos + arm64
```

### With Inclusions

```yaml
# Input
matrix:
  os: [ubuntu, windows]
  include:
    - os: macos
      arch: arm64
      xcode: "14.2"

jobs:
  test:
    os: "${os}"
    arch: "${arch}"
    xcode: "${xcode}"

# Output (3+ jobs, depending on base dimensions)
jobs:
  test:
    os: "ubuntu"
    arch: null      # Not in dimension
    xcode: null
  
  test:
    os: "windows"
    arch: null
    xcode: null
  
  test:
    os: "macos"
    arch: "arm64"   # From include rule
    xcode: "14.2"   # From include rule
```

### Complex Example with Conditions

```yaml
# Input
matrix:
  os: [ubuntu, windows, macos]
  node: [18, 20, 22]
  exclude:
    - os: windows
      node: 18
  include:
    - os: docker
      node: 20
      container: "node:20"

build:
  os: "${os}"
  node: "${node}"
  
  install_deps:
    if: "${os != 'docker'}"
    command: "npm install"
  
  use_container:
    if: "${container}"
    image: "${container}"
  
  name: "Build on ${os} / Node ${node}"
```

## API Functions

### Create Matrix Handler

**`NewMatrixHandler() DirectiveHandler`**

Creates a handler for the `matrix:` directive.

```go
handler := NewMatrixHandler()
// Use with: expr.WithDirectiveHandler("matrix", handler, 5)
```

### Parse Matrix Directive

**`parseMatrixDirective(m map[string]any) (*MatrixDirective, error)`**

Internal function that parses the matrix map into dimensions, includes, and excludes.

```go
// Returns: MatrixDirective with parsed structure
```

### Expand Matrix Base

**`expandMatrixBase(md *MatrixDirective) []map[string]any`**

Generates cartesian product of all dimensions.

```go
// Input: Dimensions {os: [linux, windows], arch: [x86_64, arm64]}
// Output: 4 job maps with all combinations
```

### Apply Exclusions

**`applyExcludes(jobs []map[string]any, excludes []map[string]any) []map[string]any`**

Removes jobs matching exclude rules.

```go
// Filters out jobs where ALL exclude keys match
```

### Apply Inclusions

**`applyEmbeds(jobs []map[string]any, embeds []map[string]any) ([]map[string]any, error)`**

Merges or adds jobs from include rules.

```go
// Merges into existing matching jobs
// Creates new jobs if no match found
```

### Expand Matrix at Root

**`ExpandMatrixAtRoot(ctx *Context, template map[string]any, processor Processor) ([]any, error)`**

Expands a root-level `matrix:` directive into multiple documents.

```go
docs, err := ExpandMatrixAtRoot(ctx, template, processor)
// Returns: Multiple documents (one per matrix combination)
```

## Root-Level Matrix

When `matrix:` is at document root, it produces multiple documents:

```yaml
# config.yaml
matrix:
  env: [dev, prod]
  region: [us-west, eu-west]

environment: "${env}"
region: "${region}"
config: "/config/${env}/${region}.yaml"

# Output (4 documents)
---
environment: "dev"
region: "us-west"
config: "/config/dev/us-west.yaml"
---
environment: "dev"
region: "eu-west"
config: "/config/dev/eu-west.yaml"
---
environment: "prod"
region: "us-west"
config: "/config/prod/us-west.yaml"
---
environment: "prod"
region: "eu-west"
config: "/config/prod/eu-west.yaml"
```

## Null Filling

Dimensions not present in all jobs are filled with `null`:

```yaml
# Input
matrix:
  os: [linux, windows]
  include:
    - special: "custom"

# Output includes:
# os: linux,   special: null
# os: windows, special: null
# os: null,    special: "custom"
```

## Error Handling

### Invalid Matrix Format

```yaml
matrix: "not a map"
# Error: matrix must be a map, got string
```

### Invalid Dimension

```yaml
matrix:
  os: "not an array"
# Error: matrix dimension 'os' must be an array, got string
```

### Invalid Include/Exclude

```yaml
matrix:
  os: [linux, windows]
  include: "not a list"
# Error: include must be an array, got string
```

### Invalid Include Item

```yaml
matrix:
  os: [linux, windows]
  include:
    - "not a map"
# Error: include[0] must be a map, got string
```

## Performance

- **Time Complexity**: O(d₁ × d₂ × ... × dₙ) where dᵢ is dimension i size
- **Space Complexity**: O(jobs) for storing combinations
- **Exclude Filtering**: O(jobs × excludes)
- **Include Merging**: O(jobs × includes)

For typical use cases:
- 3 dimensions with 3 values each: 27 combinations
- 4 dimensions with 3 values each: 81 combinations
- 5 dimensions with 2 values each: 32 combinations

## Best Practices

### Organize Dimensions

```yaml
# Good - semantic grouping
matrix:
  environment: [dev, staging, prod]
  region: [us, eu, asia]
  architecture: [x86_64, arm64]
```

### Use Meaningful Variable Names

```yaml
# Good
matrix:
  database_version: [13, 14, 15]
  os: [ubuntu-20, ubuntu-22]

# Less clear
matrix:
  d1: [13, 14, 15]
  d2: [ubuntu-20, ubuntu-22]
```

### Document Exclusions

```yaml
matrix:
  os: [ubuntu, windows, macos]
  arch: [x86_64, arm64]
  exclude:
    # Windows doesn't support ARM64 builds
    - os: windows
      arch: arm64
```

### Use Includes for Special Cases

```yaml
matrix:
  os: [ubuntu, windows]
  include:
    # Add special docker-based job
    - os: docker
      arch: arm64
      image: "alpine:latest"
```

## See Also

- [For Loops (`for:`)](/docs/handlers/for.md) - iterate over arrays
- [Conditionals (`if:`)](/docs/handlers/if.md) - filter matrix jobs
- [Interpolation (`${}`)](/docs/handlers/interpolation.md) - use matrix variables
- [Syntax Reference](/docs/syntax.md) - complete directive syntax guide
