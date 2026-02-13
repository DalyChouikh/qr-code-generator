.PHONY: build run clean test lint snapshot release-dry

BINARY   := qrgen
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE     := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

## build: Compile the binary
build:
	@echo "Building $(BINARY) $(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/qrgen

## run: Build and run the application
run: build
	./$(BINARY)

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/

## test: Run tests
test:
	go test -v ./...

## lint: Run linter
lint:
	@command -v golangci-lint > /dev/null 2>&1 || echo "Install golangci-lint: https://golangci-lint.run/welcome/install/"
	golangci-lint run ./...

## snapshot: Build a snapshot release (no publish)
snapshot:
	goreleaser release --snapshot --clean

## release-dry: Dry-run release (validate config)
release-dry:
	goreleaser check
	goreleaser release --skip=publish --clean

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
