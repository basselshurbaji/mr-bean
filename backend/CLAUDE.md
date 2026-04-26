# Mr. Bean — Claude Instructions

This file is for AI assistants working on this repo. It describes architecture rules, coding conventions, and patterns to follow. For project setup and stack overview, see `backend/README.md`.

> Always format Markdown tables with aligned columns — pad cells with spaces so all `|` separators line up vertically. Apply this to every table in every file, including `CLAUDE.md` itself.

---

## Project Layout

```
mr_bean/
├── backend/
│   ├── cmd/server/         # binary entrypoint — wiring only, no logic
│   ├── config/             # env-based config, loaded once at startup
│   ├── db/
│   │   ├── migrations/     # plain SQL files, applied in order
│   │   ├── queries/        # sqlc input — named .sql files, one per domain
│   │   └── sqlc/           # sqlc output — never edit by hand
│   └── internal/
│       ├── handler/        # Handler[Req,Res] interface + AppError + NoContent
│       ├── router/         # chi wiring: Router struct, NewRouter, Register[Req,Res]
│       ├── middleware/     # Tag-based middleware registry (Register, Resolve, TagAuthenticated)
│       ├── principal/      # shared context helper for user ID (no auth logic)
│       ├── auth/           # JWT tokens, login/refresh/register handlers, middleware
│       ├── health/         # health check endpoint
│       ├── user/           # user repo, service, and /user/me /user/update /user/change-password handlers
│       ├── gear/           # gear + station repo, service, and CRUD handlers
│       ├── bean/           # bean repo, service, and CRUD handlers
│       └── mailer/         # Mailer interface + SMTP implementation
└── frontend/               # TBD
```

`internal/` follows standard Go convention — nothing inside is importable from outside this module.

---

## Architecture

### Folder-by-feature

Each feature lives in its own package under `internal/`. All three layers (handler, service, repo) are co-located in that folder:

```
internal/auth/
    login.go      ← handler
    refresh.go    ← handler
    service.go    ← business logic
    token.go      ← token service
    middleware.go ← chi middleware

internal/user/
    me.go         ← handler
    service.go    ← business logic
    repo.go       ← database access + domain struct
```

Do not put handlers in `internal/handler/`, services in `internal/service/`, etc. That pattern is not used.

### Request flow

```
Handler → Service → Repo
```

Each layer depends on the **interface** of the layer below, never the concrete type. No layer imports another layer's package directly — they communicate through interfaces and domain structs defined locally or in `internal/principal`.

---

## Go Conventions

### Interface satisfaction

When a method exists solely to satisfy an interface, document it with a standard comment:

```go
// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) { ... }
```

Format: `// MethodName implements InterfaceName.`

---

## Handler Contract

```go
// internal/handler/handler.go
type Handler[Req, Res any] interface {
    Method() string
    Pattern() string
    Middlewares() []middleware.Tag        // middleware applied to this route
    Validate(req Req) error              // 422 before Serve if non-nil
    Serve(ctx context.Context, req Req) (Res, error)
}

// Return an *AppError from Serve to send a specific HTTP status instead of 500.
type AppError struct {
    Code int
    Msg  string
}

// Return NoContent from Serve to send HTTP 204 with no body.
type NoContent struct{}
```

- `context.Context` is always the first argument to `Serve` — it carries the request context and any middleware-injected values (e.g. user ID via `principal`).
- No HTTP types (`http.Request`, `http.ResponseWriter`) appear in handler implementations.
- `chi` and `net/http` are fully contained in `internal/router/router.go`.
- GET request structs use `schema` tags; body request structs use `json` tags; chi URL path params use `url` tags.
- DELETE handlers receive no body.

### Middleware tags

Handlers declare which middleware they need by returning `[]middleware.Tag` from `Middlewares()`. The router resolves tags to concrete middleware at startup — unregistered tags panic immediately.

Currently defined tags:

| Tag                            | Effect                                           |
|--------------------------------|--------------------------------------------------|
| `middleware.TagAuthenticated`  | Validates Bearer token, sets user ID in context  |

Register middleware once at startup before constructing the router:

```go
middleware.Register(middleware.TagAuthenticated, auth.Middleware(tokenSvc))
```

### Registration

`router.Register` is a generic free function — it takes a `*Router` and any `handler.Handler[Req, Res]`. No intermediate `Route` type or `Adapt` call is needed.

All routes use the same flat registration — no separate "protected" group:

```go
r := router.NewRouter()

// public — Middlewares() returns nil
router.Register(r, auth.NewLoginHandler(authSvc))

// authenticated — Middlewares() returns []middleware.Tag{middleware.TagAuthenticated}
router.Register(r, user.NewMeHandler(userSvc))
```

`*Router` implements `http.Handler`, so pass it directly to `http.ListenAndServe`.

Adding an endpoint = one `router.Register(r, ...)` call. No other wiring required.

---

## Context Values

`internal/principal` is the only package that puts values into `context.Context`. Currently it stores the authenticated user ID.

- Set by: `auth.Middleware` via `principal.WithUserID`
- Read by: handlers via `principal.UserIDFromContext`

Do not create additional context-value packages. If a new value needs to be shared via context, add it to `principal`.

---

## Build & Linting

Run `make build` from `backend/` before committing. It runs `golangci-lint` then `sqlc generate`. The linter **must pass** — do not commit code with lint errors.

- Install golangci-lint: https://golangci-lint.run/usage/install/
- Install sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- Install goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`

`make migrate` applies pending migrations (goose tracks state in `goose_db_version`). `make test` runs all tests. `make clean` removes the compiled binary.

---

## Database

No ORM. SQL is written explicitly.

- Write named queries in `db/queries/<feature>.sql` using sqlc annotations
- Run `$(go env GOPATH)/bin/sqlc generate` after changing queries or schema
- Never edit files in `db/sqlc/` by hand
- `sqlc.yaml` has type overrides: `uuid → string`, `timestamptz → time.Time`
- Repo interfaces are defined before implementations, enabling testing without a live database

---

## Configuration

All values come from environment variables. Time durations are expressed as **integer minutes** in env vars (e.g. `JWT_EXPIRY=15` means 15 minutes). Parsed by `getMinutes(key, fallbackMinutes)` in `config.go`.

| Env var          | Default | Unit    |
|------------------|---------|---------|
| `JWT_EXPIRY`     | `1`     | minutes |
| `REFRESH_EXPIRY` | `1440`  | minutes |

`JWT_SECRET` defaults to `mrbean` if not set.

---

## Authentication

JWT-only. No OAuth2 yet (planned).

- **Access tokens**: short-lived (default 1 min), HS256-signed JWT with `typ:"access"` claim
- **Refresh tokens**: long-lived (default 1440 min), same signing, `typ:"refresh"` claim
- `ValidateAccessToken` and `ValidateRefreshToken` each reject the wrong token type
- Auth middleware validates the Bearer token and sets user ID in context via `principal`
- Individual handlers never perform auth checks — that belongs in middleware

---

## Testing

**Rule:** unit test all logic and handlers; integration tests (real Postgres via testcontainers) are the path for the repo layer when needed. Never skip testing entirely.

### Unit tests

- Test files use the external test package (`package foo_test`) — tests access only exported surface
- Call handler methods (`Validate`, `Serve`) directly — do not go through the HTTP layer
- Use `context.Background()` for tests that don't need auth; use `principal.WithUserID(ctx, id)` for tests that do
- For expired-token tests: instantiate `TokenService` with a negative expiry (e.g. `-time.Second`)

### Mock pattern

Define a small private struct in the test file that satisfies the interface. No mocking libraries.

```go
type mockUserStore struct {
    user *auth.StoredUser
    err  error
}

func (m *mockUserStore) GetByEmail(_ context.Context, _ string) (*auth.StoredUser, error) {
    return m.user, m.err
}
```

All mocks live in the test file that uses them — never in a shared mock package.

### Integration tests (future)

When repo-layer tests are needed: use `testcontainers-go` to spin up a real `postgres:16-alpine` container, apply migrations from `db/migrations/`, seed data, and test against it. A shared helper `internal/testutil/db.go` will own container setup and `t.Cleanup` teardown.

---

## Git Conventions

Imperative plain English. Short subject line, no period.

```
Add health endpoint
Fix query param decoding for GET requests
Remove unused middleware
```

Every commit made with Claude ends with a co-author trailer:

```
Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

The model name in the trailer reflects the model used in that session.

### What to stage

Stage source files explicitly by name — never `git add .` or `git add -A`.

**Never stage `db/sqlc/*.go`** — these files are gitignored and regenerated by `make build`. Attempting to add them will fail with an "ignored by .gitignore" error. Only stage the inputs that drive generation:

| Generated by    | Stage these instead         |
|-----------------|-----------------------------|
| `sqlc generate` | `db/queries/*.sql`          |
| `sqlc generate` | `db/migrations/*.sql`       |
