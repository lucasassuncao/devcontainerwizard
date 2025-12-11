.PHONY: help build release install docs clean test test-coverage fmt lint deps run run-with-config generate-docs show-docs install-tools fix-permissions

# Tool versions
GORELEASER_VERSION := v2@latest
GOLANGCI_LINT_VERSION := v2.5.0
GOMARKDOC_VERSION := latest

# Project variables
BINARY_NAME := devcontainer
BUILD_DIR := bin
MAIN_PATH := main.go
TOOLS_DIR := .binaries

# Export GOBIN so go install uses our local tools directory
export GOBIN := $(shell pwd)/$(TOOLS_DIR)

# Tool binaries
GORELEASER := $(TOOLS_DIR)/goreleaser
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint
GOMARKDOC := $(TOOLS_DIR)/gomarkdoc

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Install goreleaser to local tools directory
$(GORELEASER):
	@echo "Installing goreleaser to $(TOOLS_DIR)..."
	@mkdir -p $(TOOLS_DIR)
	@go install github.com/goreleaser/goreleaser/$(GORELEASER_VERSION)

# Install golangci-lint to local tools directory
$(GOLANGCI_LINT):
	@echo "Installing golangci-lint to $(TOOLS_DIR)..."
	@mkdir -p $(TOOLS_DIR)
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# Install gomarkdoc to local tools directory
$(GOMARKDOC):
	@echo "Installing gomarkdoc to $(TOOLS_DIR)..."
	@mkdir -p $(TOOLS_DIR)
	@go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@$(GOMARKDOC_VERSION)

install-tools: $(GORELEASER) $(GOLANGCI_LINT) $(GOMARKDOC) ## Install all tools to .binaries directory

build: $(GORELEASER) ## Build binary with goreleaser (current platform only)
	@echo "Building..."
	@$(GORELEASER) build --skip=validate --single-target --snapshot --clean

build-all: $(GORELEASER) ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@$(GORELEASER) build --skip=validate --snapshot --clean

release: $(GORELEASER) ## Create a release with goreleaser
	@echo "Creating release..."
	@$(GORELEASER) release --timeout 360s

install: ## Install binary globally
	@go install

fmt: ## Format code
	@go fmt ./...

lint: $(GOLANGCI_LINT) ## Run linter checks
	@$(GOLANGCI_LINT) -v run ./...

test: ## Run tests
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

deps: ## Download and tidy dependencies
	@go mod download
	@go mod tidy

fix-permissions: ## Fix file permissions (useful in WSL)
	@echo "Fixing file permissions..."
	@sudo chown -R $(USER):$(USER) . 2>/dev/null || true

docs: $(GOMARKDOC) ## Generate documentation with gomarkdoc
	@sudo $(GOMARKDOC) -e -o '{{.Dir}}/README.md' ./...

generate-docs: ## Generate docs using the app
	@go run $(MAIN_PATH) generate-docs

show-docs: ## Show documentation in terminal
	@go run $(MAIN_PATH) show-docs

run: ## Run the application
	@go run $(MAIN_PATH)

run-with-config: ## Run with custom config (CONFIG=path OUTPUT=path)
	@go run $(MAIN_PATH) -c $(CONFIG) -o $(OUTPUT)

clean: ## Remove build artifacts and cache
	@rm -rf $(BUILD_DIR) dist/ coverage.out $(TOOLS_DIR)
	@go clean -cache -testcache