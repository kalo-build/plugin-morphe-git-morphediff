set GOOS=wasip1
set GOARCH=wasm
go build -o ../dist/morphe-git-morphediff-v1.0.0.wasm ../cmd/plugin/main.go