#!/bin/bash
# Generates include/composition feature documentation
# Usage: ./include.md.sh > include.md

cat << 'EOF'
# Composition with `include:`

## Syntax Cheat Sheet

```yaml
# Include single file at top level
include: "_base.yaml"
config:
  name: "myapp"

# Include in nested structure
database:
  include: "_db-config.yaml"
  pool_size: 10

# Include multiple files as array
imports:
  include:
    - "_monitoring.yaml"
    - "_logging.yaml"
  environment: production

# Inline include (same key name)
services:
  include: "_services.yaml"
  timeout: 30
```

## Description

The `include:` directive enables composition by merging external YAML files into the current document. This allows reusable components, shared configurations, and modular YAML structures.

Files are resolved relative to the base directory (filesystem) provided to `Expr.New()`. Includes can appear at any level and combine with other directives like `for:` and `if:`.

## Core Concepts

- **File resolution**: Files are resolved relative to the base filesystem directory
- **Merging**: Included content replaces the `include:` directive at that location
- **Composition**: Can be combined with for loops, conditionals, and other features
- **Reusability**: Share common configurations across multiple files
- **Nesting**: Includes can appear at any depth in the document structure

## Basic Include

**Files:**

`_base.yaml`:
```yaml
version: "1.0"
metadata:
  author: "team"
  created: "2024-01-01"
```

`config.yaml`:
```yaml
include: "_base.yaml"
app:
  name: "myapp"
  port: 8080
```

**Result:**
```yaml
version: "1.0"
metadata:
  author: "team"
  created: "2024-01-01"
app:
  name: "myapp"
  port: 8080
```

## Nested Includes

Includes work at any level in the document structure:

**Files:**

`_db-defaults.yaml`:
```yaml
host: "localhost"
port: 5432
timeout: 30
```

`config.yaml`:
```yaml
app:
  name: "myapp"
database:
  include: "_db-defaults.yaml"
  pool_size: 20
  ssl: true
```

**Result:**
```yaml
app:
  name: "myapp"
database:
  host: "localhost"
  port: 5432
  timeout: 30
  pool_size: 20
  ssl: true
```

## Multiple Includes

Include multiple files at once using an array:

**Files:**

`_monitoring.yaml`:
```yaml
prometheus:
  enabled: true
  port: 9090
```

`_logging.yaml`:
```yaml
logs:
  level: info
  format: json
```

`config.yaml`:
```yaml
app:
  include:
    - "_monitoring.yaml"
    - "_logging.yaml"
  name: "myapp"
```

**Result:**
```yaml
app:
  prometheus:
    enabled: true
    port: 9090
  logs:
    level: info
    format: json
  name: "myapp"
```

## Include with For Loops

Includes work well with for loops to create multiple instances with shared base config:

**Files:**

`_service-defaults.yaml`:
```yaml
timeout: 30
retries: 3
health_check: enabled
```

`config.yaml`:
```yaml
services:
  - for: svc in service_list
    name: "${svc.name}"
    port: ${svc.port}
    include: "_service-defaults.yaml"
service_list:
  - name: "api"
    port: 8080
  - name: "worker"
    port: 9000
```

**Result:**
```yaml
services:
  - name: "api"
    port: 8080
    timeout: 30
    retries: 3
    health_check: enabled
  - name: "worker"
    port: 9000
    timeout: 30
    retries: 3
    health_check: enabled
service_list:
  - name: "api"
    port: 8080
  - name: "worker"
    port: 9000
```

## Include with Conditionals

Includes can be conditional:

**Files:**

`_postgres.yaml`:
```yaml
driver: "postgresql"
port: 5432
```

`_mysql.yaml`:
```yaml
driver: "mysql"
port: 3306
```

`config.yaml`:
```yaml
app:
  name: "myapp"
database:
  if: ${enable_database}
  host: "db.local"
  include: "${use_postgres ? '_postgres.yaml' : '_mysql.yaml'}"
enable_database: true
use_postgres: true
```

Note: Variable interpolation in include paths is not yet supported, but conditional includes work:

```yaml
database:
  postgres:
    if: ${use_postgres}
    include: "_postgres.yaml"
  mysql:
    if: ${use_mysql}
    include: "_mysql.yaml"
```

## Building Reusable Components

Create a components directory structure:

```
config/
├── base/
│   ├── _database.yaml
│   ├── _cache.yaml
│   └── _logging.yaml
├── services/
│   ├── _api-service.yaml
│   ├── _worker-service.yaml
│   └── _scheduler-service.yaml
├── environments/
│   ├── _dev.yaml
│   ├── _staging.yaml
│   └── _production.yaml
└── app.yaml
```

**_database.yaml:**
```yaml
database:
  host: "localhost"
  port: 5432
  pool_size: 10
  timeout: 30
```

**_cache.yaml:**
```yaml
cache:
  type: "redis"
  host: "localhost"
  port: 6379
  ttl: 3600
```

**app.yaml:**
```yaml
include:
  - "base/_database.yaml"
  - "base/_cache.yaml"
app:
  name: "myapp"
  environment: ${APP_ENV}
services:
  - for: svc in service_names
    name: "${svc}"
    include: "services/_${svc}-service.yaml"
service_names:
  - "api"
  - "worker"
```

## Layered Configuration

Build configuration by layering includes:

**base.yaml:**
```yaml
version: "1.0"
defaults:
  timeout: 30
  retries: 3
```

**environment-specific.yaml:**
```yaml
include: "_base.yaml"
environment: production
debug: false
resources:
  cpu: "4"
  memory: "8Gi"
```

**service-config.yaml:**
```yaml
include: "_environment-specific.yaml"
services:
  include: "_services.yaml"
  replicas: 3
```

## Common Use Cases

- **Base configurations**: Shared settings used across multiple configs
- **Component libraries**: Reusable service definitions
- **Environment configs**: Layer base → environment → specific settings
- **Service templates**: Repeated service structures with shared defaults
- **Feature toggles**: Include different features based on conditions
- **Multi-tenant setup**: Share base config, customize per tenant

## File Resolution

Files are resolved relative to the filesystem root provided to `Expr.New()`:

```go
// Assuming directory structure:
// configs/
//   ├── _base.yaml
//   ├── _services.yaml
//   └── app.yaml

expr := yamlexpr.New(os.DirFS("configs"))
docs, err := expr.Load("app.yaml")
// Files are resolved relative to "configs" directory
```

When using `app.yaml` with `include: "_base.yaml"`, it resolves to `configs/_base.yaml`.

## Include Chain Prevention

Circular includes are detected and reported as errors:

```yaml
# a.yaml
include: "b.yaml"

# b.yaml
include: "a.yaml"  # ERROR: circular include detected
```

## Merging Behavior

When an include is processed:
1. The include file is loaded and processed (recursively)
2. The resulting content replaces the `include:` directive
3. All keys from the included file are merged at that location
4. Existing keys are preserved (included content doesn't override)

EOF
