# Makefile for FamGo monorepo
# Streamlined developer workflow

.PHONY: sync deps quick test cover clean lint fmt

# Sync go.work with all modules
sync:
    go work sync

# Pre-cache all dependencies
deps:
    go mod download all

# Quick tests (no coverage, short mode)
quick:
    go test ./internal/... -short

# Full test run with coverage, parallelized
test:
    go test ./... -cover -p=8

# Coverage report (HTML)
cover:
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html

# Clean test cache
clean:
    go clean -testcache

# Lint code (requires golangci-lint installed)
lint:
    golangci-lint run ./...

# Format code
fmt:
    go fmt ./...
