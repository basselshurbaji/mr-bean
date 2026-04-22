# internal

All application packages. Nothing here is importable from outside this module.

Each subdirectory is a feature package containing all layers for that feature (handler, service, repo). Do not create layer-named packages (`handler/`, `service/`, `repo/`) — use feature names instead.

Current packages:

| Package | Purpose |
|---|---|
| `handler/` | `Handler[Req,Res]` interface definition only |
| `router/` | chi wiring — `Adapt`, `Register`, `RegisterProtected`, `NewRouter` |
| `principal/` | shared context helper for the authenticated user ID |
| `auth/` | JWT token service, login/refresh handlers, auth middleware |
| `health/` | health check endpoint |
| `user/` | user repo, service, and `/user/me` handler |

Add new features as new subdirectories here. Each new package gets its own `README.md`.
