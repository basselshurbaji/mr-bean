# Mr. Bean

Mono-repo. Backend in Go, frontend TBD.

```
mr_bean/
├── backend/   Go HTTP API
└── frontend/  TBD
```

---

## Backend

### Stack

| Concern     | Choice                        |
|-------------|-------------------------------|
| Language    | Go                            |
| HTTP router | chi                           |
| Database    | PostgreSQL                    |
| Query layer | sqlc (codegen)                |
| Auth        | JWT (access + refresh tokens) |

### Running

```bash
export JWT_SECRET=your-secret-here
cd backend
go run ./cmd/server
```

### Environment Variables

| Variable         | Default      | Description                                      |
|------------------|--------------|--------------------------------------------------|
| `PORT`           | `8080`       | HTTP listen port                                 |
| `DB_HOST`        | `localhost`  | Postgres host                                    |
| `DB_PORT`        | `5432`       | Postgres port                                    |
| `DB_USER`        | `postgres`   | Postgres user                                    |
| `DB_PASSWORD`    | _(empty)_    | Postgres password                                |
| `DB_NAME`        | `mr_bean`    | Postgres database name                           |
| `DB_SSLMODE`     | `disable`    | Postgres SSL mode                                |
| `JWT_SECRET`     | **required** | Signing key for JWTs                             |
| `JWT_EXPIRY`     | `1`          | Access token lifetime in minutes                 |
| `REFRESH_EXPIRY` | `1440`       | Refresh token lifetime in minutes (1440 = 1 day) |

### Database Setup

Apply migrations in order before running the server:

```bash
psql -d mr_bean -f backend/db/migrations/001_create_users.sql
```

---

## Frontend

Stack TBD.

---

## Git Conventions

Commit messages are imperative plain English, short subject line, no period.

```
Add health endpoint
Fix query param decoding for GET requests
Remove unused middleware
```

Every commit made with Claude includes a co-author trailer:

```
Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```
