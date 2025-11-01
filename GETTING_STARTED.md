# Getting Started with Morphe Git Diff Plugin

## Quick Start (5 minutes)

### 1. Build the Plugin

```bash
# Linux/Mac
./scripts/build.sh

# Windows
scripts\build.bat
```

Output: `dist/morphe-git-morphediff-v1.0.0.wasm`

### 2. Prepare Your Morphe Schemas

You need two versions of your morphe schema:

```
my-project/
├── base/           # Original version
│   ├── models/
│   ├── entities/
│   ├── enums/
│   └── structures/
└── head/           # Modified version
    ├── models/
    ├── entities/
    ├── enums/
    └── structures/
```

### 3. Generate Diff

```bash
# Using the CLI (if integrated with kalo)
kalo compile --plugin @kalo-build/plugin-morphe-git-morphediff \
  --base ./base \
  --head ./head \
  --output ./schema-diff.yaml
```

Or run directly:
```bash
./plugin '{"baseInputPath":"./base","headInputPath":"./head","outputPath":"./diff.yaml","verbose":true}'
```

### 4. Review Output

```yaml
# diff.yaml
metadata:
  spec_version: KA:MD1:YAML1
  summary:
    total_changes: 5
    breaking: 1
    additive: 4
    safe: 0

changes:
  - operation: add
    type: field
    target:
      model: User
      field: PhoneNumber
    classification: additive
  
  - operation: modify
    type: field
    target:
      model: User
      field: Email
    changes:
      attributes:
        before: [nullable]
        after: [mandatory]
    classification: breaking  # ⚠️ Warning!
```

## Example Scenarios

### Scenario 1: Adding a New Field

**Base**:
```yaml
# user.mod
name: User
fields:
  ID:
    type: UUID
  Name:
    type: String
```

**Head**:
```yaml
# user.mod
name: User
fields:
  ID:
    type: UUID
  Name:
    type: String
  Email:  # New field
    type: String
    attributes:
      - nullable
```

**Generated Diff**:
```yaml
changes:
  - operation: add
    type: field
    target:
      model: User
      field: Email
    definition:
      type: String
      attributes: [nullable]
    classification: additive  # ✅ Safe to add
```

### Scenario 2: Breaking Change

**Base**:
```yaml
fields:
  Age:
    type: String
```

**Head**:
```yaml
fields:
  Age:
    type: Integer  # Type changed!
```

**Generated Diff**:
```yaml
changes:
  - operation: modify
    type: field
    target:
      model: User
      field: Age
    changes:
      type:
        before: String
        after: Integer
    classification: breaking  # ⚠️ Requires migration
```

### Scenario 3: Adding Enum Entry

**Base**:
```yaml
# user-role.enum
name: UserRole
type: String
entries:
  Admin: ADMIN
  Viewer: VIEWER
```

**Head**:
```yaml
# user-role.enum
name: UserRole
type: String
entries:
  Admin: ADMIN
  Editor: EDITOR  # New entry
  Viewer: VIEWER
```

**Generated Diff**:
```yaml
changes:
  - operation: modify
    type: enum
    target:
      enum: UserRole
    changes:
      entries:
        added:
          Editor: EDITOR
    classification: additive  # ✅ Backward compatible
```

## Running Tests

### Unit Tests
```bash
go test ./pkg/compile/... -v
```

### Integration Tests
```bash
go test ./pkg/compile/compile_test.go -v
```

### All Tests with Coverage
```bash
go test ./... -v -cover
```

Expected output:
```
PASS
ok      github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile  1.2s  coverage: 83.1%
```

## Understanding the Output

### Metadata Section
```yaml
metadata:
  spec_version: KA:MD1:YAML1  # Format version
  source:
    version: base              # Base version identifier
    timestamp: "..."           # When base was captured
  target:
    version: head              # Target version identifier
    timestamp: "..."           # When target was captured
  summary:
    total_changes: 8           # Total number of changes
    breaking: 1                # Breaking changes count
    additive: 7                # Additive changes count
    safe: 0                    # Safe changes count
    by_type:
      model: 1                 # Changes by artifact type
      field: 4
      enum: 1
      relationship: 2
```

### Change Operations
Each change includes:
- `operation`: add, remove, or modify
- `type`: model, entity, enum, structure, field, or relationship
- `target`: Location of the change (model name, field name, etc.)
- `definition`: Complete definition (for add operations)
- `changes`: Before/after details (for modify operations)
- `classification`: breaking, additive, or safe

## Common Use Cases

### 1. Pre-Deployment Safety Check
```bash
# Generate diff from your feature branch
kalo diff --base main --head feature/add-user-fields

# Review breaking changes before deploying
grep "classification: breaking" diff.yaml
```

### 2. SQL Migration Generation
```bash
# Generate semantic diff
kalo diff --base v1.0 --head v1.1 -o schema-changes.yaml

# Feed to SQL migration plugin
kalo compile --plugin psql-migration --input schema-changes.yaml -o migration.sql
```

### 3. API Documentation
```bash
# Generate diff
kalo diff --base v1.0 --head v2.0 -o changelog-data.yaml

# Generate markdown changelog
kalo compile --plugin changelog-generator --input changelog-data.yaml -o CHANGELOG.md
```

## Troubleshooting

### Issue: "failed to load base registry"
**Solution**: Ensure base registry path has valid morphe files
```bash
ls base/models/  # Should show .mod files
ls base/enums/   # Should show .enum files
```

### Issue: "No changes detected"
**Solution**: Verify schemas are actually different
```bash
diff -r base/ head/  # Should show differences
```

### Issue: Test failures
**Solution**: Run go mod tidy
```bash
go mod tidy
go test ./... -v
```

## Next Steps

1. ✅ Generate your first diff
2. ✅ Review the classifications
3. ✅ Feed to downstream plugins
4. 🚀 Integrate into CI/CD pipeline

## Support

Check the comprehensive documentation:
- `README.md` - Overview
- `REQUIREMENTS.md` - Technical details
- `SAMPLE_OUTPUT.md` - Real examples
- `QUICK_REFERENCE.md` - Command cheat sheet

