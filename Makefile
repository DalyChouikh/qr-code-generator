.PHONY: build run clean test lint snapshot release-dry wasm website website-dev

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

## wasm: Compile the generator to WebAssembly for the website
wasm:
	mkdir -p website/public
	GOOS=js GOARCH=wasm go build -ldflags "-s -w" -o website/public/qrgen.wasm ./cmd/qrgen-wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" website/public/wasm_exec.js

## website: Build the static website into website/out
website: wasm
	cd website && npm ci && npm run build

## website-dev: Run the website dev server
website-dev: wasm
	cd website && npm run dev

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/
	rm -rf website/out website/.next
	rm -f website/public/qrgen.wasm website/public/wasm_exec.js

## test: Run tests
test:
	go test -v ./...

## lint: Run go vet and golangci-lint
lint:
	@echo "Running go vet..."
	@go vet ./...
	@if command -v golangci-lint > /dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping (install: https://golangci-lint.run)"; \
	fi

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
