# scripts

Helper scripts invoked by the Makefile.

| Script      | Purpose                                     |
|-------------|---------------------------------------------|
| `build.sh`  | Run `golangci-lint`, then `sqlc generate`   |
| `test.sh`   | Run `go test ./...`                         |

Use `make` from the `backend/` directory rather than calling scripts directly.
