package diffdef

import "time"

// DiffDocument represents a complete morphe diff artifact
type DiffDocument struct {
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Changes  []Change `yaml:"changes" json:"changes"`
}

// Metadata contains information about the diff context
type Metadata struct {
	SpecVersion string      `yaml:"spec_version" json:"spec_version"`
	Source      VersionInfo `yaml:"source" json:"source"`
	Target      VersionInfo `yaml:"target" json:"target"`
	Summary     Summary     `yaml:"summary" json:"summary"`
	GeneratedAt string      `yaml:"generated_at" json:"generated_at"`
	Generator   string      `yaml:"generator,omitempty" json:"generator,omitempty"`
}

// VersionInfo describes a schema version with git provenance
type VersionInfo struct {
	Ref       string `yaml:"ref" json:"ref"`                             // Git ref name (e.g., "base", "head", "main", "feature/xyz")
	Commit    string `yaml:"commit,omitempty" json:"commit,omitempty"`   // Git commit hash for reproducibility
	Timestamp string `yaml:"timestamp" json:"timestamp"`                 // When this version was captured
	Version   string `yaml:"version,omitempty" json:"version,omitempty"` // Deprecated: use Ref instead
}

// Summary provides change statistics
type Summary struct {
	TotalChanges int            `yaml:"total_changes" json:"total_changes"`
	Breaking     int            `yaml:"breaking" json:"breaking"`
	Additive     int            `yaml:"additive" json:"additive"`
	Safe         int            `yaml:"safe" json:"safe"`
	ByType       map[string]int `yaml:"by_type,omitempty" json:"by_type,omitempty"`
}

// Change represents a single delta operation
type Change struct {
	Operation      string                 `yaml:"operation" json:"operation"`
	Type           string                 `yaml:"type" json:"type"`
	Target         map[string]string      `yaml:"target,omitempty" json:"target,omitempty"`
	Source         map[string]string      `yaml:"source,omitempty" json:"source,omitempty"`
	Destination    map[string]string      `yaml:"destination,omitempty" json:"destination,omitempty"`
	Definition     map[string]interface{} `yaml:"definition,omitempty" json:"definition,omitempty"`
	Changes        map[string]interface{} `yaml:"changes,omitempty" json:"changes,omitempty"`
	EntitySnapshot map[string]interface{} `yaml:"entity_snapshot,omitempty" json:"entity_snapshot,omitempty"`
	RenamedTo      string                 `yaml:"renamed_to,omitempty" json:"renamed_to,omitempty"`
	Fingerprint    string                 `yaml:"fingerprint,omitempty" json:"fingerprint,omitempty"`
	Reason         string                 `yaml:"reason,omitempty" json:"reason,omitempty"`
	Classification string                 `yaml:"classification" json:"classification"`
}

// Operation types
const (
	OperationAdd       = "add"
	OperationRemove    = "remove"
	OperationModify    = "modify"
	OperationRename    = "rename"
	OperationMove      = "move"
	OperationDeprecate = "deprecate"
)

// Artifact types
const (
	TypeModel        = "model"
	TypeEntity       = "entity"
	TypeEnum         = "enum"
	TypeStructure    = "structure"
	TypeField        = "field"
	TypeRelationship = "relationship"
	TypeEnumEntry    = "enum_entry"
)

// Classifications
const (
	ClassificationBreaking = "breaking"
	ClassificationAdditive = "additive"
	ClassificationSafe     = "safe"
)

// DiffVersionConfig holds version info passed to NewDiffDocument
type DiffVersionConfig struct {
	Ref       string // Git ref name (e.g., "base", "main", "feature/xyz")
	Commit    string // Git commit hash (optional but recommended)
	Timestamp string // When this version was captured (optional, defaults to now)
}

// NewDiffDocument creates a new diff document with metadata
func NewDiffDocument(baseVersion, headVersion string) *DiffDocument {
	return NewDiffDocumentWithConfig(
		DiffVersionConfig{Ref: baseVersion},
		DiffVersionConfig{Ref: headVersion},
	)
}

// isValidTimestamp checks if a timestamp string is valid and reasonable
// (not empty and not a zero/near-zero time value)
func isValidTimestamp(ts string) bool {
	if ts == "" {
		return false
	}
	// Parse and check if it's a reasonable date (after year 2000)
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return false
	}
	// Reject zero-ish times (before year 2000)
	return t.Year() >= 2000
}

// NewDiffDocumentWithConfig creates a new diff document with full version config
func NewDiffDocumentWithConfig(source, target DiffVersionConfig) *DiffDocument {
	now := time.Now().UTC().Format(time.RFC3339)

	sourceTimestamp := source.Timestamp
	if !isValidTimestamp(sourceTimestamp) {
		sourceTimestamp = now
	}

	targetTimestamp := target.Timestamp
	if !isValidTimestamp(targetTimestamp) {
		targetTimestamp = now
	}

	return &DiffDocument{
		Metadata: Metadata{
			SpecVersion: "KA:MD1:YAML1",
			Source: VersionInfo{
				Ref:       source.Ref,
				Commit:    source.Commit,
				Timestamp: sourceTimestamp,
			},
			Target: VersionInfo{
				Ref:       target.Ref,
				Commit:    target.Commit,
				Timestamp: targetTimestamp,
			},
			Summary: Summary{
				ByType: make(map[string]int),
			},
			GeneratedAt: now,
			Generator:   "morphe-git-morphediff@1.0.0",
		},
		Changes: make([]Change, 0),
	}
}

// AddChange adds a change and updates the summary
func (d *DiffDocument) AddChange(change Change) {
	d.Changes = append(d.Changes, change)
	d.Metadata.Summary.TotalChanges++

	// Update classification counts
	switch change.Classification {
	case ClassificationBreaking:
		d.Metadata.Summary.Breaking++
	case ClassificationAdditive:
		d.Metadata.Summary.Additive++
	case ClassificationSafe:
		d.Metadata.Summary.Safe++
	}

	// Update type counts
	d.Metadata.Summary.ByType[change.Type]++
}
