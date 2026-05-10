.PHONY: build run clean

build:
	@GOLANGCI_LINT="$$(go env GOPATH)/bin/golangci-lint"; \
	if [ ! -f "$$GOLANGCI_LINT" ]; then \
	    echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"; \
	    exit 1; \
	fi; \
	echo "→ Running linter..."; \
	"$$GOLANGCI_LINT" run ./...
	@echo "→ Running go build..."
	@go build -o mr-bean-mcp .
	@echo "✓ built successfully"

run:
	@go run .

clean:
	@rm -f mr-bean-mcp
