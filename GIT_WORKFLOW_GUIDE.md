# Git Workflow Guide - Morphe Diff Plugin

## The Challenge

This plugin compares two Morphe registry states. In real projects, you typically have:
- **One morphe directory** in your repository
- **Multiple git versions** (branches, tags, commits)

The plugin needs both states as directories to compare, but git stores them as history.

## Current Implementation

✅ **Plugin is a pure comparator**: Takes two directory paths, generates diff  
✅ **Works perfectly for**: Testing, manual workflows, pre-extracted states  
❌ **Doesn't handle**: Git extraction, ref resolution, temp file management  

**This is by design** - WASM plugins are sandboxed and shouldn't have git access.

## Solutions

### 🚀 Short-Term: Shell Script (Available Now)

Use the provided helper script for git-based diffing:

**File**: `scripts/diff-with-git.sh`
```bash
#!/bin/bash
# Generate morphe diff using git refs

set -e

BASE_REF=${1:-main}
HEAD_REF=${2:-HEAD}
MORPHE_PATH=${3:-morphe}
OUTPUT=${4:-morphe-diff.yaml}

echo "📦 Comparing $MORPHE_PATH: $BASE_REF → $HEAD_REF"

# Create temp directory for base extraction
TEMP_BASE=$(mktemp -d)
trap "rm -rf $TEMP_BASE" EXIT

# Extract base version from git
echo "  Extracting base ($BASE_REF)..."
git archive $BASE_REF $MORPHE_PATH/ | tar -x -C $TEMP_BASE

# Determine head path
if [ "$HEAD_REF" = "HEAD" ]; then
    # Use working directory (unstaged changes included)
    HEAD_PATH="./$MORPHE_PATH"
else
    # Extract head from git too
    TEMP_HEAD=$(mktemp -d)
    trap "rm -rf $TEMP_BASE $TEMP_HEAD" EXIT
    echo "  Extracting head ($HEAD_REF)..."
    git archive $HEAD_REF $MORPHE_PATH/ | tar -x -C $TEMP_HEAD
    HEAD_PATH="$TEMP_HEAD/$MORPHE_PATH"
fi

# Run diff plugin
echo "  Generating diff..."
kalo compile --plugin @kalo-build/plugin-morphe-git-morphediff \
  --config "{\"baseInputPath\":\"$TEMP_BASE/$MORPHE_PATH\",\"headInputPath\":\"$HEAD_PATH\",\"outputPath\":\"$OUTPUT\",\"verbose\":false}"

echo "✅ Diff generated: $OUTPUT"

# Show summary
if command -v yq &> /dev/null; then
    echo ""
    echo "Summary:"
    yq '.metadata.summary' $OUTPUT
fi
```

**File**: `scripts/diff-with-git.bat` (Windows)
```batch
@echo off
setlocal

set BASE_REF=%1
set HEAD_REF=%2
set MORPHE_PATH=%3
set OUTPUT=%4

if "%BASE_REF%"=="" set BASE_REF=main
if "%HEAD_REF%"=="" set HEAD_REF=HEAD
if "%MORPHE_PATH%"=="" set MORPHE_PATH=morphe
if "%OUTPUT%"=="" set OUTPUT=morphe-diff.yaml

echo Comparing %MORPHE_PATH%: %BASE_REF% to %HEAD_REF%

REM Create temp directory
set TEMP_BASE=%TEMP%\morphe-diff-base-%RANDOM%
mkdir %TEMP_BASE%

REM Extract base
git archive %BASE_REF% %MORPHE_PATH% | tar -x -C %TEMP_BASE%

REM Run diff
kalo compile --plugin @kalo-build/plugin-morphe-git-morphediff ^
  --config "{\"baseInputPath\":\"%TEMP_BASE%\\%MORPHE_PATH%\",\"headInputPath\":\".\\%MORPHE_PATH%\",\"outputPath\":\"%OUTPUT%\"}"

REM Cleanup
rmdir /s /q %TEMP_BASE%

echo Diff generated: %OUTPUT%
```

**Usage:**
```bash
# Compare current changes against main
./scripts/diff-with-git.sh main HEAD

# Compare two releases
./scripts/diff-with-git.sh v1.0.0 v2.0.0

# Custom morphe path
./scripts/diff-with-git.sh main HEAD src/schema

# Custom output
./scripts/diff-with-git.sh main HEAD morphe ./migrations/schema-diff.yaml
```

**Advantages:**
- ✅ Works immediately
- ✅ Simple to understand
- ✅ Easy to debug (can inspect temp files)
- ✅ Cross-platform (bash + batch versions)
- ✅ No kalo CLI changes needed

**Disadvantages:**
- ❌ Separate script to remember
- ❌ Not integrated into kalo workflow
- ❌ Manual cleanup if script fails

---

### 🎯 Medium-Term: Kalo Utility Execution (`kalo kx`)

Future vision for a generic utility execution system:

#### Concept

Similar to `npx` for npm, `kalo kx` runs utilities from the kalo registry:

```bash
# Execute utilities from registry
kalo kx @kamd/diff --base main --head HEAD
kalo kx @kamd/migrate --direction up
kalo kx @kamd/validate --strict
kalo kx @kamd/init --template backend-api
```

#### Architecture

**Two-Tier Plugin System:**

**Tier 1: Compile Plugins (WASM, sandboxed)**
- Purpose: Pure data transformation
- Sandbox: Yes (WASM)
- Access: File I/O only within mounted paths
- Examples: morphe→sql, morphe→typescript
- Registry: `@kalo-build/plugin-*`

**Tier 2: Utilities (Native binaries, system access)**
- Purpose: Orchestration, git integration, system operations
- Sandbox: No (native executable)
- Access: Full system, git, network, etc.
- Examples: Diff generators, migration runners, validators
- Registry: `@kamd/*`, `@kaoa/*`, etc.

#### Registry Schema Extension

```sql
-- New table: utilities (alongside plugins)
CREATE TABLE utilities (
  id UUID PRIMARY KEY,
  namespace VARCHAR(50),      -- @kamd, @kaoa, etc.
  name VARCHAR(100),           -- diff, migrate, validate
  description TEXT,
  version VARCHAR(50),
  
  -- Multi-platform binaries
  binary_windows_amd64_url TEXT,
  binary_linux_amd64_url TEXT,
  binary_darwin_arm64_url TEXT,
  binary_darwin_amd64_url TEXT,
  
  -- Metadata
  output_format VARCHAR(50),   -- KA:MD1:YAML1, etc.
  requires_git BOOLEAN,
  created_at TIMESTAMP
);
```

#### Utility Structure

```
@kamd/diff/
├── cmd/
│   └── diff/
│       └── main.go          # Native Go binary
├── internal/
│   ├── git/
│   │   └── extractor.go     # Git file extraction
│   └── plugin/
│       └── runner.go        # WASM plugin execution
├── dist/
│   ├── diff-v1.0.0-windows-amd64.exe
│   ├── diff-v1.0.0-linux-amd64
│   └── diff-v1.0.0-darwin-arm64
└── README.md
```

#### Developer Experience

**First Use (Auto-Install):**
```bash
$ kalo kx @kamd/diff --base main --head HEAD

📦 Installing @kamd/diff@1.0.0...
✓ Downloaded (2.3MB)
✓ Verified checksum
✓ Cached to ~/.kalo/utilities/

🔍 Extracting morphe files from git...
  Base: main (abc123)
  Head: HEAD (def456)
  
✅ Generated morphe-diff.yaml
   8 changes (1 breaking, 7 additive)
```

**Subsequent Uses (Cached):**
```bash
$ kalo kx @kamd/diff --base main

🔍 Comparing schemas...
✅ Generated morphe-diff.yaml
```

#### kalo.yaml Integration

```yaml
# Optional: Configure default refs
config:
  "@kamd/diff":
    default_base_ref: "main"
    default_head_ref: "HEAD"
    morphe_path: "morphe/"
    output_path: "morphe-diff.yaml"

# Use in pipelines
pipelines:
  pre-migration:
    stages:
      - name: "generate-diff"
        steps:
          - utility: "@kamd/diff"
            args: ["--base", "main", "--head", "HEAD"]
```

#### Kalo CLI Implementation

```go
// cmd/kalo/kx.go
func kxCommand(utilityName string, args []string) error {
    // 1. Parse utility identifier
    parsed := parseUtility(utilityName)  // @kamd/diff@1.0.0
    
    // 2. Check cache
    cachedPath := checkCache(parsed)
    if cachedPath == "" {
        // 3. Download from registry
        cachedPath = downloadUtility(parsed)
    }
    
    // 4. Execute with args
    cmd := exec.Command(cachedPath, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
```

**Cache Structure:**
```
~/.kalo/
├── plugins/                 # WASM plugins (existing)
│   └── morphe-psql-types/
└── utilities/               # Native utilities (new)
    └── @kamd/
        └── diff/
            ├── 1.0.0/
            │   ├── diff        # Native binary
            │   └── metadata.json
            └── 1.1.0/
                └── diff
```

---

### 🔄 Comparison: Script vs kalo kx

| Aspect | Shell Script | kalo kx |
|--------|-------------|---------|
| **Install** | Copy to project | Auto-download from registry |
| **Update** | Manual | `kalo kx @kamd/diff@latest` |
| **Cross-platform** | Need .sh + .bat | Single command |
| **Discovery** | Have to know it exists | `kalo search utilities` |
| **Versioning** | Manual | Built-in: `@kamd/diff@1.2.0` |
| **Integration** | Script in git | Declarative in kalo.yaml |
| **Development Time** | 1 hour | 3-5 days |
| **Maintenance** | Per-project | Centralized |

---

## 🎯 Recommended Path Forward

### Phase 1: Now (This Week)
**Provide quality shell scripts** in the plugin:

```
plugin-morphe-git-morphediff/
├── scripts/
│   ├── build.sh
│   ├── build.bat
│   ├── diff-with-git.sh      # NEW: Git-aware wrapper
│   └── diff-with-git.bat     # NEW: Windows version
```

**Document in README:**
```markdown
## Git-Based Diffing

For comparing git refs, use the helper script:

    ./scripts/diff-with-git.sh main HEAD

This automatically extracts files from git and generates the diff.
```

### Phase 2: Next Month (When Registry Needs Expansion)
**Build `kalo kx` infrastructure** when:
- 3+ utilities identified that need this
- Registry schema update planned anyway
- Clear use cases for other ecosystems (OpenAPI, Protobuf, etc.)

**Utilities to implement:**
- `@kamd/diff` - Git-aware morphe diff generator
- `@kamd/migrate` - Database migration runner
- `@kamd/validate` - Schema validator
- `@kamd/init` - Project initializer

### Phase 3: Eventually (As Pattern Matures)
**Extend kalo kx capabilities:**
- Utility compositions (diff → migrate pipeline)
- Custom utility repositories (private utilities)
- Plugin-utility dependencies (plugin declares required utilities)

---

## Implementation Details (Future Reference)

### @kamd/diff Utility Internals

```go
// cmd/diff/main.go
package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

func main() {
    baseRef := flag.String("base", "main", "Base git ref")
    headRef := flag.String("head", "HEAD", "Head git ref") 
    morphePath := flag.String("path", "morphe/", "Morphe directory in repo")
    output := flag.String("output", "morphe-diff.yaml", "Output file")
    flag.Parse()

    // Extract base from git
    baseTempDir, err := extractGitRef(*baseRef, *morphePath)
    if err != nil {
        fatal("Failed to extract base: %v", err)
    }
    defer os.RemoveAll(baseTempDir)

    // Determine head path
    var headPath string
    if *headRef == "HEAD" {
        // Use working directory
        headPath = *morphePath
    } else {
        // Extract from git
        headTempDir, err := extractGitRef(*headRef, *morphePath)
        if err != nil {
            fatal("Failed to extract head: %v", err)
        }
        defer os.RemoveAll(headTempDir)
        headPath = filepath.Join(headTempDir, *morphePath)
    }

    // Run WASM plugin
    pluginPath := resolvePlugin("@kalo-build/plugin-morphe-git-morphediff")
    runWasmPlugin(pluginPath, map[string]string{
        "baseInputPath": baseTempDir,
        "headInputPath": headPath,
        "outputPath": *output,
    })
}

func extractGitRef(ref string, path string) (string, error) {
    tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("kamd-diff-%d", os.Getpid()))
    os.MkdirAll(tempDir, 0755)
    
    cmd := exec.Command("git", "archive", ref, path)
    cmd.Dir = "."
    
    // Pipe to tar extraction
    tarCmd := exec.Command("tar", "-x", "-C", tempDir)
    tarCmd.Stdin, _ = cmd.StdoutPipe()
    
    tarCmd.Start()
    cmd.Run()
    tarCmd.Wait()
    
    return tempDir, nil
}
```

### Kalo CLI kx Implementation

```go
// cmd/kalo/kx.go
package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
)

// UtilityExecutor handles kalo kx commands
type UtilityExecutor struct {
    registryClient *RegistryClient
    cacheDir       string  // ~/.kalo/utilities/
}

func (e *UtilityExecutor) Execute(utilitySpec string, args []string) error {
    // Parse: @kamd/diff@1.0.0 or @kamd/diff (latest)
    parsed := parseUtilitySpec(utilitySpec)
    
    // Check cache
    cachedBinary := e.checkCache(parsed)
    
    if cachedBinary == "" {
        // Fetch from registry
        utility, err := e.registryClient.GetUtility(parsed)
        if err != nil {
            return fmt.Errorf("utility not found: %w", err)
        }
        
        // Download platform-specific binary
        cachedBinary, err = e.downloadAndCache(utility)
        if err != nil {
            return fmt.Errorf("download failed: %w", err)
        }
    }
    
    // Execute
    cmd := exec.Command(cachedBinary, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    
    return cmd.Run()
}

func (e *UtilityExecutor) downloadAndCache(utility Utility) (string, error) {
    // Determine platform-specific URL
    var downloadURL string
    switch runtime.GOOS {
    case "windows":
        downloadURL = utility.BinaryWindowsURL
    case "linux":
        downloadURL = utility.BinaryLinuxURL
    case "darwin":
        downloadURL = utility.BinaryDarwinURL
    }
    
    // Download to cache
    cacheDir := filepath.Join(e.cacheDir, utility.Name, utility.Version)
    os.MkdirAll(cacheDir, 0755)
    
    binaryName := utility.Name
    if runtime.GOOS == "windows" {
        binaryName += ".exe"
    }
    
    binaryPath := filepath.Join(cacheDir, binaryName)
    
    // Download
    resp, err := http.Get(downloadURL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    // Write and verify checksum
    f, _ := os.Create(binaryPath)
    defer f.Close()
    
    h := sha256.New()
    io.Copy(io.MultiWriter(f, h), resp.Body)
    
    if hex.EncodeToString(h.Sum(nil)) != utility.Checksum {
        return "", fmt.Errorf("checksum mismatch")
    }
    
    // Make executable
    os.Chmod(binaryPath, 0755)
    
    return binaryPath, nil
}
```

---

## Use Cases Comparison

### Use Case 1: Daily Development

**With Script:**
```bash
# Make schema changes
vim morphe/models/user.mod

# Generate diff
./scripts/diff-with-git.sh main

# Review
cat morphe-diff.yaml
```

**With kalo kx:**
```bash
# Make schema changes
vim morphe/models/user.mod

# Generate diff (auto-installs if needed)
kalo kx @kamd/diff --base main

# Review
cat morphe-diff.yaml
```

### Use Case 2: CI/CD Pipeline

**With Script:**
```yaml
# .github/workflows/schema-check.yml
- name: Generate diff
  run: ./scripts/diff-with-git.sh origin/main HEAD

- name: Check breaking changes
  run: |
    if grep -q "classification: breaking" morphe-diff.yaml; then
      echo "::error::Breaking schema changes detected"
      exit 1
    fi
```

**With kalo kx:**
```yaml
# .github/workflows/schema-check.yml
- name: Generate diff
  run: kalo kx @kamd/diff --base origin/main --head HEAD

- name: Check breaking changes  
  run: kalo kx @kamd/validate-breaking --input morphe-diff.yaml
```

### Use Case 3: Release Changelog

**With Script:**
```bash
# Compare two releases
./scripts/diff-with-git.sh v1.0.0 v2.0.0 morphe changelog-data.yaml

# Generate markdown
kalo compile --plugin changelog-generator \
  --input changelog-data.yaml \
  --output CHANGELOG.md
```

**With kalo kx:**
```bash
# One command
kalo kx @kamd/changelog --from v1.0.0 --to v2.0.0

# Internally calls @kamd/diff, then @kamd/format-changelog
```

---

## Decision Framework

### Use Shell Script If:
- ✅ Need solution immediately
- ✅ Only morphe-diff needs git integration
- ✅ Want to validate workflow before investing
- ✅ Team is comfortable with bash scripts

### Build kalo kx If:
- ✅ 3+ utilities identified that need system access
- ✅ Want polished, integrated experience
- ✅ Building utility ecosystem is strategic goal
- ✅ Ready to expand registry infrastructure

---

## Current Recommendation

### This Week: Ship with Script

Include quality shell scripts that:
1. Handle common cases (main vs HEAD)
2. Support custom refs and paths
3. Provide clear error messages
4. Clean up temp files properly

**Document clearly:**
```markdown
## Quick Start

### For Git-Based Diffing (Recommended)
./scripts/diff-with-git.sh main HEAD

### For Directory-Based Diffing (Testing)
kalo compile --plugin morphe-git-morphediff \
  --base ./testdata/base \
  --head ./testdata/head
```

### Next Quarter: Evaluate kalo kx

After developers use the script for a month:
- Gather feedback on pain points
- Identify other utilities that need similar patterns
- Decide if generic execution is worth the investment

### Key Insight

The **plugin is perfect as-is** because it's a pure comparator. The git integration is a **developer convenience layer**, not core functionality. 

Scripts provide that convenience cheaply. `kalo kx` would provide it elegantly but requires infrastructure investment.

---

## Final Thought

Maybe the real answer is: **Both, in sequence**

1. **Week 1**: Ship script (quick value)
2. **Month 1**: Collect usage data
3. **Month 2**: Build `kalo kx` if pattern validates
4. **Month 3**: Deprecate script in favor of `kalo kx @kamd/diff`

This lets you:
- ✅ Ship immediately
- ✅ Validate assumptions
- ✅ Avoid premature abstraction
- ✅ Have migration path when ready

Sound good?

