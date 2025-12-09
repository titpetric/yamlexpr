# yamlexpr Handlers

Custom directive handlers for yamlexpr, providing extensible YAML composition functionality.

## Built-in Handlers

### embed
File composition and merging. Loads external YAML files and merges their content.

**Priority:** 1000 (highest - runs first)
**Usage:** 
```yaml
embed: "base.yaml"
# or
embed:
  - "base.yaml"
  - "overrides.yaml"
```
**Automatically registered:** Enabled by default when creating an Expr instance.

### if
Conditional block inclusion. Omits a block if the condition is false.

**Priority:** 100
**Usage:** `if: "${condition}"` or `if: "expression > 5"`

### for
Loop expansion. Repeats a template for each item in a collection.

**Priority:** 100
**Usage:** `for: "item in items"` or `for: "(idx, item) in items"`

### interpolation
String variable interpolation. Replaces `${variable}` syntax with stack values.

**Priority:** 50 (automatic, always applied to strings)
**Usage:** Automatically applied to all string values during document processing.
```yaml
greeting: "Hello ${name}"  # Becomes "Hello Alice" if name=Alice
config: ${config_path}     # Becomes the value of config_path variable
```
**Lenient mode:** Uses `InterpolateString()` for lenient interpolation (undefined variables left unchanged)
**Strict mode:** Uses `InterpolateStringWithContext()` for strict interpolation (errors on undefined variables)

### discard
Simple conditional omit. Omits a block if set to true (opposite of if).

**Priority:** 10
**Usage:** `discard: true`

### matrix
GitHub Actions-compatible matrix expansion with exclude/include support.

**Priority:** 5
**Usage:**
```yaml
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
  include:
    - os: macos
      arch: arm64
```

## Handler Priority System

Higher priority runs first:
- 1000: `include` (file merging - core)
- 100: `if`, `for` (built-in)
- 10: `discard` (example custom)
- 5: `matrix` (example custom)
- 0: user-defined defaults

## Creating Custom Handlers

Implement `yamlexpr.DirectiveHandler`:

```go
type DirectiveHandler func(
    value any,
    mapContext map[string]any,
    ctx *ExprContext,
) (result any, consumed bool, error)
```

Register with:

```go
e := yamlexpr.New(fs,
    yamlexpr.WithDirectiveHandler("name", handler, priority),
)
```

Return values:
- `nil, true, nil` → Omit block, skip normal processing
- `map[string]any, true, nil` → Single item, skip normal processing
- `[]any, true, nil` → Multiple items, skip normal processing
- `nil, false, nil` → Continue normal processing

## Handler Implementation Notes

### include Handler

Implemented in `expr.go` as a private method (`newIncludeHandler()`) since it needs access to:
- Private `loadAndMergeFileWithContext()` for recursive file loading
- The Expr's filesystem (fs.FS)
- Context propagation through the evaluation pipeline

Automatically registered with highest priority (1000) when creating an Expr instance.
Can be overridden by passing a custom include handler via `WithDirectiveHandler()`.

### if & for Handlers

These handlers are tightly integrated with `expr.go` since they need access to:
- Private `processMapWithContext()` for recursive template processing
- Context propagation through the evaluation pipeline

Current implementation:
- Handler interface defined in `handlers/`
- Core logic remains in `expr.go`
- Extraction into handlers package is for semantic clarity

### interpolation Handler

String variable interpolation is applied automatically to all string values 
during document processing (in `expr.processWithContext()`). Two modes:

- **Lenient mode:** `InterpolateString()` from `util.go` - returns unchanged placeholders 
  for undefined variables (used in `if.go`)
- **Strict mode:** `InterpolateStringWithContext()` - errors on undefined variables 
  (used in document processing)

The handler provides `InterpolateValue()` for explicit control over interpolation 
with error handling based on context (path information).

### discard Handler

Simple standalone handler that checks a boolean value. Fully independent.

### matrix Handler

Complex handler that returns expanded job definitions. The actual template 
processing is done in `expr.go` using the returned job data.

Returns `[]any` of job configurations with matrix variables. The parent
`processMapWithContext()` handles template expansion for each job.

## Testing

Unit tests verify:
- Input validation
- Error handling
- Edge cases

Fixture tests verify:
- Full YAML processing
- Integration with other directives
- Real-world examples

Run tests:
```bash
go test ./handlers -v
go test ./... -run TestFixtures
```

## Integration with expr.go

Handlers are called from `processMapWithContext()`:

```go
for directive, handler := range e.config.handlers {
    if value, ok := m[directive]; ok {
        result, consumed, err := handler(ctx, m, value)
        // ... process result ...
    }
}
```

Order of operation:
1. Parse YAML map
2. Check custom handlers (in order of registration)
3. Check built-in `for` directive (loop expansion)
4. Check built-in `if` directive (conditional inclusion)
5. Process remaining keys normally

Note: The `include` handler is automatically registered as a custom handler with the highest priority,
so it runs before all other directives.

## Files

- `discard.go` - Discard handler (standalone)
- `if.go` - If handler with condition evaluation
- `for.go` - For handler helpers and parsing
- `interpolation.go` - Interpolation handler and utilities
- `matrix.go` - Matrix handler with cartesian product
- `util.go` - Shared utilities (lenient interpolation)
- `models.go` - Type aliases for handlers
- `*_test.go` - Unit tests
