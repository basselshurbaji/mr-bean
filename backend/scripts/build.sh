#!/usr/bin/env zsh
set -euo pipefail

cd "$(dirname "$0")/.."

[[ -f .env ]] && set -a && source .env && set +a

echo "→ Generating sqlc..."s
SQLC="$(go env GOPATH)/bin/sqlc"
if [ ! -f "$SQLC" ]; then
    echo "sqlc not found. Run: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
    exit 1
fi
"$SQLC" generate

echo "→ Running linter..."
if ! command -v golangci-lint &>/dev/null; then
    echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"
    exit 1
fi
golangci-lint run ./...

echo "✓ built successfully"
