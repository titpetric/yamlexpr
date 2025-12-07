# Development Workflow for yamlexpr

This document explains the SDLC (Software Development Life Cycle) workflow for yamlexpr features.

## Feature Status States

Each feature transitions through these states:

1. **waiting** - Feature is planned but not started
2. **doing** - Feature is under active development
3. **testing** - Implementation complete, all tests passing
4. **iterating** - Testing revealed issues, iterating fixes
5. **done** - Complete and ready for production

## Before Starting Work

1. Check `README.md` for feature status
2. If starting a new feature, update status: `waiting` → `doing`
3. Create a new todo item in this session if needed

Example:

```markdown
- [ ] **Feature 1 (doing)**: Include composition
  - Placeholder items to complete
```

## During Development

### Phase 1: Planning & Test Fixtures
- Create test fixtures in `testdata/fixtures/NNN-description.yaml`
- Use `---` delimiter to separate input from expected output
- Follow lessgo/vuego fixture pattern
- Example: `020-include-single-file.yaml`

### Phase 2: Write Black Box Tests
- Create `xxx_test.go` in the same package as implementation
- Use `package xxx_test` (black box testing)
- Follow naming: `TestXxx_Feature` or `TestFixtures`
- Use `github.com/stretchr/testify/require` assertions
- Example: `expr_fixtures_test.go`

### Phase 3: Implementation
- Add code to appropriate package (stack/ or expr/)
- Update stubs in expr.go if needed (e.g., parseYAML, evaluateCondition)
- Keep imports minimal - no external dependencies for stack/
- Follow established code style (see AGENTS.md)

### Phase 4: Documentation
- Update DESIGN.md with approach and decisions
- Add examples to README.md if user-facing
- Update function godoc comments
- Document any new directives or syntax

## After Implementation - SDLC Checklist

Before marking as `done`, verify:

### Testing ✓
- [ ] All new tests passing
- [ ] No regression in existing tests
- [ ] Edge cases covered (empty, null, missing values)
- [ ] Error cases tested and documented
- [ ] Run: `go test -v ./...`

### Code Quality ✓
- [ ] Code style consistent (see AGENTS.md)
- [ ] No unused imports
- [ ] No external dependencies for stack/
- [ ] Godoc comments on exported items
- [ ] Run: `go build ./...`

### Security ✓
- [ ] No panic on invalid input
- [ ] Proper error handling and messages
- [ ] Input validation for file paths (relative path attacks)
- [ ] Stack limits respected (pathCacheLimit = 256)

### Consistency ✓
- [ ] Matches patterns from vuego/lessgo
- [ ] Variable naming is clear
- [ ] Error messages start with lowercase
- [ ] Comments explain "why" not "what"

### Documentation ✓
- [ ] README.md updated with feature status
- [ ] Examples added for user-facing features
- [ ] DESIGN.md updated with approach
- [ ] Inline comments for complex logic
- [ ] AGENTS.md section updated if conventions changed

### Final Step

When all above complete:

```markdown
- [x] **Feature N (done)**: Description
  - [x] Implementation complete
  - [x] All tests passing
  - [x] Documentation complete
  - [x] Code reviewed
  - [x] Security verified
```

## Example: Implementing Feature 1 (Include Composition)

### Step 1: Update README (waiting → doing)

```bash
# Edit README.md
- [ ] **Feature 1 (doing)**: Include composition
```

### Step 2: Create Fixtures

```bash
# testdata/fixtures/020-include-single-file.yaml
config:
  include: "other.yaml"
  name: "main"
---
config:
  name: "main"
  included_key: "included_value"

# testdata/fixtures/021-include-list.yaml
files:
  include:
    - "file1.yaml"
    - "file2.yaml"
---
files:
  result1: "value1"
  result2: "value2"
```

### Step 3: Add Tests

```go
// expr/expr_fixtures_test.go
func TestFixtures(t *testing.T) {
	fixtures := []string{
		"020-include-single-file",
		"021-include-list",
	}
	// Load and test each fixture
}
```

### Step 4: Implement Feature

```go
// expr/expr.go - update stubs
func (e *Expr) handleInclude(incl any, result map[string]any, st *stack.Stack) error {
	// Implement include logic
}

func (e *Expr) loadAndMergeFile(filename string, result map[string]any, st *stack.Stack) error {
	// Implement file loading
}
```

### Step 5: Run Tests

```bash
go test -v ./expr
```

### Step 6: Update Documentation
- Add examples to README.md
- Update DESIGN.md with implementation details
- Add security note about path validation

### Step 7: Final Checklist
- Testing: ✓ All tests passing
- Code Quality: ✓ No linter warnings
- Security: ✓ Path validation in place
- Consistency: ✓ Matches existing patterns
- Documentation: ✓ README, DESIGN.md, comments updated

### Step 8: Mark Done

```markdown
- [x] **Feature 1 (done)**: Include composition
  - [x] File loading from fs.FS
  - [x] YAML merging
  - [x] Relative path resolution
  - [x] 2 test fixtures passing
  - [x] Documentation complete
```

## Commands for Development

```bash
# Build all packages
go build ./...

# Run all tests
go test -v ./...

# Run specific package
go test -v ./stack
go test -v ./expr

# Run specific test
go test -v -run TestStack_Resolve

# Run fixture tests only
go test -v -run TestFixtures

# Check for unused imports
go mod tidy

# Build and test in one shot
go build ./... && go test -v ./...
```

## File Organization

```
yamlexpr/
├── stack/                    # Core variable scoping (reusable)
│   ├── stack.go             # Implementation
│   └── stack_test.go        # Black box tests
│
├── expr/                     # YAML expression evaluation
│   ├── expr.go              # Main Expr type
│   ├── expr_test.go         # Basic tests
│   ├── interpolate.go       # String interpolation
│   ├── interpolate_test.go  # Interpolation tests
│   └── expr_fixtures_test.go # Fixture-based tests (planned)
│
├── testdata/
│   └── fixtures/            # Test fixtures
│       ├── 001-*.yaml       # Basic pass-through
│       ├── 030-*.yaml       # Conditionals (if directives)
│       ├── 040-*.yaml       # For loops
│       ├── 050-*.yaml       # Combined for/if/interpolation
│       ├── 060-*.yaml       # Include composition
│       ├── 070-*.yaml       # Nested structures
│       ├── 080-*.yaml       # Advanced features
│       ├── 090-*.yaml       # Unquoted syntax variants
│       └── _*.yaml          # Base files for includes
│
├── README.md                # Feature status and usage
├── AGENTS.md                # Development conventions
├── DESIGN.md                # Architecture decisions
├── DEVELOPMENT.md           # This file - SDLC workflow
├── go.mod                   # Dependencies
└── LICENSE
```

## Notes

- Stack package has NO external dependencies (only stdlib)
- Expr package depends only on stack/ and stdlib
- All tests use testify/require for assertions
- All tests use black box approach (xxx_test.go packages)
- Fixtures are ground truth - expected output is source of truth
- Every feature needs tests before it's marked "done"
- Security review is part of SDLC, not optional
