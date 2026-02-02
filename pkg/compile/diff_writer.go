package compile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
	"gopkg.in/yaml.v3"
)

// OutputFormat represents the output format for the diff document
type OutputFormat string

const (
	OutputFormatYAML OutputFormat = "yaml"
	OutputFormatJSON OutputFormat = "json"
)

// MorpheDiffWriter handles writing diff documents to files
type MorpheDiffWriter struct {
	OutputPath string
	Format     OutputFormat
}

// NewMorpheDiffWriter creates a new diff writer with YAML format (default)
func NewMorpheDiffWriter(outputPath string) *MorpheDiffWriter {
	return &MorpheDiffWriter{
		OutputPath: outputPath,
		Format:     OutputFormatYAML,
	}
}

// NewMorpheDiffWriterWithFormat creates a new diff writer with specified format
func NewMorpheDiffWriterWithFormat(outputPath string, format OutputFormat) *MorpheDiffWriter {
	// Auto-detect format from file extension if not specified
	if format == "" {
		if strings.HasSuffix(outputPath, ".json") {
			format = OutputFormatJSON
		} else {
			format = OutputFormatYAML
		}
	}
	return &MorpheDiffWriter{
		OutputPath: outputPath,
		Format:     format,
	}
}

// WriteDiff writes a diff document to a file in the configured format
func (w *MorpheDiffWriter) WriteDiff(diffDoc *diffdef.DiffDocument) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(w.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var data []byte
	var err error

	switch w.Format {
	case OutputFormatJSON:
		data, err = json.MarshalIndent(diffDoc, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal diff document to JSON: %w", err)
		}
	default: // YAML
		data, err = yaml.Marshal(diffDoc)
		if err != nil {
			return fmt.Errorf("failed to marshal diff document to YAML: %w", err)
		}
	}

	// Write to file
	if err := os.WriteFile(w.OutputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write diff file: %w", err)
	}

	return nil
}
