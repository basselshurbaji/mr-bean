#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

[[ -f .env ]] && set -a && source .env && set +a

echo "→ Generating sqlc..."
SQLC="$(go env GOPATH)/bin/sqlc"
if [ ! -f "$SQLC" ]; then
    echo "sqlc not found. Run: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
    exit 1
fi
"$SQLC" generate

echo "→ Running linter..."
GOLANGCI_LINT="$(go env GOPATH)/bin/golangci-lint"
if [ ! -f "$GOLANGCI_LINT" ]; then
    echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"
    exit 1
fi
"$GOLANGCI_LINT" run ./...

echo "→ Running go build..."
go build -o bin/server ./cmd/server

echo "✓ built successfully"
