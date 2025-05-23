.PHONY: all build test lint clean tidy run help 

# Default target
all: help

# Show help for each target
help:
	@echo "Available targets:"
	@echo "  build         Build the main server binary (bin/groupcache-server)"
	@echo "  test          Run all tests (go test ./...)"
	@echo "  lint          Run linter (golangci-lint run ./...)"
	@echo "  clean         Remove build artifacts (bin directory)"
	@echo "  tidy          Ensure dependencies are tidy (go mod tidy)"
	@echo "  run           Run the server (go run ./cmd/server)"
	@echo "  install-lint  Install golangci-lint if not present"
	@echo "  help          Show this help message"

# Build the main server binary
build:
	go build -o bin/groupcache-server ./cmd/server
	go build -o bin/groupcache-benchmark ./cmd/benchmark

# Run all tests
test:
	go test ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Remove build artifacts
clean:
	rm -rf bin

# Ensure dependencies are tidy
tidy:
	go mod tidy

# Run the server (example)
run:
	go run ./cmd/server

# Install golangci-lint if not present
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run tests with coverage reporting
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html