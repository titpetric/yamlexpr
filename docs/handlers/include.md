# Include Handler (`include:`)

The **`include:` directive** loads and merges external YAML files into the current document structure.

## Overview

The `include:` directive enables YAML composition by loading external files and merging their content into the current document. Files are processed in order, with later files overriding earlier values for duplicate keys.

## Syntax

### Single File

```yaml
include: "base.yaml"
```

### Multiple Files

```yaml
include:
  - "base.yaml"
  - "overrides.yaml"
```

## File Resolution

Files are resolved relative to the **base directory** provided to `yamlexpr.New()`:

```go
expr := yamlexpr.New(yamlexpr.WithFS(os.DirFS("./config")))
// "base.yaml" resolves to "./config/base.yaml"
// "subdir/overrides.yaml" resolves to "./config/subdir/overrides.yaml"
```

## Merge Behavior

### Key Merging

When embedding files with overlapping keys:
- **Last file wins**: Later embeds override earlier values
- **Deep merge**: Keys are merged recursively for maps
- **Array handling**: Arrays are **replaced** (not merged)

### Example

```yaml
# base.yaml
database:
  host: "localhost"
  port: 5432
features: ["auth", "cache"]

# config.yaml
embed: "base.yaml"
database:
  port: 3306      # Override
  username: "app" # New key
features: ["api"] # Replace array

# Result
database:
  host: "localhost"
  port: 3306      # From config.yaml
  username: "app" # From config.yaml
features: ["api"] # From config.yaml (replaced)
```

## Processing Order

1. **Include files are loaded first** (highest priority: 1000)
2. **Files are merged in order** (left-to-right)
3. **Current document keys are merged last** (highest precedence)
4. **Other directives are processed** (if, for, etc.)

```yaml
# Load files first
embed: "base.yaml"
embed:
  - "defaults.yaml"
  - "overrides.yaml"

# Current keys take precedence
key: "from_current"
```

## Single File Include

```yaml
# Input - config.yaml
include: "base.yaml"
debug: true

# Input - base.yaml
database:
  host: "localhost"
server:
  port: 8080

# Output
database:
  host: "localhost"
server:
  port: 8080
debug: true
```

## Multiple File Include

```yaml
# Input - config.yaml
include:
  - "base.yaml"
  - "env.yaml"

# Input - base.yaml
database:
  host: "localhost"
  pool_size: 10
server:
  port: 8080

# Input - env.yaml (production overrides)
database:
  host: "db.example.com"
  pool_size: 50

# Output
database:
  host: "db.example.com"  # From env.yaml
  pool_size: 50           # From env.yaml
server:
  port: 8080              # From base.yaml
```

## Conditional Including

Use `include:` with conditions to conditionally load files:

```yaml
# config.yaml
include: "base.yaml"

env: ${environment}

debug_settings:
  if: "${env == 'dev'}"
  include: "debug.yaml"

prod_settings:
  if: "${env == 'prod'}"
  include: "prod.yaml"
```

## Including in Arrays

```yaml
# config.yaml
services:
  - include: "service_api.yaml"
  - include: "service_worker.yaml"
  - name: "cache"
    port: 6379

# Result
services:
  - # keys from service_api.yaml
  - # keys from service_worker.yaml
  - name: "cache"
    port: 6379
```

## Recursive Including

Included files can themselves use `include:` directives:

```yaml
# main.yaml
include: "config.yaml"

# config.yaml
include: "base.yaml"
database:
  host: "localhost"

# base.yaml
environment: "dev"
version: "1.0"

# Final result contains all merged keys from the chain
```

## Circular Reference Protection

Circular references (A → B → A) are detected and prevented:

```go
// If detected, returns error:
// "circular include detected: main.yaml -> config.yaml -> main.yaml"
```

## Variable Interpolation in Include Paths

Include paths can use variable interpolation:

```yaml
env: "production"
environment_file: "config-${env}.yaml"

config:
  include: "${environment_file}"
```

## Error Handling

### File Not Found

```yaml
include: "missing.yaml"
# Error: error reading file missing.yaml: file does not exist
```

### Invalid YAML

```yaml
include: "invalid.yaml"
# Error: error parsing YAML file invalid.yaml: yaml: line 1: mapping values are not allowed in this context
```

### Circular Reference

```yaml
include: "self.yaml"  # self.yaml includes itself
# Error: circular include detected: config.yaml -> self.yaml -> config.yaml
```

### Non-String Path

```yaml
include: 123
# Error: include must be a string or list of strings, got int
```

## Deep Merging

Maps are merged recursively, while other types are replaced:

```yaml
# Input - base.yaml
server:
  http:
    port: 8080
    timeout: 30
  https:
    enabled: false

# Input - overrides.yaml (included by base.yaml)
server:
  http:
    port: 3000    # Override

# Result
server:
  http:
    port: 3000    # From overrides
    timeout: 30   # From base
  https:
    enabled: false # From base
```

## Best Practices

### File Organization

```
config/
├── base.yaml          # Common defaults
├── dev.yaml           # Development overrides
├── prod.yaml          # Production overrides
├── schemas/
│   ├── database.yaml
│   └── server.yaml
└── secrets/
    └── credentials.yaml
```

### Use Cases

**Shared Defaults:**

```yaml
# base.yaml - shared across all environments
database:
  pool_size: 10
  timeout: 5
server:
  keepalive: true
```

**Environment Overrides:**

```yaml
# prod.yaml
include: "base.yaml"
database:
  pool_size: 100
  timeout: 30
```

**Feature Flags:**

```yaml
include:
  - "base.yaml"
  
analytics:
  if: "${enable_analytics}"
  include: "features/analytics.yaml"
```

**Multi-Stage Composition:**

```yaml
# final.yaml
include:
  - "defaults.yaml"        # Stage 1: baseline
  - "environment.yaml"     # Stage 2: environment-specific
  - "overrides.yaml"       # Stage 3: local overrides
```

## Implementation Details

### Priority

The `include:` directive has **highest priority (1000)** in the handler system, ensuring:
1. Files are loaded and merged first
2. Other directives see merged content
3. Interpolation uses merged values

### Merge Algorithm

```
result = load(include_files) → merge left-to-right → merge current keys
```

### Filesystem Access

- All file operations use the provided `fs.FS` interface
- Paths are always relative to the filesystem root
- Absolute paths (starting with `/`) are resolved within the filesystem

## Examples

### Configuration Base + Environment Override

```yaml
# base.yaml
app:
  name: "myapp"
  version: "1.0"
database:
  host: "localhost"
  port: 5432
server:
  port: 8080

# production.yaml
include: "base.yaml"
database:
  host: "db.prod.internal"
  pool_size: 50
server:
  port: 443

# Result
app:
  name: "myapp"
  version: "1.0"
database:
  host: "db.prod.internal"  # Overridden
  port: 5432
  pool_size: 50             # Added
server:
  port: 443                 # Overridden
```

### Shared Schema Components

```yaml
# schemas/database.yaml
type: "postgresql"
version: "14"
ssl_required: true

# schemas/redis.yaml
type: "redis"
version: "7"

# config.yaml
services:
  db:
    include: "schemas/database.yaml"
    host: "localhost"
  cache:
    include: "schemas/redis.yaml"
    host: "localhost"

# Result
services:
  db:
    type: "postgresql"
    version: "14"
    ssl_required: true
    host: "localhost"
  cache:
    type: "redis"
    version: "7"
    host: "localhost"
```

### Conditional Feature Including

```yaml
# main.yaml
include: "base.yaml"
features:
  analytics:
    if: "${enable_analytics}"
    include: "features/analytics.yaml"
  monitoring:
    if: "${enable_monitoring}"
    include: "features/monitoring.yaml"

# features/analytics.yaml
provider: "mixpanel"
events:
  - "page_view"
  - "user_signup"

# Output (if both features enabled)
features:
  analytics:
    provider: "mixpanel"
    events:
      - "page_view"
      - "user_signup"
  monitoring:
    # monitoring keys...
```

## See Also

- [Syntax Reference](/docs/syntax.md) - complete directive syntax guide
- [API Reference](/docs/api.md) - file loading and merging API
- [For Loops (`for:`)](/docs/handlers/for.md) - iterate with included schemas
- [Conditionals (`if:`)](/docs/handlers/if.md) - conditional including
