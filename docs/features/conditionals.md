## If Condition True

When if: evaluates to true, the block is included in output.

**Input:**

```yaml
config:
  debug:
    if: true
    enabled: true
    level: "verbose"
```

**Output:**

```yaml
config:
  debug:
    enabled: true
    level: "verbose"
```

## If Condition False

When if: evaluates to false, the block is omitted from output.

**Input:**

```yaml
config:
  debug:
    if: false
    enabled: true
    level: "verbose"
  production: true
```

**Output:**

```yaml
config:
  production: true
```

## If with Variable Reference

Use variables from the context in if conditions.

**Input:**

```yaml
environment: "production"
config:
  backup:
    if: ${environment == "production"}
    enabled: true
    retention_days: 30
  debug:
    if: ${environment == "development"}
    enabled: true
```

**Output:**

```yaml
environment: "production"
config:
  backup:
    enabled: true
    retention_days: 30
```

## If Condition with For Loop

Combine if: and for: to conditionally filter items during iteration.

**Input:**

```yaml
active_ports:
  - for: port in all_ports
    if: ${port.active}
    name: "${port.name}"
    number: ${port.number}
```

**Output:**

```yaml
active_ports:
  - name: "http"
    number: 80
  - name: "https"
    number: 443
```

## Nested If Conditions

If conditions can be nested at multiple levels.

**Input:**

```yaml
debug_enabled: true
detailed_logging: true
config:
  logging:
    if: ${debug_enabled}
    level: "debug"
    verbose:
      if: ${detailed_logging}
      include_timestamps: true
      include_caller: true
```

**Output:**

```yaml
debug_enabled: true
detailed_logging: true
config:
  logging:
    level: "debug"
    verbose:
      include_timestamps: true
      include_caller: true
```
