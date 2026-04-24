#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "→ cleaning go binaries..."
rm -rf bin/

echo "→ cleaning sqlc generated files..."
find db/sqlc -name "*.go" -maxdepth 1 -delete


echo "✓ cleaned successfully"