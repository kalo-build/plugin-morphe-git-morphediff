# Morphe Git MorpheDiff Plugin

A Kalo plugin that generates semantic schema diffs between two versions of a Morphe registry, producing `KA:MD1:YAML1` format diff artifacts.

## Overview

This plugin compares two Morphe registry states (base and head) and generates semantic diff artifacts that describe:
- **Add** operations (new models, fields, relationships, etc.)
- **Remove** operations (deleted artifacts)
- **Modify** operations (changed types, attributes, cardinality)
- **Rename** operations (name changes with preserved structure)

## Purpose

The generated diff artifacts can be consumed by downstream plugins to:
- Generate SQL migration scripts (`plugin-morphediff-psql`)
- Update TypeScript type definitions incrementally
- Create API changelog documentation
- Validate breaking changes in CI/CD pipelines

## Usage

### With Kalo CLI

Configure in your `kalo.yaml`:

```yaml
stores:
  # Base version from git (e.g., main branch)
  KA_GIT_BASE:
    format: "KA:MO1:YAML1"
    type: "gitRepository"
    options:
      repoRoot: "."
      ref: "main"
      subPath: "morphe/registry"

  # Head version (current working directory)
  KA_MO_YAML:
    format: "KA:MO1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe/registry"

  # Output directory for diff files
  KA_MORPHE_DIFFS:
    format: "KA:MD1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe-diffs"

plugins:
  "@kalo-build/plugin-morphe-git-morphediff":
    version: "v1.0.0"
    inputs:
      base:
        format: "KA:MO1:YAML1"
        store: "KA_GIT_BASE"
      head:
        format: "KA:MO1:YAML1"
        store: "KA_MO_YAML"
    output:
      format: "KA:MD1:YAML1"
      store: "KA_MORPHE_DIFFS"

pipelines:
  diff:
    description: "Generate diff from git refs"
    stages:
      - name: "morphe-diff"
        steps:
          - "plugin: @kalo-build/plugin-morphe-git-morphediff"
```

Then run:

```bash
kalo run diff
```

## Output

Generates a `morphe-diff.yaml` file containing:
- **Metadata** - Version information, git provenance, timestamps, change summary
- **Changes** - Ordered list of delta operations with classifications

### Example Output

```yaml
metadata:
  spec_version: KA:MD1:YAML1
  source:
    ref: main
    commit: "5bb463a6fa647f06516858686495e78cbf570a45"
    timestamp: "2025-06-11T07:27:03Z"
  target:
    ref: head
    timestamp: "2026-01-27T10:30:00Z"
  summary:
    total_changes: 3
    breaking: 1
    additive: 2
    safe: 0
  generated_at: "2026-01-27T10:30:00Z"

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

### Git Provenance

The output includes git commit hashes and timestamps for reproducibility:

- **`ref`**: Git ref name (e.g., "main", "HEAD", branch name)
- **`commit`**: Full git commit hash (only for `gitRepository` stores)
- **`timestamp`**: Commit timestamp for git refs, or generation time for local stores

### Archive Mode

By default, diff files are ephemeral and should be gitignored. For historical tracking, enable archive mode in your pipeline config:

```yaml
pipelines:
  diff:
    stages:
      - name: "morphe-diff"
        steps:
          - "plugin: @kalo-build/plugin-morphe-git-morphediff"
        config:
          archiveDiffs: true  # Creates timestamped files: 20260127123456_morphe-diff.yaml
```

## Features

- **Models** - Detect model additions, removals, and modifications
- **Entities** - Track entity structure changes
- **Enums** - Identify enum entry additions/removals
- **Structures** - Compare structure field changes
- **Fields** - Detect field type changes, attribute modifications
- **Relationships** - Track relationship additions, removals, type changes
- **Change Classification** - Automatic breaking/additive/safe classification

## Building

```bash
# Build the WASM plugin
./scripts/build.sh

# Output: dist/morphe-git-morphediff-v1.0.0.wasm

# Windows
scripts\build.bat
```

## Input Stores

The plugin requires two Morphe registry stores:

| Input | Description |
|-------|-------------|
| `base` | The original/previous state (typically from git) |
| `head` | The new/current state (typically local filesystem) |

Each store should contain:
- Models directory (`.mod` files)
- Entities directory (`.ent` files)
- Enums directory (`.enum` files)
- Structures directory (`.str` files)

## License

MIT License
