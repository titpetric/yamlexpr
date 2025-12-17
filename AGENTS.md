# Agent Guidelines for yamlexpr

This file contains conventions and preferences for AI agents working on this codebase.

## Package Architecture

### Package Organization

- **stack/**: Variable scope stack management (shared with vuego/lessgo patterns)
  - `stack.go`: Stack implementation for variable lookup and resolution
  - `stack_test.go`: Black box tests for Stack API

- **expr/**: YAML expression evaluation and composition
  - `expr.go`: Main Expr type for evaluating YAML documents with stack values
  - `expr_test.go`: Black box tests for Expr API
  - `context.go`: ExprContext for carrying evaluation state through processing
  - `context_test.go`: Tests for ExprContext API
  - `processor.go`: Core document processing (for, if, include handlers)
  - `processor_test.go`: Integration tests for document transformation
  - `interpolate.go`: ${} syntax interpolation in string values
  - `interpolate_test.go`: Interpolation tests

- **testdata/fixtures/**: Test fixtures matching lessgo pattern
  - `NNN-description.yaml`: Input YAML source
  - `NNN-description.yaml.expected`: Expected output YAML after processing

### Design Principles

- **Minimize external dependencies**: Only import go-expr when evaluating `if` conditions
- **Follow vuego/lessgo patterns**: Use black box testing, fixture-based testing, same code style
- **Separation of concerns**:
  - `stack` handles variable lookup (reusable across projects)
  - `expr` handles YAML composition (for, if, include, interpolation)
  - `ExprContext` carries evaluation state (stack, path, include chain) through processing
- **Process YAML as a data structure**: Parse once, transform in-place, serialize once
- **Fixtures are ground truth**: YAML.expected files are the source of truth for behavior
- **Context propagation**: Pass `ExprContext` through all processing functions for consistent path tracking and error reporting

## Feature Implementation (Feature 1: Include)

Feature 1 (composition with `include:`) is implemented in `expr` package:

```yaml
# In YAML, include directive pulls in another file
config:
  include: "other.yaml"
  
# Or as a list to include multiple files
includes:
  - "file1.yaml"
  - "file2.yaml"
```

Files are resolved relative to the base directory (fs.FS) passed to Expr.New().

## Testing Conventions

### Black Box Testing Philosophy
- **Prefer black box tests** using `package yamlexpr_test` or specific `package expr_test` instead of internal packages
- Tests should only interact with exported APIs
- This allows running individual test files: `go test -v stack_test.go`

### Test File Organization
- **Each major `.go` file must have a corresponding `_test.go` file** in the same package
- Example: `expr.go` → `expr_test.go`, `processor.go` → `processor_test.go`
- Group tests for related functionality in the same test file as the implementation
- **Fixture-based tests** should be in dedicated test files:
  - Example: `expr_fixtures_test.go` for processing all fixtures in `testdata/fixtures/`
  - These load `.yaml` and `.expected` file pairs and assert output matches

### Test Naming Convention

Follow the pattern `Test[Receiver_]Function`:
- Methods: `TestExpr_Evaluate`, `TestStack_Push`
- Functions: `TestNewExpr`, `TestInterpolate`
- Fixture tests: `TestFixtures`, `TestFixtures_Include`, `TestFixtures_Conditionals`

### Assertion Style

```go
// Good: Use testify/require
require.NoError(t, err)
require.Equal(t, expected, actual)
require.Contains(t, str, substring)

// Avoid: t.Fatal, t.Error, t.Errorf
if err != nil {
	t.Fatal(err) // Don't do this
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test file
go test -v expr_test.go

# Run specific test function
go test -v -run TestExpr_Evaluate

# Run fixture-based tests only
go test -v -run TestFixtures
```

## Development Commands

### Building and Testing

```bash
# Build the package
go build ./...

# Run all tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Test a specific package
go test -v ./stack ./expr
```

### Fixture Testing

Test fixtures follow the lessgo pattern with `---` delimiter:

**File: `testdata/fixtures/001-simple.yaml`**

```yaml
# Input YAML (before processing)
items:
  - name: "item1"
    count: 1
  - name: "item2"
    count: 2
---
# Expected YAML (after processing)
items:
  - name: "item1"
    count: 1
  - name: "item2"
    count: 2
```

**Fixture test runner**:

```bash
# All fixtures are tested via TestFixtures function
go test -v -run TestFixtures ./expr
```

### Generating Documentation

Feature documentation is generated from shell script templates. All `.md` files in `docs/features/` are generated from corresponding `.md.sh` scripts.

```bash
# Generate all feature documentation
task generate

# Or manually generate a single feature
cd docs/features
bash interpolation.md.sh > interpolation.md
```

**Important**: Always edit the `.md.sh` scripts, not the generated `.md` files. The `.md` files are regenerated on each `task generate` run.

## Documentation Standards

### Code Block Requirements

All code blocks MUST have opening language tags (no bare ````` fences):

- ````yaml` - YAML configuration examples and syntax patterns
- ````go` - Go code examples (API usage, type definitions)
- ````text` - File structures, plain text, directory trees
- ````bash` - Shell commands and scripts
- `python`, `json`, etc. for other languages as needed

**Rule**: Every opening code fence must have a language tag. Closing fences are bare `````.

### Feature Documentation Structure

Each feature doc in `docs/features/` should follow this structure:

```markdown
# Feature Name

## Syntax Cheat Sheet

```yaml
# Show common syntax patterns
example: value
```

## Description

Brief explanation of what the feature does and why you'd use it.

## Core Concepts

- **Concept 1**: Explanation
- **Concept 2**: Explanation
- **Concept 3**: Explanation

## Examples

### Example Title

**Input:**

```yaml
# Show example input
key: value
```

**Output:**

```yaml
# Show expected output
key: processed_value
```

## Common Use Cases

- Use case 1: Description
- Use case 2: Description

## Edge Cases

### Edge Case Name

Description and example.

```

### Markdown File Naming

**Generated files** (`docs/features/`):
- `feature-name.md.sh` - Generator script (edit this)
- `feature-name.md` - Generated output (regenerated by `task generate`)

**Manual files**:
- `docs/syntax.md` - Complete syntax reference
- `docs/api.md` - API documentation
- `docs/tutorial.md` - Practical tutorials
- `docs/DESIGN.md` - Architecture and design decisions

### Code Block Best Practices

**Input/Output Examples**:
```markdown
**Input:**

```yaml
config:
  name: "example"
```

**Output:**

```yaml
config:
  name: "example"
  processed: true
```

```

**File Structures**:
```markdown
```text
config/
├── base/
│   ├── _database.yaml
│   └── _cache.yaml
├── services/
│   ├── _api-service.yaml
│   └── _worker-service.yaml
└── app.yaml
```

```

**API Examples**:
```markdown
```go
expr := yamlexpr.New(os.DirFS("."))
docs, err := expr.Load("config.yaml")
if err != nil {
    log.Fatal(err)
}
```

```

### Template Files

When generating documentation from scripts (`.md.sh`), ensure:

1. **Wrap dynamic output in fences**: 
   - Print opening fence: ` cat << 'EOF'\n```yaml\nEOF `
   - Include command that produces output
   - Print closing fence: ` cat << 'EOF'\n```\nEOF `

2. **Match fence pairs**: Every ` ```language ` opening must have corresponding ` ``` ` closing

3. **Validate syntax**: Run `bash -n script.md.sh` before using

4. **Test generation**: 
   ```bash
   bash docs/features/feature.md.sh | grep -c '^```'  # Should be even number
```

### Documentation Regeneration Workflow

1. **Edit the template**: Modify `docs/features/feature.md.sh`
2. **Verify syntax**: `bash -n docs/features/feature.md.sh`
3. **Test locally**: `bash docs/features/feature.md.sh > /tmp/test.md`
4. **Regenerate all**: `task generate`
5. **Verify output**: Check generated `.md` files for proper formatting

### Documentation Validation

Run before committing:

```bash
# Generate all docs
task generate

# Check for unmatched code blocks
for file in docs/features/*.md; do
  total=$(grep -c '^```' "$file" 2>/dev/null || echo "0")
  if [ $((total % 2)) -ne 0 ]; then
    echo "ERROR: $file has unmatched code blocks"
  fi
done

# Verify all opening fences have language tags
for file in docs/features/*.md; do
  opens=$(grep -c '^```' "$file")
  withLang=$(grep '^```[a-z]' "$file" | wc -l)
  if [ "$withLang" -lt "$opens" ]; then
    echo "WARNING: $file might have bare opening fences"
  fi
done
```

### Common Mistakes to Avoid

- ❌ Editing `.md` files directly (they'll be overwritten)
- ❌ Using bare ````` for opening fences (use `yaml`, `text`, etc.)
- ❌ Forgetting to run `task generate` after modifying `.md.sh` files
- ❌ Mixing code block styles (use consistent formatting across docs)
- ❌ Not testing generated output before committing

### Cross-Referencing

Link between documentation files using relative paths:

```markdown
See [Syntax Reference](syntax.md) for complete directive documentation.
See [For Loops](features/for-loops.md) for iteration patterns.
```

Files in `docs/features/` should reference main docs:

```markdown
For more details, see [Syntax Reference](../syntax.md).
```

## Code Style

### Godoc Comments
- **Always** add godoc comments for exported types, functions, and methods
- Start with the name of the item being documented
- Be concise but descriptive

**Examples:**

```go
// Stack provides stack-based variable lookup for templates.
type Stack struct { ... }

// NewStack constructs a Stack with optional initial root map.
func NewStack(root map[string]any) *Stack { ... }

// Resolve returns the value at the given path expression (e.g., "user.name").
func (s *Stack) Resolve(expr string) (any, bool) { ... }
```

### Error Messages
- Use lowercase for error messages (e.g., `"error loading %s"`)
- Use `fmt.Errorf` with `%w` to wrap errors
- Include context: filenames, paths, variable names
- Make errors actionable

**Examples:**

```go
return fmt.Errorf("error including %s: %w", filename, err)
return fmt.Errorf("undefined variable '%s' in if condition", varName)
```

### File Organization
- Group related functionality in separate files (e.g., `processor_*.go` for processing logic)
- Keep type definitions and constructors together
- Place helper functions near their usage

## Feature Checklist (In README.md)

Use this pattern in README.md for feature tracking:

```markdown
## Features

- [ ] **Feature 1 (waiting)**: Include composition - load external YAML files
  - Include directive resolves files from base directory
  - Merge external YAML into current structure
  
- [x] **Feature 2 (done)**: For loops - iterate over arrays/objects
  - `for:` syntax expands arrays with loop variable in scope
```

Each status is one of: `waiting`, `doing`, `testing`, `iterating`, `done`

## ExprContext Usage

The `ExprContext` type encapsulates evaluation state and replaces direct stack/path parameter passing:

```go
// Create context for top-level processing
ctx := expr.NewExprContext(&expr.ExprContextOptions{
	Stack:        myStack,
	Path:         "",
	IncludeChain: []string{},
})

// Navigate deeper into document
childCtx := ctx.AppendPath("database")
deepCtx := childCtx.AppendPath("config")

// Track include chain
includedCtx := ctx.WithInclude("config.yaml")

// Manage variable scope (for loops)
ctx.PushStackScope(map[string]any{"item": currentItem})
// ... process item ...
ctx.PopStackScope()
```

Methods on ExprContext:
- `Stack()`: Returns the variable resolution stack
- `Path()`: Returns current path in document
- `WithPath(newPath)`: Creates new context with updated path
- `AppendPath(segment)`: Appends to path (handles both keys and array indices)
- `WithInclude(filename)`: Creates new context with extended include chain
- `FormatIncludeChain()`: Formats chain for error messages
- `PushStackScope(map)`: Push variable scope
- `PopStackScope()`: Pop variable scope

## Key Principles

- **Fixtures are source of truth**: All behavior documented in testdata/fixtures/*.yaml
- **Stack is reusable**: Keep it minimal and importable by other packages
- **ExprContext propagates state**: Pass it through all processing functions for consistent path tracking
- **Expr is self-contained**: Handle composition, interpolation, conditionals here
- **Black box testing**: Only test exported APIs, not implementation details
- **Follow Go conventions**: Use standard library, avoid unnecessary dependencies
