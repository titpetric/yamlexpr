# Custom Syntax Configuration

By default, yamlexpr uses standard directive keywords: `if`, `for`, and `include`. However, you can customize these keywords to match your preferred syntax or framework conventions (e.g., Vue.js-style directives).

## Overview

The `Syntax` type defines custom directive keywords. Empty fields in the syntax struct retain their default values, allowing you to customize directives partially or fully.

```go
type Syntax struct {
	If      string `json:"if" yaml:"if"`           // e.g., "v-if" (default: "if")
	For     string `json:"for" yaml:"for"`         // e.g., "v-for" (default: "for")
	Include string `json:"include" yaml:"include"` // e.g., "v-include" (default: "include")
}
```

The `Syntax` struct includes JSON and YAML tags, making it easy to serialize/deserialize configuration from JSON or YAML files.

## Basic Usage

Configure custom syntax using `WithSyntax` when creating an Expr evaluator:

```go
import (
    "github.com/titpetric/yamlexpr"
)

// Use Vue-style directives
e := yamlexpr.New(myFS, yamlexpr.WithSyntax(yamlexpr.Syntax{
    If:      "v-if",
    For:     "v-for",
    Include: "v-include",
}))

// Now use the custom directives in YAML processing
doc := map[string]any{
    "items": []any{1, 2, 3},
}

template := map[string]any{
    "results": map[string]any{
        "v-for": "item in items",
        "value": "${item}",
        "v-if":  "item > 1",
    },
}

result, _ := e.Process(template, doc)
// result: {"results": [{"value": "2"}, {"value": "3"}]}
```

## Partial Customization

You can customize individual directives while keeping others at their defaults. Empty string values retain the default keywords:

```go
// Only customize the if directive
e := yamlexpr.New(myFS, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If: "v-if",
	// For remains "for"
	// Include remains "include"
}))

// Only customize for and include
e := yamlexpr.New(myFS, yamlexpr.WithSyntax(yamlexpr.Syntax{
	For:     "v-for",
	Include: "v-include",
	// If remains "if"
}))
```

## Examples

### Vue.js-Style Directives

```go
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:      "v-if",
	For:     "v-for",
	Include: "v-include",
}))
```

YAML using Vue directives:

```yaml
include: _base.yaml

services:
  - v-for: service in services
    name: "${service.name}"
    v-if: "${service.enabled}"
    config:
      v-include: "_service-defaults.yaml"
```

### Custom Prefix (Angular-style)

```go
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:      "*ngIf",
	For:     "*ngFor",
	Include: "*ngInclude",
}))
```

YAML using Angular directives:

```yaml
*ngInclude: _base.yaml

items:
  - *ngFor: item in list
    name: "${item}"
    *ngIf: "item.active"
```

### Minimal Customization

```go
// Only change the for directive to "each"
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	For: "each",
}))
```

YAML mixing default and custom directives:

```yaml
include: _base.yaml  # uses default "include"

items:
  - each: item in list
    name: "${item}"
    if: "item.active"  # uses default "if"
```

## All Features with Custom Syntax

All yamlexpr features work identically with custom syntax. The directive keywords are the only difference.

### Variable Interpolation

Variable interpolation with `${variable}` syntax remains unchanged:

```yaml
name: "${user.name}"
email: "${user.email}"
```

### Conditionals with Custom `if`

```yaml
config:
  v-if: "${is_production}"
  debug: false
  timeout: 30
```

### Loops with Custom `for`

```yaml
servers:
  - v-for: server in server_list
    hostname: "${server.name}"
    port: "${server.port}"
```

### Nested Structures

```yaml
v-include: _base.yaml

services:
  - v-for: service in services
    v-if: "${service.enabled}"
    name: "${service.name}"
    v-include: "_service-config.yaml"
    resources:
      v-for: resource in service.resources
      type: "${resource.type}"
      v-if: "resource.required"
```

## Design Rationale

Custom syntax configuration uses a typed `Syntax` struct with `WithSyntax()` functional option because:

1. **Type Safety**: The struct provides clear, documented fields with IDE autocomplete
2. **Partial Configuration**: Empty fields automatically preserve defaults, simplifying partial customization
3. **Flexibility**: Struct-based approach scales better if additional configuration options are needed in the future
4. **Clarity**: The `Syntax` struct name makes it explicit what is being configured

## Migration from Default Syntax

If you decide to switch from default to custom syntax:

1. Update your `New()` call with `WithSyntax()`
2. Replace old directive keywords with new ones in your YAML templates
3. No changes needed to interpolation syntax or expression evaluation

```go
// Before
e := yamlexpr.New(fs)

// After
e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
	If:      "v-if",
	For:     "v-for",
	Include: "v-include",
}))
```

## Serializing Configuration

Since the `Syntax` struct includes JSON and YAML tags, you can easily load configuration from files:

### From JSON

```go
import (
    "encoding/json"
    "os"
    "github.com/titpetric/yamlexpr"
)

// Load syntax from JSON file
data, _ := os.ReadFile("syntax.json")
var syntax yamlexpr.Syntax
json.Unmarshal(data, &syntax)

e := yamlexpr.New(fs, yamlexpr.WithSyntax(syntax))
```

JSON file example:

```json
{
  "if": "v-if",
  "for": "v-for",
  "include": "v-include"
}
```

### From YAML

```go
import (
    "gopkg.in/yaml.v3"
    "os"
    "github.com/titpetric/yamlexpr"
)

// Load syntax from YAML file
data, _ := os.ReadFile("syntax.yaml")
var syntax yamlexpr.Syntax
yaml.Unmarshal(data, &syntax)

e := yamlexpr.New(fs, yamlexpr.WithSyntax(syntax))
```

YAML file example:

```yaml
if: v-if
for: v-for
include: v-include
```

## Performance

Custom syntax configuration has no runtime performance impact. The directive keywords are resolved once during `New()` initialization.
