package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile"
)

// CompileConfig represents the configuration passed to the plugin
type CompileConfig struct {
	BaseInputPath string                 `json:"baseInputPath"`
	HeadInputPath string                 `json:"headInputPath"`
	OutputPath    string                 `json:"outputPath"`
	Config        map[string]interface{} `json:"config,omitempty"`
	Verbose       bool                   `json:"verbose,omitempty"`
}

// Exit codes
const (
	ExitSuccess         = 0
	ExitCompileFailed   = 1
	ExitMissingConfig   = 3
	ExitInvalidConfig   = 4
	ExitInputPathError  = 12
	ExitOutputPathError = 13
)

// logInfo prints info messages only when verbose mode is enabled
func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	// Check command line arguments
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-git-morphediff <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON string with baseInputPath, headInputPath, outputPath")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Example:")
		fmt.Fprintln(os.Stderr, `  plugin-morphe-git-morphediff '{"baseInputPath":"./base","headInputPath":"./head","outputPath":"./diff.yaml","verbose":true}'`)
		os.Exit(ExitMissingConfig)
	}

	// Parse configuration
	rawConfig := os.Args[1]
	var compileConfig CompileConfig
	if err := json.Unmarshal([]byte(rawConfig), &compileConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ExitInvalidConfig)
	}

	// Validate required fields
	if compileConfig.BaseInputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: baseInputPath is required")
		os.Exit(ExitInputPathError)
	}

	if compileConfig.HeadInputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: headInputPath is required")
		os.Exit(ExitInputPathError)
	}

	if compileConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: outputPath is required")
		os.Exit(ExitOutputPathError)
	}

	// Convert to absolute paths
	baseInputAbs, err := filepath.Abs(compileConfig.BaseInputPath)
	if err == nil {
		compileConfig.BaseInputPath = baseInputAbs
	}

	headInputAbs, err := filepath.Abs(compileConfig.HeadInputPath)
	if err == nil {
		compileConfig.HeadInputPath = headInputAbs
	}

	outputAbs, err := filepath.Abs(compileConfig.OutputPath)
	if err == nil {
		compileConfig.OutputPath = outputAbs
	}

	logInfo(compileConfig.Verbose, "Processing base registry from: '%s'", compileConfig.BaseInputPath)
	logInfo(compileConfig.Verbose, "Processing head registry from: '%s'", compileConfig.HeadInputPath)
	logInfo(compileConfig.Verbose, "Output diff to: '%s'", compileConfig.OutputPath)

	// Initialize the diff configuration
	diffConfig := compile.MorpheDiffConfig{
		BaseRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     filepath.Join(compileConfig.BaseInputPath, "models"),
			RegistryEntitiesDirPath:   filepath.Join(compileConfig.BaseInputPath, "entities"),
			RegistryEnumsDirPath:      filepath.Join(compileConfig.BaseInputPath, "enums"),
			RegistryStructuresDirPath: filepath.Join(compileConfig.BaseInputPath, "structures"),
		},
		HeadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     filepath.Join(compileConfig.HeadInputPath, "models"),
			RegistryEntitiesDirPath:   filepath.Join(compileConfig.HeadInputPath, "entities"),
			RegistryEnumsDirPath:      filepath.Join(compileConfig.HeadInputPath, "enums"),
			RegistryStructuresDirPath: filepath.Join(compileConfig.HeadInputPath, "structures"),
		},
		OutputPath:  compileConfig.OutputPath,
		BaseVersion: "base",
		HeadVersion: "head",
	}

	// Run diff generation
	logInfo(compileConfig.Verbose, "Starting diff generation...")
	if err := compile.MorpheToMorpheDiff(diffConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Diff generation failed:", err)
		os.Exit(ExitCompileFailed)
	}

	logInfo(compileConfig.Verbose, "Diff generation completed successfully")
	os.Exit(ExitSuccess)
}
