package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile"
)

// PluginConfig represents the configuration passed to the plugin by Kalo CLI
type PluginConfig struct {
	// New store-based paths (mounted by CLI)
	Stores map[string]StoreConfig `json:"stores,omitempty"`

	// Legacy direct paths (for backward compatibility)
	BaseInputPath string `json:"baseInputPath,omitempty"`
	HeadInputPath string `json:"headInputPath,omitempty"`
	OutputPath    string `json:"outputPath,omitempty"`

	// Plugin-specific config
	Config  map[string]interface{} `json:"config,omitempty"`
	Verbose bool                   `json:"verbose,omitempty"`
}

// StoreConfig represents a store configuration from Kalo CLI
type StoreConfig struct {
	ID        uint32 `json:"id"`
	Type      string `json:"type"`
	MountPath string `json:"mountPath,omitempty"`

	// Git provenance (for gitRepository stores)
	GitRef       string `json:"gitRef,omitempty"`
	GitCommit    string `json:"gitCommit,omitempty"`
	GitTimestamp string `json:"gitTimestamp,omitempty"`
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
		fmt.Fprintln(os.Stderr, "  config: JSON string with store configurations")
		os.Exit(ExitMissingConfig)
	}

	// Parse configuration
	rawConfig := os.Args[1]
	var config PluginConfig
	if err := json.Unmarshal([]byte(rawConfig), &config); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ExitInvalidConfig)
	}

	// Determine paths - prefer store mounts, fall back to legacy paths
	var baseInputPath, headInputPath, outputPath string

	// Check for store-based configuration (new approach)
	if config.Stores != nil {
		// Look for /base mount (from "base" named input)
		for _, store := range config.Stores {
			if store.MountPath == "/base" {
				baseInputPath = "/base"
			}
			if store.MountPath == "/head" {
				headInputPath = "/head"
			}
			if store.MountPath == "/output" {
				outputPath = "/output"
			}
			// Also support legacy /input mount as head
			if store.MountPath == "/input" && headInputPath == "" {
				headInputPath = "/input"
			}
		}
	}

	// Fall back to legacy direct paths
	if baseInputPath == "" && config.BaseInputPath != "" {
		baseInputPath = config.BaseInputPath
	}
	if headInputPath == "" && config.HeadInputPath != "" {
		headInputPath = config.HeadInputPath
	}
	if outputPath == "" && config.OutputPath != "" {
		outputPath = config.OutputPath
	}

	// Validate required paths
	if baseInputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: base input path is required (mount /base store or provide baseInputPath)")
		os.Exit(ExitInputPathError)
	}

	if headInputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: head input path is required (mount /head store or provide headInputPath)")
		os.Exit(ExitInputPathError)
	}

	if outputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: output path is required (mount /output store or provide outputPath)")
		os.Exit(ExitOutputPathError)
	}

	// If outputPath is a directory (from store mount), append default filename
	// Store mounts like /output need a filename appended
	if outputPath == "/output" || outputPath == "/output/" {
		outputPath = filepath.Join(outputPath, "morphe-diff.yaml")
	}

	// Get verbose from nested config if present
	verbose := config.Verbose
	if config.Config != nil {
		if v, ok := config.Config["verbose"].(bool); ok {
			verbose = v
		}
	}

	logInfo(verbose, "Processing base registry from: '%s'", baseInputPath)
	logInfo(verbose, "Processing head registry from: '%s'", headInputPath)
	logInfo(verbose, "Output diff to: '%s'", outputPath)

	// Extract git provenance info from store configs (passed by CLI)
	var baseRef, baseCommit, baseTimestamp string
	var headRef, headCommit, headTimestamp string
	var archiveDiffs bool

	// Get git info from store configs
	if config.Stores != nil {
		for _, store := range config.Stores {
			if store.MountPath == "/base" {
				baseRef = store.GitRef
				baseCommit = store.GitCommit
				baseTimestamp = store.GitTimestamp
			}
			if store.MountPath == "/head" {
				headRef = store.GitRef
				headCommit = store.GitCommit
				headTimestamp = store.GitTimestamp
			}
		}
	}

	// Check for archiveDiffs and outputFormat in plugin config
	var outputFormat compile.OutputFormat
	if config.Config != nil {
		if v, ok := config.Config["archiveDiffs"].(bool); ok {
			archiveDiffs = v
		}
		if v, ok := config.Config["outputFormat"].(string); ok {
			switch v {
			case "json":
				outputFormat = compile.OutputFormatJSON
			case "yaml":
				outputFormat = compile.OutputFormatYAML
			default:
				outputFormat = compile.OutputFormatYAML
			}
		}
	}

	// Default refs if not provided
	if baseRef == "" {
		baseRef = "base"
	}
	if headRef == "" {
		headRef = "head"
	}

	logInfo(verbose, "Base ref: %s (commit: %s)", baseRef, baseCommit)
	logInfo(verbose, "Head ref: %s (commit: %s)", headRef, headCommit)

	// Adjust output path extension for JSON format
	if outputFormat == compile.OutputFormatJSON && filepath.Ext(outputPath) == ".yaml" {
		outputPath = outputPath[:len(outputPath)-5] + ".json"
	}

	// Initialize the diff configuration
	diffConfig := compile.MorpheDiffConfig{
		BaseRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     filepath.Join(baseInputPath, "models"),
			RegistryEntitiesDirPath:   filepath.Join(baseInputPath, "entities"),
			RegistryEnumsDirPath:      filepath.Join(baseInputPath, "enums"),
			RegistryStructuresDirPath: filepath.Join(baseInputPath, "structures"),
		},
		HeadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     filepath.Join(headInputPath, "models"),
			RegistryEntitiesDirPath:   filepath.Join(headInputPath, "entities"),
			RegistryEnumsDirPath:      filepath.Join(headInputPath, "enums"),
			RegistryStructuresDirPath: filepath.Join(headInputPath, "structures"),
		},
		OutputPath:    outputPath,
		BaseRef:       baseRef,
		BaseCommit:    baseCommit,
		BaseTimestamp: baseTimestamp,
		HeadRef:       headRef,
		HeadCommit:    headCommit,
		HeadTimestamp: headTimestamp,
		ArchiveDiffs:  archiveDiffs,
		OutputFormat:  outputFormat,
	}

	// Run diff generation
	logInfo(verbose, "Starting diff generation...")
	if err := compile.MorpheToMorpheDiff(diffConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Diff generation failed:", err)
		os.Exit(ExitCompileFailed)
	}

	logInfo(verbose, "Diff generation completed successfully")
	os.Exit(ExitSuccess)
}
