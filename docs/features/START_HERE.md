# yamlexpr Features - Start Here

Welcome to yamlexpr feature documentation! This guide will help you find what you need.

## 🚀 Quick Start

**New to yamlexpr?** Start with:
1. [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - One-page syntax cheat sheet
2. [docs/tutorial.md](../tutorial.md) - Comprehensive tutorial

**Need help with a specific feature?** Jump to the section below.

## 📚 Feature Documentation

Each feature has detailed documentation with syntax, examples, and use cases.

### 1. Interpolation

**Variable substitution with `${variable}` syntax**

- 📖 [Read interpolation.md](interpolation.md)
- 🎯 Use case: Reference variables in YAML strings
- ⏱️ Read time: 5 minutes

**Syntax:**

```yaml
name: "Hello, ${user_name}!"
connection: "postgres://${db.host}:${db.port}/mydb"
```

### 2. Conditionals (`if:`)

**Include/exclude blocks based on conditions**

- 📖 [Read conditionals.md](conditionals.md)
- 🎯 Use cases: Feature flags, environment-specific config
- ⏱️ Read time: 10 minutes

**Syntax:**

```yaml
database:
  if: ${enable_database}
  host: "localhost"
```

### 3. For Loops (`for:`)

**Iterate arrays and create multiple items**

- 📖 [Read for-loops.md](for-loops.md)
- 🎯 Use cases: Service enumeration, platform builds
- ⏱️ Read time: 15 minutes

**Syntax:**

```yaml
servers:
  - for: server in server_list
    name: "${server}"
```

### 4. Matrix (`matrix:`)

**Generate cartesian products of dimensions**

- 📖 [Read matrix.md](matrix.md)
- 🎯 Use cases: CI/CD test matrices, cross-platform builds
- ⏱️ Read time: 12 minutes

**Syntax:**

```yaml
matrix:
  os: [linux, windows]
  arch: [x86_64, arm64]
  exclude:
    - os: windows
      arch: arm64
job: "${os}/${arch}"
```

### 5. Include (`include:`)

**Compose external YAML files**

- 📖 [Read include.md](include.md)
- 🎯 Use cases: Reusable configs, component libraries
- ⏱️ Read time: 12 minutes

**Syntax:**

```yaml
include: "_base.yaml"
database:
  include: "_db-config.yaml"
```

### 6. Document Expansion

**Root-level directives creating multiple documents**

- 📖 [Read document-expansion.md](document-expansion.md)
- 🎯 Use cases: Multi-environment generation, test matrices
- ⏱️ Read time: 10 minutes

**Syntax:**

```yaml
for: env in [staging, production]
environment: "${env}"
```

## 🔗 Navigation

### By Use Case

**Building reusable configurations?** → [Include](include.md) + [components/](components/)

**Creating multi-environment configs?** → [Document Expansion](document-expansion.md) + [For Loops](for-loops.md)

**Setting up CI/CD matrices?** → [Matrix](matrix.md)

**Filtering configurations?** → [Conditionals](conditionals.md) + [For Loops](for-loops.md)

**Quick syntax lookup?** → [QUICK_REFERENCE.md](QUICK_REFERENCE.md)

### By Learning Style

**Visual learner?** → Start with [QUICK_REFERENCE.md](QUICK_REFERENCE.md) for syntax examples

**Learn by doing?** → Check [components/](components/) for working examples

**Comprehensive learner?** → Read full feature documentation in order

**API user?** → See [../api.md](../api.md) for Go usage

## 📂 Component Examples

Reusable YAML components demonstrating best practices:

- **[_service-base.yaml](components/_service-base.yaml)** - Common service defaults
- **[_database-config.yaml](components/_database-config.yaml)** - Database template
- **[example-compose.yaml](components/example-compose.yaml)** - Full composition example

## 🔍 Quick Reference

**Syntax Cheat Sheet:** → [QUICK_REFERENCE.md](QUICK_REFERENCE.md)

**Common Patterns:**
- Reusable defaults with includes
- Layered configuration
- Environment-specific services
- Feature flags and conditionals

## 📋 Feature Comparison

| Feature                | Use                   | Output              |
|------------------------|-----------------------|---------------------|
| **Interpolation**      | Reference variables   | String substitution |
| **Conditionals**       | Filter by condition   | Omit keys           |
| **For Loops**          | Iterate arrays        | Multiple items      |
| **Matrix**             | Generate combinations | Multiple items      |
| **Include**            | Load files            | Merged YAML         |
| **Document Expansion** | Root-level for/matrix | Multiple documents  |

## 🎓 Learning Path

### Beginner
1. [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - 5 min overview
2. [Interpolation](interpolation.md) - Variable substitution
3. [Conditionals](conditionals.md) - If directives

### Intermediate
4. [For Loops](for-loops.md) - Iterations
5. [Include](include.md) - Composition

### Advanced
6. [Matrix](matrix.md) - Complex patterns
7. [Document Expansion](document-expansion.md) - Multi-document
8. [components/](components/) - Real examples

## 📖 Related Documentation

- **[Tutorial](../tutorial.md)** - Step-by-step guide
- **[Syntax Reference](../syntax.md)** - Technical reference
- **[API Reference](../api.md)** - Go API documentation
- **[Custom Syntax](../custom-syntax.md)** - Customize directives

## 🤔 FAQ

**Q: What's the difference between for and matrix?** A: `for:` iterates a single array linearly. `matrix:` generates a cartesian product of multiple dimensions.

**Q: Can I combine features?** A: Yes! You can use `for:` + `if:`, `include:` + `for:`, etc.

**Q: How are variables available?** A: All root-level keys become variables available in `${}` and expressions.

**Q: Can includes be conditional?** A: Yes, wrap the `include:` in an `if:` block.

**Q: How do I reuse configurations?** A: Use `include:` with helper files in `components/` directory.

## 💡 Tips

- **Use quotes** for `if:` and `for:` directives: `for: "item in items"`
- **Nested access**: Use dot notation: `${config.server.host}`
- **Filter before expanding**: Combine `for:` with `if:` for filtering
- **Layered configs**: Use multiple `include:` statements for inheritance
- **Index access**: Use `(idx, item)` syntax for both index and value

## 🔗 Links

- [Root README](../../README.md) - Main project documentation
- [Feature Directory README](README.md) - This directory overview
- [Documentation Updates](../../DOCUMENTATION_UPDATES.md) - Implementation details
- [GitHub Repository](https://github.com/titpetric/yamlexpr) - Source code

---

**Last Updated**: December 17, 2025 **Total Documentation**: 1,836 lines across 6 feature docs **Examples**: 40+ working examples with input/output pairs
