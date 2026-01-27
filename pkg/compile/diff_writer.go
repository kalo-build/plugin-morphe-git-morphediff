package compile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
	"gopkg.in/yaml.v3"
)

// MorpheDiffWriter handles writing diff documents to files
type MorpheDiffWriter struct {
	OutputPath string
}

// NewMorpheDiffWriter creates a new diff writer
func NewMorpheDiffWriter(outputPath string) *MorpheDiffWriter {
	return &MorpheDiffWriter{
		OutputPath: outputPath,
	}
}

// WriteDiff writes a diff document to a YAML file
func (w *MorpheDiffWriter) WriteDiff(diffDoc *diffdef.DiffDocument) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(w.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(diffDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal diff document: %w", err)
	}

	// Write to file
	if err := os.WriteFile(w.OutputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write diff file: %w", err)
	}

	return nil
}


