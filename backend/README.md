# backend

Go HTTP API server.

## Stack

| Concern      | Choice                        |
|--------------|-------------------------------|
| Language     | Go                            |
| HTTP router  | chi                           |
| Database     | PostgreSQL & Goose            |
| Query layer  | sqlc (codegen)                |
| Auth         | JWT (access + refresh tokens) |

## Running

```bash
# make sure you're in backend dir
cd backend

# run a clean build
make clean build

# run docker compose for needed dependencies (postgres)
docker compose up -d --build

# run db migrations to create needed tables
make migrate

# run web server (or through an ide for debugging)
make run
```

Environment variables can be set in a `.env` file in `backend/`.

## Environment Variables

| Variable         | Default      | Description                                      |
|------------------|--------------|--------------------------------------------------|
| `PORT`           | `8080`       | HTTP listen port                                 |
| `DB_HOST`        | `localhost`  | Postgres host                                    |
| `DB_PORT`        | `5432`       | Postgres port                                    |
| `DB_USER`        | `postgres`   | Postgres user                                    |
| `DB_PASSWORD`    | _(empty)_    | Postgres password                                |
| `DB_NAME`        | `mr_bean`    | Postgres database name                           |
| `DB_SSLMODE`     | `disable`    | Postgres SSL mode                                |
| `JWT_SECRET`     | `mrbean`     | Signing key for JWTs                             |
| `JWT_EXPIRY`     | `1`          | Access token lifetime in minutes                 |
| `REFRESH_EXPIRY` | `1440`       | Refresh token lifetime in minutes (1440 = 1 day) |

## Database Setup

Start Postgres and apply migrations:

```bash
docker compose up -d   # start Postgres
make migrate           # apply pending migrations
```

Migrations live in `db/migrations/` and are managed by [goose](https://github.com/pressly/goose). Each file follows the naming convention `NNN_description.sql` and must include goose annotations:

```sql
-- +goose Up
-- SQL to apply the migration

-- +goose Down
-- SQL to reverse the migration
```

Never modify an existing migration — write a new one instead.

## Structure

```
backend/
├── cmd/server/     binary entrypoint — wiring only, no logic
├── config/         environment-based configuration
├── db/
│   ├── migrations/ versioned SQL files applied in order (goose)
│   ├── queries/    sqlc input — named queries, one file per domain
│   └── sqlc/       sqlc output — never edit by hand
└── internal/
    ├── handler/    Handler[Req,Res] interface definition only
    ├── router/     chi wiring
    ├── principal/  shared context helper
    ├── auth/       JWT tokens, login/refresh handlers, middleware
    ├── middleware/ Custom middleware registration
    ├── health/     health check endpoint
    └── user/       user repo, service, and handlers
```

## Adding a New Feature

1. Create `internal/<feature>/` with `repo.go`, `service.go`, and one file per handler
2. Add SQL queries to `db/queries/<feature>.sql`, run `make build` (regenerates sqlc output)
3. Wire up in `cmd/server/main.go` using `router.Adapt` and `router.Register` / `router.RegisterProtected`