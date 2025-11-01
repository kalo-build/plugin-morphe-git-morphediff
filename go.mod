module github.com/kalo-build/plugin-morphe-git-morphediff

go 1.21.6

toolchain go1.24.2

require (
	github.com/kalo-build/morphe-go v0.0.0
	github.com/stretchr/testify v1.9.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kalo-build/clone v0.0.0-20250329082958-41db0353412f // indirect
	github.com/kalo-build/go-util v0.0.0-20250329083327-00e97aeff9b7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

replace github.com/kalo-build/morphe-go => ../morphe-go
