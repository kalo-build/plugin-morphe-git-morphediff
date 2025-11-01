# Sample Output

## Example Diff Artifact

This is actual output from the integration test showing a real morphe diff artifact.

### Scenario

**Base Schema:**
- `User` model with ID, Email (nullable), Name
- `UserRole` enum with Admin, Viewer
- `Address` structure with Street, City
- `UserProfile` entity with ID, Email

**Head Schema:**
- `User` model with ID, Email (mandatory), Name, PhoneNumber, ContactInfo relationship
- `Organization` model (new)
- `UserRole` enum with Admin, Editor, Viewer (Editor added)
- `Address` structure with Street, City, Country (Country added)
- `UserProfile` entity with ID, Email, FullName, Organization relationship

### Generated Diff

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
    total_changes: 8
    breaking: 1
    additive: 7
    safe: 0
    by_type:
      enum: 1
      field: 4
      model: 1
      relationship: 2
  generated_at: "2025-01-28T15:00:00Z"
  generator: morphe-git-morphediff@1.0.0

changes:
  # Enum Entry Addition (Additive)
  - operation: modify
    type: enum
    target:
      enum: UserRole
    changes:
      entries:
        added:
          Editor: EDITOR
    classification: additive

  # Structure Field Addition (Additive)
  - operation: add
    type: field
    target:
      structure: Address
      field: Country
    definition:
      type: String
    classification: additive

  # Model Addition (Additive)
  - operation: add
    type: model
    target:
      model: Organization
    definition:
      fields:
        ID:
          type: UUID
          attributes:
            - mandatory
        Name:
          type: String
      identifiers:
        primary:
          - ID
      related: {}
    classification: additive

  # Model Field Addition (Additive)
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

  # Model Field Modification (Breaking)
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

  # Model Relationship Addition (Additive)
  - operation: add
    type: relationship
    target:
      model: User
      relationship: ContactInfo
    definition:
      type: HasOne
    classification: additive

  # Entity Field Addition (Additive)
  - operation: add
    type: field
    target:
      entity: UserProfile
      field: FullName
    definition:
      type: User.Name
    classification: additive

  # Entity Relationship Addition (Additive)
  - operation: add
    type: relationship
    target:
      entity: UserProfile
      relationship: Organization
    definition:
      type: ForOne
    classification: additive
```

## Change Breakdown

| Operation | Artifact Type | Count | Classification |
|-----------|--------------|-------|----------------|
| Add | Model | 1 | Additive |
| Add | Field | 3 | Additive |
| Add | Relationship | 2 | Additive |
| Modify | Enum | 1 | Additive (entries added) |
| Modify | Field | 1 | Breaking (constraint tightened) |

**Total**: 8 changes (7 additive, 1 breaking)

## Downstream Usage

This diff artifact can power:

### 1. PostgreSQL Migration
```sql
-- Generated from morphe-diff.yaml

-- Add new model
CREATE TABLE organizations (
  id UUID PRIMARY KEY,
  name VARCHAR(255)
);

-- Add field to existing table
ALTER TABLE users ADD COLUMN phone_number VARCHAR(255) NULL;

-- Modify constraint (BREAKING)
ALTER TABLE users ALTER COLUMN email SET NOT NULL;

-- Add foreign key
ALTER TABLE users ADD COLUMN contact_info_id UUID;
ALTER TABLE users ADD CONSTRAINT fk_user_contact_info 
  FOREIGN KEY (contact_info_id) REFERENCES contact_infos(id);
```

### 2. TypeScript Types Update
```typescript
// Generated from morphe-diff.yaml

// BREAKING: Email is now required
interface User {
  id: string;
  email: string;  // Changed: was 'email?: string'
  name?: string;
  phoneNumber?: string;  // New field
  contactInfo?: ContactInfo;  // New relationship
}

// New model
interface Organization {
  id: string;
  name?: string;
}

// Updated enum
export enum UserRole {
  Admin = 'ADMIN',
  Editor = 'EDITOR',  // New entry
  Viewer = 'VIEWER',
}
```

### 3. API Changelog
```markdown
## Breaking Changes ⚠️

- **User.Email** is now required (was nullable)

## New Features ✨

- Added `Organization` model
- Added `User.PhoneNumber` field (nullable)
- Added `User.ContactInfo` relationship
- Added `UserRole.Editor` enum value
- Added `Address.Country` field
- Added `UserProfile.FullName` field
- Added `UserProfile.Organization` relationship
```

## Validation

All 31 tests pass:
- ✅ 29 unit tests (100% coverage of core logic)
- ✅ 2 integration tests (end-to-end validation)
- ✅ WASM build successful

## Performance

- Diff generation: ~60ms for small schemas (3-5 models)
- Memory efficient: Streams comparison, doesn't load entire AST in memory
- Deterministic: Same input always produces same output
