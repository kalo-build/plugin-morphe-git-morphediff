#!/bin/bash

set -e

echo "Building Morphe Git Diff Plugin..."

# Ensure dist directory exists
mkdir -p dist

# Build WASM binary
GOOS=wasip1 GOARCH=wasm go build -o dist/morphe-git-morphediff-v1.0.0.wasm ./cmd/plugin

echo "Build complete: dist/morphe-git-morphediff-v1.0.0.wasm"
