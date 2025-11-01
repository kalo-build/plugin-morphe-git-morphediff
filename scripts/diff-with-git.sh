#!/bin/bash
# Morphe Git Diff - Generate semantic schema diffs using git refs
# Usage: ./diff-with-git.sh [BASE_REF] [HEAD_REF] [MORPHE_PATH] [OUTPUT]

set -e

BASE_REF=${1:-main}
HEAD_REF=${2:-HEAD}
MORPHE_PATH=${3:-morphe}
OUTPUT=${4:-morphe-diff.yaml}

echo "📦 Morphe Git Diff"
echo "   Comparing: $MORPHE_PATH"
echo "   Base: $BASE_REF"
echo "   Head: $HEAD_REF"
echo ""

# Validate git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not a git repository"
    exit 1
fi

# Validate base ref exists
if ! git rev-parse "$BASE_REF" > /dev/null 2>&1; then
    echo "❌ Error: Base ref '$BASE_REF' not found"
    exit 1
fi

# Create temp directory for base extraction
TEMP_BASE=$(mktemp -d)
trap "rm -rf $TEMP_BASE" EXIT

# Extract base version from git
echo "📂 Extracting base from $BASE_REF..."
if ! git archive "$BASE_REF" "$MORPHE_PATH/" 2>/dev/null | tar -x -C "$TEMP_BASE" 2>/dev/null; then
    echo "❌ Error: Failed to extract $MORPHE_PATH from $BASE_REF"
    echo "   Make sure the path exists in the base ref"
    exit 1
fi

# Determine head path
if [ "$HEAD_REF" = "HEAD" ]; then
    # Use working directory (includes uncommitted changes)
    HEAD_PATH="./$MORPHE_PATH"
    echo "📂 Using working directory for head"
    
    if [ ! -d "$HEAD_PATH" ]; then
        echo "❌ Error: $HEAD_PATH directory not found"
        exit 1
    fi
else
    # Extract head from git too
    TEMP_HEAD=$(mktemp -d)
    trap "rm -rf $TEMP_BASE $TEMP_HEAD" EXIT
    
    echo "📂 Extracting head from $HEAD_REF..."
    if ! git archive "$HEAD_REF" "$MORPHE_PATH/" 2>/dev/null | tar -x -C "$TEMP_HEAD" 2>/dev/null; then
        echo "❌ Error: Failed to extract $MORPHE_PATH from $HEAD_REF"
        exit 1
    fi
    HEAD_PATH="$TEMP_HEAD/$MORPHE_PATH"
fi

# Run diff plugin
echo "🔍 Generating diff..."

# Check if plugin binary exists (for local development)
if [ -f "./dist/morphe-git-morphediff-v1.0.0.wasm" ]; then
    PLUGIN_PATH="./dist/morphe-git-morphediff-v1.0.0.wasm"
    echo "   Using local plugin: $PLUGIN_PATH"
else
    # Assume kalo will resolve from registry
    PLUGIN_PATH="@kalo-build/plugin-morphe-git-morphediff"
fi

# Build config JSON
CONFIG=$(cat <<EOF
{
  "baseInputPath": "$TEMP_BASE/$MORPHE_PATH",
  "headInputPath": "$HEAD_PATH",
  "outputPath": "$OUTPUT",
  "verbose": false
}
EOF
)

# Execute plugin
if command -v kalo &> /dev/null; then
    kalo compile --plugin "$PLUGIN_PATH" --config "$CONFIG"
else
    echo "❌ Error: kalo CLI not found in PATH"
    echo "   Install kalo CLI or run the plugin directly"
    exit 1
fi

echo ""
echo "✅ Diff generated: $OUTPUT"
echo ""

# Show summary if yq is available
if command -v yq &> /dev/null; then
    echo "Summary:"
    yq '.metadata.summary' "$OUTPUT"
    echo ""
    
    # Warn about breaking changes
    BREAKING=$(yq '.metadata.summary.breaking' "$OUTPUT")
    if [ "$BREAKING" != "0" ]; then
        echo "⚠️  WARNING: $BREAKING breaking change(s) detected!"
        echo "   Review carefully before deploying."
    fi
fi

