@echo off

echo Building Morphe Git Diff Plugin...

REM Ensure dist directory exists
if not exist dist mkdir dist

REM Build WASM binary
set GOOS=wasip1
set GOARCH=wasm
go build -o dist/morphe-git-morphediff-v1.0.0.wasm ./cmd/plugin

echo Build complete: dist/morphe-git-morphediff-v1.0.0.wasm
