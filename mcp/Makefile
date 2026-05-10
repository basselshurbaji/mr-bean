.PHONY: build run clean

build:
	@if ! command -v golangci-lint &>/dev/null; then \
	    echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"; \
	    exit 1; \
	fi
	@echo "→ Running linter..."
	@golangci-lint run ./...
	@echo "→ Running go build..."
	@go build -o mr-bean-mcp .
	@echo "✓ built successfully"

run:
	@go run .

clean:
	@rm -f mr-bean-mcp
