package compile

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
)

// MorpheDiffConfig holds configuration for diff generation
type MorpheDiffConfig struct {
	BaseRegistryConfig rcfg.MorpheLoadRegistryConfig
	HeadRegistryConfig rcfg.MorpheLoadRegistryConfig
	OutputPath         string
	BaseVersion        string
	HeadVersion        string
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

	// Create diff document
	diffDoc := diffdef.NewDiffDocument(config.BaseVersion, config.HeadVersion)

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

	// Write output
	writer := NewMorpheDiffWriter(config.OutputPath)
	if err := writer.WriteDiff(diffDoc); err != nil {
		return fmt.Errorf("failed to write diff: %w", err)
	}

	fmt.Printf("Generated morphe diff with %d changes (%d breaking, %d additive, %d safe)\n",
		diffDoc.Metadata.Summary.TotalChanges,
		diffDoc.Metadata.Summary.Breaking,
		diffDoc.Metadata.Summary.Additive,
		diffDoc.Metadata.Summary.Safe)

	return nil
}
