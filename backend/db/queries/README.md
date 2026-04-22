# db/queries

sqlc input files. One `.sql` file per domain area, containing named queries.

Each query uses a sqlc annotation:

```sql
-- name: GetUserByEmail :one
SELECT ...
```

After adding or changing a query, run `$(go env GOPATH)/bin/sqlc generate` to regenerate `db/sqlc/`.

One file per domain (e.g. `users.sql`, `posts.sql`). Do not mix domains in a single file.
