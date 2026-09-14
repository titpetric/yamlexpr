## Simple matrix

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

Rendering produces **4** documents:

```yaml
arch: "x86_64"
name: "linux-x86_64-v18"
os: "linux"
version: 18
```

```yaml
arch: "x86_64"
name: "linux-x86_64-v20"
os: "linux"
version: 20
```

```yaml
arch: "arm64"
name: "linux-arm64-v18"
os: "linux"
version: 18
```

```yaml
arch: "arm64"
name: "linux-arm64-v20"
os: "linux"
version: 20
```

## Simple matrix

When a matrix is used in a list item, the values of the iteration are carried forward.

**Input:**

```yaml
- matrix:
    os: [linux]
    arch: [x86_64, arm64]
    version: [18, 20]
  name: "${os}-${arch}-v${version}"
```

**Output:**

```yaml
- name: "linux-x86-18"
  os: "linux"
  arch: "x86_64"
  version: 18
- name: "linux-arm64-18"
  os: "linux"
  arch: "arm64"
  version: 18
- name: "linux-x86-20
  os: "linux"
  arch: "x86"
  version: 20
- name: "linux-arm64-20"
  os: "linux"
  arch: "arm64"
  version: 20
```
