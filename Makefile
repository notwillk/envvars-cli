.PHONY: help build run test clean lint format deps

# Default target
help:
	@echo "Available commands:"
	@echo "  build   - Build the application"
	@echo "  run     - Run the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean build artifacts"
	@echo "  lint    - Run linter"
	@echo "  format  - Format code"
	@echo "  deps    - Download dependencies"

# Build the application
build:
	go build -o bin/envvars-cli .

# Run the application
run:
	go run .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run linter
lint:
	golangci-lint run

# Format code
format:
	go fmt ./...
	goimports -w .

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install development tools
tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/lint/golint@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

# Run with hot reload (requires air)
dev:
	air

# Build for different platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/envvars-cli-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/envvars-cli-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/envvars-cli-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/envvars-cli-windows-amd64.exe .

