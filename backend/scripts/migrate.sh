#!/usr/bin/env zsh
set -euo pipefail

cd "$(dirname "$0")/.."

GOOSE="$(go env GOPATH)/bin/goose"
if [ ! -f "$GOOSE" ]; then
    echo "goose not found. Run: go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

DSN="host=${DB_HOST:-localhost} port=${DB_PORT:-5432} user=${DB_USER:-postgres} password=${DB_PASSWORD:-postgres} dbname=${DB_NAME:-mr_bean} sslmode=${DB_SSLMODE:-disable}"

"$GOOSE" -dir db/migrations postgres "$DSN" up
