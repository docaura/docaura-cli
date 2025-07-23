.PHONY: build install clean test docs watch-docs lint dev-deps help

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X github.com/docaura/docaura-cli/cmd.version=$(VERSION) -X github.com/docaura/docaura-cli/cmd.commit=$(COMMIT) -X github.com/docaura/docaura-cli/cmd.date=$(DATE)"

# Default target
help: ## Show this help message
	@echo "Docaura CLI - AI-powered Go documentation generator"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the docaura binary
	go build $(LDFLAGS) -o bin/docaura .

install: ## Install docaura to GOPATH/bin
	go install $(LDFLAGS) .

clean: ## Clean build artifacts and documentation
	rm -rf bin/ docs/

test: ## Run all tests
	go test -v ./...

docs: ## Generate documentation for the current project
	./bin/docaura generate -dir . -output ./docs

watch-docs: ## Watch for changes and regenerate documentation
	./bin/docaura generate -dir . -output ./docs --watch

lint: ## Run golangci-lint
	golangci-lint run

# Development helpers
dev-deps: ## Install development dependencies
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Examples for common usage
examples: build ## Show usage examples
	@echo "Docaura CLI Usage Examples:"
	@echo ""
	@echo "1. Initialize a new configuration:"
	@echo "   ./bin/docaura init"
	@echo ""
	@echo "2. Generate documentation:"
	@echo "   ./bin/docaura generate"
	@echo ""
	@echo "3. Generate with watch mode:"
	@echo "   ./bin/docaura generate --watch"
	@echo ""
	@echo "4. Generate HTML documentation:"
	@echo "   ./bin/docaura generate --style html"
	@echo ""
	@echo "5. Show configuration:"
	@echo "   ./bin/docaura config show"
	@echo ""
	@echo "6. Show help:"
	@echo "   ./bin/docaura --help"
	@echo ""
	@echo "For more information, run: ./bin/docaura --help"