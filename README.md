# Mr. Bean

Mono-repo. Backend in Go, frontend TBD.

```
mr_bean/
├── backend/   Go HTTP API
└── frontend/  TBD
```

---

## Backend

Chi-based HTTP API with PostgreSQL. Features live in `internal/<feature>/` with handler, service, and repo co-located. Handlers implement a typed `Handler[Req, Res]` interface and are fully decoupled from `net/http`.

See [`backend/README.md`](backend/README.md) for setup, environment variables, and full details.

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
