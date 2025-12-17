# Matrix Expansion with `matrix:`

## Syntax Cheat Sheet

```yaml
# Simple matrix: creates cartesian product of dimensions
matrix:
  os: [linux, macos, windows]
  arch: [x86_64, arm64]
  version: [18, 20]
name: "${os}-${arch}-v${version}"

# With variables (non-array values)
matrix:
  platform: [ubuntu, fedora]
  timeout: 300
  retries: 3
name: "${platform}"
build_timeout: ${timeout}
max_retries: ${retries}

# With exclude: remove specific combinations
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
name: "${os}/${arch}"

# With include: add custom combinations
matrix:
  os: [linux, windows]
  arch: [x86_64]
  include:
    - os: macos
      arch: arm64
      xcode: "14"
name: "${os}/${arch}"
xcode: "${xcode}"
```

## Description

The `matrix:` directive generates a cartesian product of dimension combinations. This is essential for CI/CD systems that need to test across multiple platforms, versions, and configurations.

Unlike `for:` loops, matrix creates all possible combinations by default, with options to exclude or include specific combinations. This is inspired by GitHub Actions matrix strategy.

## Core Concepts

- **Dimensions**: Array values in the matrix map become dimensions
- **Variables**: Non-array values are variables added to each combination
- **Cartesian product**: By default, all combinations of dimensions are generated
- **Exclude**: Filter out specific combinations that shouldn't be generated
- **Include**: Add additional custom combinations beyond the cartesian product
- **Scope**: All dimension and variable values become available for interpolation

## Basic Matrix

The simplest matrix creates all combinations of dimensions:

**Input:**

```yaml
- matrix:
    os: [linux, windows]
    version: [12, 14]
    run: steps
  name: "Test ${os} v${version}"
  os: ${os}
  version: ${version}
```

**Output:**

```yaml
- name: "Test linux v12"
  os: linux
  version: 12
- name: "Test linux v14"
  os: linux
  version: 14
- name: "Test windows v12"
  os: windows
  version: 12
- name: "Test windows v14"
  os: windows
  version: 14
```

The matrix creates 2 × 2 = 4 combinations. The `run: steps` variable is available in all combinations.

## Three-Dimensional Matrix

**Input:**

```yaml
- matrix:
    language: [go, python, rust]
    version: ["1.0", "2.0"]
    os: [linux, windows]
  job_name: "Test ${language} v${version} on ${os}"
  language: ${language}
  version: "${version}"
  os: ${os}
```

**Output:**

```yaml
- job_name: "Test go v1.0 on linux"
  language: go
  version: "1.0"
  os: linux
- job_name: "Test go v1.0 on windows"
  language: go
  version: "1.0"
  os: windows
# ... 14 more combinations (3 × 2 × 2 = 12 total)
```

## With Exclude

Filter out specific combinations that shouldn't be generated:

**Input:**

```yaml
matrix:
  os: [linux, macos, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
    - os: macos
      arch: x86_64
name: "${os}/${arch}"
---
- arch: x86_64
  name: linux/x86_64
  os: linux
- arch: arm64
  name: linux/arm64
  os: linux
- arch: arm64
  name: macos/arm64
  os: macos
- arch: x86_64
  name: windows/x86_64
  os: windows
```

Without exclude, this would generate 3 × 2 = 6 combinations. The exclude section removes 2 combinations, leaving 4.

## With Include

Add custom combinations beyond the cartesian product:

**Input:**

```yaml
matrix:
  os: [linux, windows]
  arch: [x86_64]
  include:
    - os: macos
      arch: arm64
      xcode: "14"
name: "${os}/${arch}"
xcode: "${xcode}"
```

**Output:**

```yaml
- arch: x86_64
  name: linux/x86_64
  os: linux
  xcode: null
- arch: x86_64
  name: windows/x86_64
  os: windows
  xcode: null
- arch: arm64
  name: macos/arm64
  os: macos
  xcode: "14"
```

The include creates 3 total combinations: 2 from the cartesian product (linux/windows × x86_64) plus 1 custom (macos/arm64 with xcode).

## Combining Exclude and Include

You can use both exclude and include together:

**Input:**

```yaml
matrix:
  os: [linux, macos, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
  include:
    - os: freebsd
      arch: amd64
      beta: true
name: "${os}/${arch}"
beta: "${beta}"
```

**Output:**

```yaml
# From cartesian product (3 × 2 = 6 minus 1 exclude = 5)
- arch: x86_64
  name: linux/x86_64
  os: linux
  beta: null
- arch: arm64
  name: linux/arm64
  os: linux
  beta: null
- arch: x86_64
  name: macos/x86_64
  os: macos
  beta: null
- arch: arm64
  name: macos/arm64
  os: macos
  beta: null
- arch: x86_64
  name: windows/x86_64
  os: windows
  beta: null

# From include
- arch: amd64
  name: freebsd/amd64
  os: freebsd
  beta: true
```

## Matrix with Variables (Non-Dimension Values)

Non-array values in the matrix become variables available in all combinations:

**Input:**

```yaml
matrix:
  os: [ubuntu, fedora]
  timeout: 300
  retries: 3
name: "${os}"
build_timeout: ${timeout}
max_retries: ${retries}
```

**Output:**

```yaml
- build_timeout: 300
  max_retries: 3
  name: ubuntu
- build_timeout: 300
  max_retries: 3
  name: fedora
```

## Common Use Cases

- **CI/CD test matrices**: Test on multiple platforms, versions, architectures
- **Build configurations**: Generate builds for different targets
- **Cross-platform testing**: Create jobs for Linux, macOS, Windows variants
- **Version compatibility**: Test against multiple language/framework versions
- **Environment variations**: Combine different regions, zones, or deployment targets
- **Hardware configurations**: Generate configs for different CPU architectures

## Comparison with For Loops

| Feature          | For Loop                  | Matrix                             |
|------------------|---------------------------|------------------------------------|
| **Source**       | Single array variable     | Multiple dimension arrays          |
| **Combinations** | Linear iteration          | Cartesian product                  |
| **Filtering**    | Use `if:` with conditions | Use `exclude:` section             |
| **Custom items** | Requires separate array   | Use `include:` section             |
| **Use case**     | Iterate known collection  | Generate all platform combinations |

**For loop**: `for: item in items` - 5 items = 5 results **Matrix**: `matrix: {a: [1,2], b: [x,y]}` - 2 × 2 = 4 results

## Edge Cases

### Empty Dimensions

An empty dimension array produces no combinations:

```yaml
matrix:
  os: []
  version: [1, 2]
name: "${os}-v${version}"
```

Result: Empty array (no combinations)

### Single Item Dimensions

Matrix works fine with single-item dimensions:

```yaml
matrix:
  language: [go]
  version: [1.19, 1.20]
name: "${language} v${version}"
```

Result: 1 × 2 = 2 combinations

### Complex Values in Include

Include entries can have nested structures:

```yaml
matrix:
  os: [linux]
  include:
    - os: windows
      env:
        key1: value1
        key2: value2
name: "${os}"
```
