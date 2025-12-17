# yamlexpr Features Documentation

This directory contains detailed documentation for each yamlexpr feature. Each feature has two components:

1. **`feature.md.sh`** - A shell script that generates markdown documentation
2. **Example YAML files** - Input examples with expected output

## Features

### Core Features

- **[Interpolation](interpolation.md.sh)** - Variable substitution with `${variable}` syntax
- **[Conditionals](conditionals.md.sh)** - Include/omit blocks with `if:` directive
- **[For Loops](for-loops.md.sh)** - Iterate and expand arrays with `for:` directive
- **[Matrix](matrix.md.sh)** - Generate cartesian products with `matrix:` directive
- **[Include](include.md.sh)** - Compose external YAML files with `include:`
- **[Document Expansion](document-expansion.md.sh)** - Root-level `for:` and `matrix:` for multi-document output

## Documentation Structure

Each feature documentation includes:

1. **Syntax Cheat Sheet** - Quick reference for syntax patterns
2. **Description** - Feature purpose and use cases
3. **Core Concepts** - Key ideas and behavior
4. **Examples** - Input/output pairs demonstrating the feature
5. **Common Use Cases** - When and why to use each feature
6. **Edge Cases** - Special behaviors and limitations

## Generating Documentation

To generate a single feature's markdown:

```bash
./docs/features/interpolation.md.sh > interpolation.md
./docs/features/conditionals.md.sh > conditionals.md
./docs/features/for-loops.md.sh > for-loops.md
./docs/features/matrix.md.sh > matrix.md
./docs/features/include.md.sh > include.md
./docs/features/document-expansion.md.sh > document-expansion.md
```

To generate all features:

```bash
cd docs/features
for script in *.md.sh; do
  feature="${script%.md.sh}"
  ./"$script" > "$feature.md"
done
```

## Component Examples

The `components/` directory contains reusable YAML components:

- **`_service-base.yaml`** - Common service configuration
- **`_database-config.yaml`** - Database settings template
- **`example-compose.yaml`** - Example showing composition with includes and for loops

These are used in documentation examples and can serve as templates for your own configurations.

## Integration with Main Docs

These feature documents should be referenced from:
- **docs/syntax.md** - Comprehensive syntax reference
- **docs/tutorial.md** - Practical tutorials
- **README.md** - Feature checklist

## Adding New Features

To add documentation for a new feature:

1. Create `docs/features/feature-name.md.sh` with:
   - Syntax cheat sheet
   - Description
   - Core concepts
   - Examples with input/output
   - Use cases
   - Edge cases

2. Add the feature to this README

3. Update main documentation to reference the new feature

## Example Format

Each script should follow this structure:

```bash
#!/bin/bash
# Generates [feature name] feature documentation
# Usage: ./[feature].md.sh > [feature].md

cat << 'EOF'
# [Feature Name]

## Syntax Cheat Sheet

```yaml
# Examples of syntax
```

## Description

...

## Core Concepts

...

## Examples

### [Example Title]

**Input:**

```yaml
...
```

**Output:**

```yaml
...
```

EOF

```

EOF
```
