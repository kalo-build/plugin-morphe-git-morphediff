# Quick Reference - Morphe Git Diff Plugin

## What It Does

Compares two versions of a Morphe schema and generates a semantic diff artifact in `KA:MD1:YAML1` format.

## Input

Two directory structures, each containing:
```
base/
├── models/      # .mod files
├── entities/    # .ent files
├── enums/       # .enum files
└── structures/  # .str files

head/
├── models/      # .mod files (modified)
├── entities/    # .ent files (modified)
├── enums/       # .enum files (modified)
└── structures/  # .str files (modified)
```

## Output

Single YAML file with semantic changes:
```yaml
metadata:
  spec_version: KA:MD1:YAML1
  source:
    version: base
    timestamp: "2025-01-28T15:00:00Z"
  target:
    version: head
    timestamp: "2025-01-28T15:00:00Z"
  summary:
    total_changes: 3
    breaking: 1
    additive: 2
    safe: 0
  generated_at: "2025-01-28T15:00:00Z"
  generator: morphe-git-morphediff@1.0.0

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
  
  - operation: modify
    type: field
    target:
      model: User
      field: Email
    changes:
      attributes:
        before:
          - nullable
        after:
          - mandatory
    classification: breaking
  
  - operation: add
    type: model
    target:
      model: Organization
    definition:
      fields:
        ID:
          type: UUID
        Name:
          type: String
      identifiers:
        primary:
          - ID
      related: {}
    classification: additive
```

## Change Types Detected

### Models
- ✅ Added/removed models
- ✅ Added/removed fields
- ✅ Field type changes
- ✅ Field attribute changes
- ✅ Added/removed relationships
- ✅ Relationship type changes
- ✅ Polymorphic relationship changes

### Entities
- ✅ Added/removed entities
- ✅ Added/removed fields
- ✅ Field path changes
- ✅ Added/removed relationships

### Enums
- ✅ Added/removed enums
- ✅ Added/removed entries
- ✅ Modified entry values
- ✅ Enum type changes

### Structures
- ✅ Added/removed structures
- ✅ Added/removed fields
- ✅ Field type changes

## Classification Rules

**Breaking Changes:**
- Any removal operation
- Field type incompatibility (String → Integer)
- Nullable → Mandatory attribute change
- Enum entry removal/modification
- Relationship cardinality change

**Additive Changes:**
- Add model/entity/structure/enum
- Add nullable/optional field
- Add new relationship
- Add enum entries

**Safe Changes:**
- (Future: renames with fingerprinting)

## Example Use Cases

### SQL Migration Generation
Use diff artifacts to generate PostgreSQL migrations:
```
[Diff Artifact] → [SQL Migration Plugin] → migration.sql
```

### API Documentation
Generate changelog from diffs:
```
[Diff Artifact] → [Changelog Plugin] → CHANGELOG.md
```

### Type Updates
Incrementally update TypeScript types:
```
[Diff Artifact] → [TS Types Plugin] → updated-types.ts
```

## Testing

Run the test suite:
```bash
go test ./pkg/compile/... -v
```

Test coverage includes:
- 7 entity comparison tests
- 7 enum comparison tests
- 9 model comparison tests
- 6 structure comparison tests
- 2 integration tests

## Building

Build WASM binary:
```bash
./scripts/build.sh      # Linux/Mac
./scripts/build.bat     # Windows
```

Output: `dist/morphe-git-morphediff-v1.0.0.wasm`
