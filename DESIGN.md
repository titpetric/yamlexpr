# yaml-expr Design Document

## Overview

`yaml-expr` is a minimalist YAML composition and expression language for Go. It provides variable interpolation, conditionals, file composition, and loop expansion for YAML documents.

## Design Philosophy

### 1. Minimalism First

The stack package contains **only** the essential functionality needed for variable lookup and resolution:

- No external dependencies (except stdlib)
- Path-based resolution with caching (vuego/lessgo pattern)
- Type-safe accessors for common types (GetString, GetInt, GetSlice, GetMap)
- Pool-based map allocation for GC efficiency

### 2. Reusable Stack Package

The `stack` package is intentionally generic and dependency-free:

- Can be imported by vuego, lessgo, yaml-expr, or any other project
- Uses only Go standard library
- Caches path resolution to avoid repeated parsing
- Provides stack-based scoping with push/pop semantics

**Key difference from vuego's Stack:**
- Removed struct field resolution (vuego-specific feature)
- Removed rootData parameter (not needed for YAML)
- Kept all core functionality: lookup, resolve, GetString/GetInt/GetSlice/GetMap

### 3. Expression Package Separation

The `expr` package handles YAML-specific concerns:

- Document processing (maps and slices)
- Interpolation (${varname} syntax)
- Directives (planned: include, for, if)
- File composition (planned)

Kept separate from `stack` to maintain reusability.

## Package Organization

```
yaml-expr/
├── stack/           # Reusable variable scope stack
│   ├── stack.go
│   └── stack_test.go
├── expr/            # YAML expression evaluation
│   ├── expr.go
│   ├── expr_test.go
│   ├── interpolate.go
│   └── interpolate_test.go
├── testdata/
│   └── fixtures/    # Test fixtures with expected output
└── AGENTS.md        # Development conventions
```

## Feature Roadmap

### Done ✓

1. **Stack-based variable scoping**
   - Push/pop for nested scopes
   - Variable shadowing
   - Path resolution (user.name, items[0].title)
   - Type conversions (GetString, GetInt, GetSlice, GetMap)

2. **Variable interpolation**
   - ${varname} syntax in strings
   - Nested path support (${obj.path})
   - Missing variable handling (returns placeholder)

### Waiting

3. **Include composition**
   - Load external YAML files
   - Merge into current document
   - Relative path resolution from fs.FS

4. **For loops**
   - Iterate over arrays
   - Loop variable in scope
   - Support for objects

5. **If conditions**
   - go-expr integration for expression evaluation
   - Include/exclude blocks based on condition
   - Support for comparisons and logical operators

## Design Decisions

### Why Pool-Based Maps?

Reduces GC pressure when processing large documents with many scopes. Matching vuego's approach for consistency.

### Why Separate Stack from Expr?

The Stack package solves a general problem (variable scoping) useful across projects. The Expr package is YAML-specific and can import Stack.

### Why ${} for Interpolation?

- Unambiguous in YAML syntax
- Matches template languages (Vue, Django)
- Easy to parse with regex
- Clear visual distinction from YAML syntax

### Why Caching Path Resolution?

Path parsing (user.name → [user, name]) is expensive. The cache (limited to 256 entries) significantly improves performance for repeated lookups of the same paths.

### Why Fixture-Based Testing?

Follows lessgo pattern:
- Source of truth is on disk
- Easy to add test cases
- Clear input/output separation
- Works with any serialization format (YAML, JSON, etc)

## Integration Points

### With go-expr

For evaluating `if:` conditions, yaml-expr will import github.com/go-expr/go-expr:

```go
// Planned
if cond, ok := m["if"]; ok {
  env := st.All()  // Get all variables as map[string]any
  result := expr.Eval(cond, env)  // Evaluate expression
}
```

### With Other Projects

The `stack` package can be imported by:

```go
import "github.com/titpetric/yaml-expr/stack"

s := stack.New(map[string]any{"x": 1})
```

## Testing Strategy

### Black Box Testing

All tests use exported APIs only:
- `stack_test.go` tests Stack API
- `expr_test.go` tests Expr API
- `interpolate_test.go` tests interpolation

### Fixture-Based Tests

Will follow lessgo pattern with `---` delimiter:

```yaml
# Input
items:
  - ${name}
---
# Expected output
items:
  - John
```

Run with: `go test -v -run TestFixtures ./expr`

## Future Improvements

1. **YAML Parser Integration**
   - Replace parseYAML stub with gopkg.in/yaml.v3
   - Support both YAML input and output

2. **Expression Evaluation**
   - Integrate go-expr for if conditions
   - Support for complex expressions

3. **For Loop Implementation**
   - Handle `for: items` syntax
   - Loop variable (item, key) in stack

4. **Error Messages**
   - Include fixture path in errors
   - Better context for debugging

5. **Performance**
   - Benchmark fixture processing
   - Profile memory usage with large documents
