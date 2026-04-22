# db/migrations

Plain SQL migration files managed by [goose](https://github.com/pressly/goose).

Naming convention: `NNN_description.sql` (e.g. `001_create_users.sql`).

Each file must begin with a goose annotation:

```sql
-- +goose Up
-- SQL to apply the migration

-- +goose Down
-- SQL to reverse the migration
```

Rules:
- Never modify an existing migration — write a new one instead
- goose tracks applied migrations in a `goose_db_version` table
- Run `make migrate` from `backend/` to apply pending migrations
