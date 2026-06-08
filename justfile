set shell := ["bash", "-uc"]

# list available recipes
default:
    @just --list

# upgrade all direct & indirect dependencies to their latest minor/patch, then tidy
upgrade:
    go get -u ./...
    go mod tidy

# upgrade dependencies including major versions (uses gotestsum-style discovery via go-mod-upgrade if available)
upgrade-major:
    go run github.com/oligot/go-mod-upgrade@latest

# bump the Go toolchain/language version pinned in go.mod to the latest release
upgrade-go:
    go get go@latest
    go mod tidy

# show which dependencies have newer versions available (updates appear in [brackets])
outdated:
    go list -u -m all | grep -F '[' || echo "all dependencies up to date"

# verify module checksums and dependency graph integrity
verify:
    go mod verify

# tidy go.mod / go.sum
tidy:
    go mod tidy

# run tests with coverage (matches CI)
test:
    go test -v ./... -covermode=atomic -coverprofile=coverage.out

# run benchmarks (matches CI)
bench:
    go test -bench=. -benchmem -run=^$ -count=5 ./...

# run golangci-lint (matches CI)
lint:
    golangci-lint run

# build all packages
build:
    go build ./...
