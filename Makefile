# ELMOS - Embedded Linux on MacOS
# Makefile for building and installing the elmos CLI

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go build flags
LDFLAGS := -X github.com/NguyenTrongPhuc552003/elmos/pkg/version.Version=$(VERSION)
LDFLAGS += -X github.com/NguyenTrongPhuc552003/elmos/pkg/version.Commit=$(COMMIT)
LDFLAGS += -X github.com/NguyenTrongPhuc552003/elmos/pkg/version.BuildDate=$(DATE)

# Binary name
BINARY := elmos

# Directories
GOBIN := $(shell go env GOPATH)/pkg
PREFIX ?= /usr/local

.PHONY: all build install clean test lint deps fmt help

## Build targets

all: build

build: ## Build the elmos binary
	@echo "Building $(BINARY)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install: build ## Install to GOPATH/bin
	@echo "Installing to $(GOBIN)..."
	cp $(BINARY) $(GOBIN)/

install-global: build ## Install to /usr/local/bin
	@echo "Installing to $(PREFIX)/bin..."
	sudo cp $(BINARY) $(PREFIX)/bin/

uninstall: ## Remove from GOPATH/bin
	rm -f $(GOBIN)/$(BINARY)

clean: ## Clean build artifacts
	rm -f $(BINARY)
	go clean

## Development targets

deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

lint: ## Run linter
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet..."; \
		go vet ./...; \
	fi

test: ## Run tests
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Release targets

release: ## Build for all supported platforms
	@echo "Building for darwin/arm64..."
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY)-darwin-arm64 .
	@echo "Building for darwin/amd64..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY)-darwin-amd64 .

## Homebrew formula generation

homebrew-formula: ## Generate Homebrew formula
	@echo "class Elmos < Formula"
	@echo "  desc \"Embedded Linux development on macOS - native kernel builds\""
	@echo "  homepage \"https://github.com/NguyenTrongPhuc552003/elmos\""
	@echo "  version \"$(VERSION)\""
	@echo "  license \"MIT\""
	@echo ""
	@echo "  depends_on \"llvm\""
	@echo "  depends_on \"lld\""
	@echo "  depends_on \"gnu-sed\""
	@echo "  depends_on \"make\""
	@echo "  depends_on \"libelf\""
	@echo "  depends_on \"qemu\""
	@echo "  depends_on \"e2fsprogs\""
	@echo "  depends_on \"coreutils\""
	@echo ""
	@echo "  def install"
	@echo "    bin.install \"elmos\""
	@echo "  end"
	@echo ""
	@echo "  test do"
	@echo "    system \"\#{bin}/elmos\", \"version\""
	@echo "  end"
	@echo "end"

## Shell completions

completions: build ## Generate shell completions
	@mkdir -p completions
	./$(BINARY) completion bash > completions/elmos.bash
	./$(BINARY) completion zsh > completions/_elmos
	./$(BINARY) completion fish > completions/elmos.fish
	@echo "Completions generated in completions/"

## Help

help: ## Show this help
	@echo "ELMOS - Embedded Linux on MacOS"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
