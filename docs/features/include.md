### Basic Include

Include external YAML file and merge contents.

**Input:**

```yaml
database:
  host: "localhost"
  port: 5432
  credentials:
    include: "_db-credentials.yaml"
```

**Output:**

```yaml
database:
  host: "localhost"
  port: 5432
  credentials:
    username: "dbuser"
    password: "secret123"
```

### Multiple Includes

Include multiple files to compose configuration from parts.

**Input:**

```yaml
metadata:
  include: "_common-labels.yaml"
database:
  credentials:
    include: "_db-credentials.yaml"
```

**Output:**

```yaml
metadata:
  app: "myapp"
  version: "1.0.0"
  managed_by: "yamlexpr"
database:
  credentials:
    username: "dbuser"
    password: "secret123"
```

### Include in For Loop

Use includes within for loop templates to compose repeated sections.

**Input:**

```yaml
environments:
  - for: env in ["staging", "production"]
    name: "${env}"
    database:
      credentials:
        include: "_db-credentials.yaml"
```

**Output:**

```yaml
environments:
  - name: "staging"
    database:
      credentials:
        username: "dbuser"
        password: "secret123"
  - name: "production"
    database:
      credentials:
        username: "dbuser"
        password: "secret123"
```
