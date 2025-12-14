# Agent Guidelines for yamlexpr

This file contains conventions and preferences for AI agents working on this codebase.

## Package Architecture

### Package Organization

- **stack/**: Variable scope stack management (reusable, shared pattern with vuego/lessgo)
- **model/**: Core interfaces and types (DirectiveHandler, Processor, Context)
- **handlers/**: Modular directive handlers (if, for, embed, interpolation, discard)
- **runner/**: Deterministic matrix expansion at document root level
- **merge/**: YAML merging utilities supporting document composition
- **testdata/fixtures/**: Black-box test fixtures (001-099: basic, 040-099: for/if, 100-199: integration, 200-299: matrix)

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

## Documentation Standards

### Handler Documentation

Each handler must have a corresponding markdown file in `docs/handlers/`:

**File**: `docs/handlers/{handler_name}.md`

**Structure**:
1. **Overview** - What the handler does
2. **Syntax** - How to use it in YAML
3. **Features** - Key capabilities
4. **Examples** - Real-world usage patterns
5. **API Functions** - Public functions with signatures and descriptions
6. **Error Handling** - Common error cases and messages
7. **See Also** - Links to related handlers

**Examples section must include**:
- Basic usage
- With variables/interpolation
- Combined with other directives (for, if)
- Error cases (when appropriate)

**Key points to document**:
- Single-line expression vs multi-line (e.g., for expressions are single-line)
- Interaction with variable scope
- Handler priority
- Type requirements (e.g., matrix dimensions must be arrays)
- Null/undefined handling behavior

### Documentation Index

Maintain indices in:
- **`docs/README.md`** - Complete documentation index with handler links
- **`README.md`** (main) - Quick reference links to handler docs and docs index
- **Handler Priority Table** - Show execution order and impact

### Keeping Documentation Accurate

- **Fixtures are source of truth**: All behavior documented in `testdata/fixtures/*.yaml`
- **Update docs when adding features**: Don't let code and docs diverge
- **Link between related handlers**: Cross-reference handlers that work together
- **Show output examples**: Include both input and output YAML for clarity

## Key Principles

- **Fixtures are source of truth**: All behavior documented in testdata/fixtures/*.yaml
- **Stack is reusable**: Keep it minimal and importable by other packages
- **ExprContext propagates state**: Pass it through all processing functions for consistent path tracking
- **Expr is self-contained**: Handle composition, interpolation, conditionals here
- **Black box testing**: Only test exported APIs, not implementation details
- **Follow Go conventions**: Use standard library, avoid unnecessary dependencies
- **Document comprehensively**: Handler docs are as important as code quality
