### Simple matrix

A top level matrix will produce multiple documents for iteration.

**Input:**

```yaml
matrix:
  os: [linux]
  arch: [x86_64, arm64]
  version: [18, 20]

name: "${os}-${arch}-v${version}"
```

**Output:**

```yaml
name: "linux-x86-18"
```
