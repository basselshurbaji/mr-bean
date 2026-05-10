#!/bin/sh
set -e

echo "==> Running database migrations..."
goose -dir ./migrations postgres \
  "host=${DB_HOST:-postgres} port=${DB_PORT:-5432} user=${DB_USER:-postgres} password=${DB_PASSWORD:-postgres} dbname=${DB_NAME:-mr_bean} sslmode=${DB_SSLMODE:-disable}" \
  up

echo "==> Starting mr-bean backend on :${PORT:-8080}..."
exec ./server