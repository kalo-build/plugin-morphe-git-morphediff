# Morphe Diff Generator Plugin

A Morphe compilation plugin that generates semantic schema diffs between two versions of a Morphe registry, producing `KA:MD1:YAML1` format diff artifacts.

## Overview

This plugin compares two Morphe registry states (base and head) and generates semantic diff artifacts that describe:
- **Add** operations (new models, fields, relationships, etc.)
- **Remove** operations (deleted artifacts)
- **Modify** operations (changed types, attributes, cardinality)
- **Rename** operations (name changes with preserved structure)

## Purpose

The generated diff artifacts can be consumed by downstream plugins to:
- Generate SQL migration scripts
- Update TypeScript type definitions incrementally
- Create API changelog documentation
- Validate breaking changes in CI/CD pipelines

## Input

The plugin requires two Morphe registries:

1. **Base Registry** - The original/previous state
   - Models directory (`.mod` files)
   - Entities directory (`.ent` files)
   - Enums directory (`.enum` files)
   - Structures directory (`.str` files)

2. **Head Registry** - The new/current state
   - Models directory (`.mod` files)
   - Entities directory (`.ent` files)
   - Enums directory (`.enum` files)
   - Structures directory (`.str` files)

## Output

Generates a single `morphe-diff.yaml` file containing:
- **Metadata** - Version information, timestamps, change summary
- **Changes** - Ordered list of delta operations with classifications

Output conforms to `KA:MD1:YAML1` specification.

## Example Output

```yaml
metadata:
  spec_version: KA:MD1:YAML1
  source:
    version: base
    timestamp: "2025-01-15T10:30:00Z"
  target:
    version: head
    timestamp: "2025-01-20T14:45:00Z"
  summary:
    total_changes: 3
    breaking: 1
    additive: 2
    safe: 0
  generated_at: "2025-01-28T15:00:00Z"

changes:
  - operation: add
    type: field
    target:
      model: User
      field: PhoneNumber
    definition:
      type: String
      attributes:
        - nullable
    classification: additive
```

## Features

- ✅ **Models** - Detect model additions, removals, and modifications
- ✅ **Entities** - Track entity structure changes
- ✅ **Enums** - Identify enum entry additions/removals
- ✅ **Structures** - Compare structure field changes
- ✅ **Fields** - Detect field type changes, attribute modifications
- ✅ **Relationships** - Track relationship additions, removals, type changes
- ✅ **Change Classification** - Automatic breaking/additive/safe classification

## Usage

### With Git (Recommended)

Use the helper script to compare git refs:

```bash
# Compare current changes against main
./scripts/diff-with-git.sh main HEAD

# Compare two releases
./scripts/diff-with-git.sh v1.0.0 v2.0.0

# Windows
scripts\diff-with-git.bat main HEAD
```

**Output**: `morphe-diff.yaml`

See [GIT_WORKFLOW_GUIDE.md](GIT_WORKFLOW_GUIDE.md) for detailed git integration strategies.

### With Directories

For testing or manual workflows:

```bash
kalo compile --plugin @kalo-build/plugin-morphe-git-morphediff \
  --config '{"baseInputPath":"./base","headInputPath":"./head","outputPath":"./diff.yaml"}'
```

## Building

```bash
# Build the WASM plugin
./scripts/build.sh

# Output: dist/morphe-git-morphediff-v1.0.0.wasm
```

## Documentation

- [GETTING_STARTED.md](GETTING_STARTED.md) - Quick start guide with examples
- [GIT_WORKFLOW_GUIDE.md](GIT_WORKFLOW_GUIDE.md) - Git integration strategies (script vs kalo kx)
- [REQUIREMENTS.md](REQUIREMENTS.md) - Technical specifications
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Command reference
- [SAMPLE_OUTPUT.md](SAMPLE_OUTPUT.md) - Example diff artifacts
- [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - What was built

## License

MIT License
