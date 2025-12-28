# Makefile for gobpftool

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Binary name
BINARY_NAME = gobpftool

# Build flags
LDFLAGS = -ldflags "-X gobpftool/cmd.Version=$(VERSION) -X gobpftool/cmd.GitCommit=$(GIT_COMMIT) -X gobpftool/cmd.BuildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Build for release (with optimizations)
.PHONY: release
release:
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME) .

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run integration tests (requires Docker)
.PHONY: integration-test
integration-test:
	./integration_test/run_test.sh

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
.PHONY: install
install: build
	cp $(BINARY_NAME) /usr/local/bin/

# Show version information that would be embedded
.PHONY: version-info
version-info:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary with version info"
	@echo "  release       - Build optimized binary for release"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  integration-test - Run integration tests (requires Docker)"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  install       - Install binary to /usr/local/bin"
	@echo "  version-info  - Show version info that would be embedded"
	@echo "  help          - Show this help message"