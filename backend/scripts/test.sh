#!/usr/bin/env zsh
set -euo pipefail

cd "$(dirname "$0")/.."

[[ -f .env ]] && set -a && source .env && set +a

echo "→ Running tests..."
go test ./...
