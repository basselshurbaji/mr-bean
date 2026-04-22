# scripts

Helper scripts invoked by the Makefile.

| Script         | Purpose                                   |
|----------------|-------------------------------------------|
| `build.sh`     | Run `golangci-lint`, then `sqlc generate` |
| `test.sh`      | Run `go test ./...`                       |
| `migrate.sh`   | Run pending migrations via goose          |

Use `make` from the `backend/` directory rather than calling scripts directly.
