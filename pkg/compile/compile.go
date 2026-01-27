package compile

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
)

// MorpheDiffConfig holds configuration for diff generation
type MorpheDiffConfig struct {
	BaseRegistryConfig rcfg.MorpheLoadRegistryConfig
	HeadRegistryConfig rcfg.MorpheLoadRegistryConfig
	OutputPath         string

	// Version info for provenance tracking
	BaseRef       string // Git ref name (e.g., "base", "main~1")
	BaseCommit    string // Git commit hash (optional)
	BaseTimestamp string // When base was captured (optional)
	HeadRef       string // Git ref name (e.g., "head", "main")
	HeadCommit    string // Git commit hash (optional)
	HeadTimestamp string // When head was captured (optional)

	// Output options
	ArchiveDiffs bool   // If true, output timestamped files for history
	OutputFile   string // Custom output filename (optional)

	// Deprecated: use BaseRef/HeadRef instead
	BaseVersion string
	HeadVersion string
}

// MorpheToMorpheDiff generates a morphe diff document from two registry states
func MorpheToMorpheDiff(config MorpheDiffConfig) error {
	// Load base registry
	baseRegistry, err := registry.LoadMorpheRegistry(registry.LoadMorpheRegistryHooks{}, config.BaseRegistryConfig)
	if err != nil {
		return fmt.Errorf("failed to load base registry: %w", err)
	}

	// Load head registry
	headRegistry, err := registry.LoadMorpheRegistry(registry.LoadMorpheRegistryHooks{}, config.HeadRegistryConfig)
	if err != nil {
		return fmt.Errorf("failed to load head registry: %w", err)
	}

	// Create diff document with provenance info
	baseRef := config.BaseRef
	if baseRef == "" {
		baseRef = config.BaseVersion // Backwards compatibility
	}
	headRef := config.HeadRef
	if headRef == "" {
		headRef = config.HeadVersion // Backwards compatibility
	}

	diffDoc := diffdef.NewDiffDocumentWithConfig(
		diffdef.DiffVersionConfig{
			Ref:       baseRef,
			Commit:    config.BaseCommit,
			Timestamp: config.BaseTimestamp,
		},
		diffdef.DiffVersionConfig{
			Ref:       headRef,
			Commit:    config.HeadCommit,
			Timestamp: config.HeadTimestamp,
		},
	)

	// Compare enums
	if baseRegistry.HasEnums() || headRegistry.HasEnums() {
		fmt.Println("Comparing enums...")
		if err := CompareEnums(baseRegistry, headRegistry, diffDoc); err != nil {
			return fmt.Errorf("failed to compare enums: %w", err)
		}
	}

	// Compare structures
	if baseRegistry.HasStructures() || headRegistry.HasStructures() {
		fmt.Println("Comparing structures...")
		if err := CompareStructures(baseRegistry, headRegistry, diffDoc); err != nil {
			return fmt.Errorf("failed to compare structures: %w", err)
		}
	}

	// Compare models
	if baseRegistry.HasModels() || headRegistry.HasModels() {
		fmt.Println("Comparing models...")
		if err := CompareModels(baseRegistry, headRegistry, diffDoc); err != nil {
			return fmt.Errorf("failed to compare models: %w", err)
		}
	}

	// Compare entities
	if baseRegistry.HasEntities() || headRegistry.HasEntities() {
		fmt.Println("Comparing entities...")
		if err := CompareEntities(baseRegistry, headRegistry, diffDoc); err != nil {
			return fmt.Errorf("failed to compare entities: %w", err)
		}
	}

	// Determine output path
	outputPath := config.OutputPath

	// If archiveDiffs is enabled, generate a timestamped filename
	if config.ArchiveDiffs {
		dir := filepath.Dir(outputPath)
		timestamp := time.Now().UTC().Format("20060102150405")
		archivedName := fmt.Sprintf("%s_morphe-diff.yaml", timestamp)
		outputPath = filepath.Join(dir, archivedName)
		fmt.Printf("Archiving diff to: %s\n", archivedName)
	}

	// Write output
	writer := NewMorpheDiffWriter(outputPath)
	if err := writer.WriteDiff(diffDoc); err != nil {
		return fmt.Errorf("failed to write diff: %w", err)
	}

	fmt.Printf("Generated morphe diff with %d changes (%d breaking, %d additive, %d safe)\n",
		diffDoc.Metadata.Summary.TotalChanges,
		diffDoc.Metadata.Summary.Breaking,
		diffDoc.Metadata.Summary.Additive,
		diffDoc.Metadata.Summary.Safe)

	// If archiving, also write to the default location for convenience
	if config.ArchiveDiffs && outputPath != config.OutputPath {
		defaultWriter := NewMorpheDiffWriter(config.OutputPath)
		if err := defaultWriter.WriteDiff(diffDoc); err != nil {
			fmt.Printf("Warning: could not write default diff file: %v\n", err)
		}
	}

	return nil
}
