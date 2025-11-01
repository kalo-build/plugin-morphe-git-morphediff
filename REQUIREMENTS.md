# Morphe Git Diff Plugin - Implementation Requirements

## Overview

This plugin generates semantic schema diffs between two Morphe registry states, producing `KA:MD1:YAML1` format diff artifacts.

## Core Functionality

### Input Requirements

The plugin requires two complete Morphe registry directories:

1. **Base Registry** - The original/source schema state
   - `models/` directory containing `.mod` files
   - `entities/` directory containing `.ent` files
   - `enums/` directory containing `.enum` files
   - `structures/` directory containing `.str` files

2. **Head Registry** - The new/target schema state
   - Same directory structure as base

### Output

Produces a single YAML file containing:
- **Metadata**: Version info, timestamps, change summary
- **Changes**: List of delta operations with classifications

## Supported Operations

### 1. Add
Detects new artifacts in the head registry:
- New models
- New fields in existing models/entities/structures
- New relationships
- New enums
- New enum entries

### 2. Remove
Detects deleted artifacts:
- Removed models, entities, structures, enums
- Removed fields
- Removed relationships
- Removed enum entries

**Classification**: Always `breaking`

### 3. Modify
Detects changes to existing artifacts:
- Field type changes
- Field attribute changes
- Relationship type changes (cardinality, polymorphism)
- Enum entry value changes
- Enum type changes

**Classification**: Usually `breaking`

## Change Classification

### Breaking Changes
- Remove any artifact
- Modify field types incompatibly
- Change nullable → mandatory
- Remove enum entries
- Modify relationship cardinality

### Additive Changes
- Add new models, entities, structures, enums
- Add new fields (nullable/optional)
- Add new relationships
- Add new enum entries

### Safe Changes
Currently not implemented (future: renames with fingerprints)

## Implementation Details

### Comparison Algorithm

1. Load both registries using `morphe-go` loader
2. Compare each artifact type:
   - **Enums**: Compare entries, type
   - **Structures**: Compare fields, attributes
   - **Models**: Compare fields, identifiers, relationships
   - **Entities**: Compare fields, identifiers, relationships
3. Generate delta operations for all differences
4. Classify each change as breaking/additive/safe
5. Output to YAML file

### Test Coverage

Unit tests cover:
- Added/removed/modified artifacts
- Edge cases (empty registries, no changes)
- Field type changes
- Attribute changes
- Relationship changes
- Polymorphic relationships
- Enum entry modifications

Integration tests verify:
- Complete diff generation workflow
- YAML output format correctness
- Metadata generation
- Change classification accuracy

## Dependencies

- `github.com/kalo-build/morphe-go` - Morphe schema loader and types
- `gopkg.in/yaml.v3` - YAML marshaling
- `github.com/stretchr/testify` - Testing framework (tests only)

## Future Enhancements

1. **Rename Detection**: Use structural fingerprints to detect renames vs remove+add
2. **Move Detection**: Detect field movements between models
3. **Deprecation Tracking**: Support deprecation metadata
4. **Custom Classification Rules**: Configurable breaking change detection
5. **Diff Optimization**: Combine related changes into single operations
