.PHONY: build install clean test docs watch-docs

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

build:
	go build $(LDFLAGS) -o bin/docaura .

install:
	go install $(LDFLAGS) .

clean:
	rm -rf bin/ docs/

test:
	go test -v ./...

docs:
	./bin/gendocs -dir . -output ./docs

watch-docs:
	./bin/gendocs -dir . -output ./docs -watch

lint:
	golangci-lint run

# Development helpers
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest