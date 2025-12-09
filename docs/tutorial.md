# yamlexpr Tutorial

Welcome! This tutorial will introduce you to yamlexpr and show you how to use it to build powerful YAML composition workflows.

## What is yamlexpr?

yamlexpr is a lightweight YAML evaluation engine designed to reduce repetition and enable composition in YAML-based configuration files. Think of it as a minimal templating language that lets you:

- **Compose** multiple YAML files together
- **Conditionally render** sections based on boolean expressions
- **Loop** over data to generate repeated structures
- **Interpolate** variables into string values

### How Does It Compare to jsonnet?

[jsonnet](https://jsonnet.org/) is a powerful, complete programming language that can generate JSON. yamlexpr is intentionally simpler and more lightweight. While jsonnet is designed to be a full-featured configuration language, yamlexpr targets a specific use case: **reducing boilerplate in existing YAML workflows**.

Instead of rewriting your configuration in a new language, yamlexpr lets you enhance your existing YAML with just three directives: `include`, `for`, and `if`.

### Real-World Use Cases

**Example 1: GitHub Actions Workflows**

You have 10 different test environments (Python 3.8, 3.9, 3.10, 3.11, 3.12; on Linux and macOS). Instead of duplicating the entire test job for each combination, use `for` and `if` to generate them:

```yaml
jobs:
  test:
    - for: (python, os) in test_matrix
      if: python != "3.8" || os != "macos"
      name: "Test Python ${python} on ${os}"
      runs-on: "${{ matrix.os }}"
      with:
        python-version: "${python}"
```

**Example 2: Kubernetes Manifests**

You maintain microservices across staging and production. Instead of separate manifest files, use `include` to pull in common defaults and `if` to include environment-specific configuration:

```yaml
include: _base-deployment.yaml
spec:
  replicas: "${replicas}"
  env:
    - include: _base-env.yaml
    - if: "${is_production}"
      include: _prod-specific-env.yaml
```

**Example 3: Docker Compose**

You want a Docker Compose file that works for both development (with volumes) and CI (without). Use conditionals:

```yaml
services:
  app:
    if: "${enable_service}"
    image: "myapp:${version}"
    volumes:
      - if: "${enable_volumes}"
        ./src:/app/src
```

## Getting Started

### 1. Basic Interpolation

The simplest use case: substitute values into strings using `${variable}` syntax:

```go
expr := yamlexpr.New(os.DirFS("."))

data := map[string]any{
	"version":     "1.0.0",
	"environment": "production",
}

template := map[string]any{
	"app_name":    "myapp",
	"app_version": "${version}",
	"env":         "${environment}",
}

result, err := expr.Process(template, data)
// Result: {"app_name": "myapp", "app_version": "1.0.0", "env": "production"}
```

### 2. Conditionals with `if`

Include or exclude configuration blocks based on conditions:

```yaml
database:
  if: "${use_database}"
  host: "localhost"
  port: 5432
```

When `use_database` is true, the `database` section is included. When false, it's completely removed (not included as null or empty).

### 3. Loops with `for`

Generate repeated structures by iterating over arrays. There are two syntax forms:

**String Syntax (reference a variable):**

```yaml
services:
  - for: service in service_list
    name: "${service.name}"
    port: "${service.port}"
```

**Inline Array Syntax (literal data):**

```yaml
services:
  - for:
      - name: "api"
        port: 3000
      - name: "worker"
        port: 3001
    name: "${item.name}"
    port: "${item.port}"
```

For each item, this generates a separate element with the data. When using inline arrays, the loop variable defaults to `item`.

### 4. Composition with `include`

Pull in external YAML files to reduce duplication:

**base-config.yaml:**

```yaml
app:
  name: "myapp"
  version: "1.0.0"
  timeout: 30
```

**config.yaml:**

```yaml
include: base-config.yaml
app:
  debug: true
```

Result: The `app` section is merged, preserving values from base-config.yaml and adding the debug flag.

## Real Examples from the Test Suite

Since yamlexpr is thoroughly tested, we can examine real examples from the test fixtures:

### Example: Filtering with Conditionals

This example (from `testdata/fixtures/050-for-with-if-filtering.yaml`) shows how to iterate over items and filter based on a condition:

**Input:**

```yaml
enabled_services:
  - for:
      - name: "api"
        active: true
      - name: "worker"
        active: false
      - name: "scheduler"
        active: true
    if: item.active
    service: "${item.name}"
```

**Result:**

```yaml
enabled_services:
  - service: "api"
  - service: "scheduler"
```

Only services with `active: true` are included. The `worker` service is filtered out.

**Note on Syntax:** This example uses inline array syntax for `for`. You can also use string syntax: `for: "service in services"`. Both forms are equivalent.

### Example: Combining Everything

Here's a more complex example combining variables, includes, loops, and conditionals:

**Input:**

```yaml
db_host: localhost
db_port: 5432
enable_debug: true
server_list:
  - name: web-1
    ip: 10.0.1.1
    enabled: true
  - name: web-2
    ip: 10.0.1.2
    enabled: false
  - name: web-3
    ip: 10.0.1.3
    enabled: true

servers:
  - for: item in server_list
    if: item.enabled
    name: "${item.name}"
    ip: "${item.ip}"
debug:
  if: ${enable_debug}
  level: "verbose"
```

**Result:**

```yaml
servers:
  - name: "web-1"
    ip: "10.0.1.1"
  - name: "web-3"
    ip: "10.0.1.3"
debug:
  level: "verbose"
```

Notice:
- Variables at the root level (`db_host`, `db_port`, etc.) become available for interpolation and expressions
- The `for` loop generates items by iterating `server_list`
- The `if: item.enabled` condition filters out `web-2`
- The `debug` section uses `${enable_debug}` as a condition

## Advanced: Custom Directive Syntax

By default, yamlexpr uses `if`, `for`, and `include` as directives. But you can customize these keywords to match your framework's conventions:

```go
// Vue.js-style directives
expr := yamlexpr.New(os.DirFS("."), yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:      "v-if",
	For:     "v-for",
	Include: "v-include",
}))
```

Now your YAML uses Vue-style syntax:

```yaml
config:
  v-if: "${is_production}"
  v-include: "_prod-config.yaml"

services:
  - v-for: service in services
    name: "${service.name}"
```

See [Custom Syntax Configuration](custom-syntax.md) for more details.

## Real-World Example: Build Configuration

Let's look at a practical example: a build configuration for a project that supports multiple languages and deployment targets.

**config.json (data):**

```json
{
  "languages": [
    {"name": "go", "version": "1.21"},
    {"name": "python", "version": "3.11"},
    {"name": "node", "version": "18"}
  ],
  "deploy_to_staging": true,
  "deploy_to_prod": false
}
```

**build.yaml (template):**

```yaml
include: _base-build-settings.yaml

builds:
  - for: lang in languages
    name: "build-${lang.name}"
    runtime: "${lang.name}:${lang.version}"
    cache: true

deployment:
  staging:
    if: "${deploy_to_staging}"
    environment: "staging"
    replicas: 2

  production:
    if: "${deploy_to_prod}"
    environment: "production"
    replicas: 5
    
  if: "${deploy_to_staging}"
  monitoring:
    enabled: true
    alert_threshold: 80
```

**Result:**

```yaml
# (base settings merged in...)

builds:
  - name: "build-go"
    runtime: "go:1.21"
    cache: true
  - name: "build-python"
    runtime: "python:3.11"
    cache: true
  - name: "build-node"
    runtime: "node:18"
    cache: true

deployment:
  staging:
    environment: "staging"
    replicas: 2
  monitoring:
    enabled: true
    alert_threshold: 80
```

Notice how the `production` deployment is completely omitted because `deploy_to_prod` is false.

## Tips and Best Practices

### 1. Use Root-Level Variables Strategically

Root-level keys become available as variables. Use this to keep your configuration DRY:

```yaml
app_version: "1.2.3"
env: "production"

services:
  - name: "api"
    image: "myapp:${app_version}"
    environment: "${env}"
  - name: "worker"
    image: "myapp:${app_version}"
    environment: "${env}"
```

### 2. Leverage Nested Conditionals

You can nest `if` conditions inside `for` loops and vice versa:

```yaml
deployments:
  - for: env in environments
    if: env.enabled
    environment: "${env.name}"
    replicas: 
      - for: zone in env.zones
        if: zone.active
        zone: "${zone.name}"
        count: "${zone.count}"
```

### 3. Use Meaningful Variable Names

When using `for` loops, prefer descriptive variable names:

```yaml
# Good
- for: database in databases
  name: "${database.name}"
  
# Less clear
- for: item in databases
  name: "${item.name}"
```

### 4. Compose with Includes

Break your configuration into logical pieces and include them:

```yaml
include: _defaults.yaml
include: _networking.yaml
include: _monitoring.yaml

app:
  name: "myapp"
  version: "1.0"
```

## Next Steps

- Read the [Syntax Reference](syntax.md) for detailed documentation of all features
- Check out the [API Reference](api.md) for Go usage details
- Explore [Custom Syntax Configuration](custom-syntax.md) if you want different directive keywords
- Review the [Design Document](DESIGN.md) to understand the architecture

## Questions?

If you encounter unexpected behavior, check the [testing coverage](testing-coverage.md) to see what scenarios are tested, and refer to the [test fixtures](../testdata/fixtures) for examples of how different features work together.

Happy templating! 🎉
