# backend

Go HTTP API server.

## Structure

```
backend/
├── cmd/server/     binary entrypoint
├── config/         environment-based configuration
├── db/             SQL migrations, queries, and sqlc-generated code
└── internal/       all application packages (not importable externally)
```

## Adding a new feature

1. Create `internal/<feature>/` with `repo.go`, `service.go`, and one file per handler
2. Add SQL queries to `db/queries/<feature>.sql`, run `sqlc generate`
3. Wire up in `cmd/server/main.go`
