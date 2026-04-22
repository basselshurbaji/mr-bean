# internal/router

chi HTTP wiring. All `net/http` and `chi` concerns are contained here — they do not leak into handler implementations.

Exports:
- `NewRouter()` — creates the base chi router with standard middleware
- `Adapt(h Handler[Req,Res]) Route` — wraps a handler into an opaque `Route`
- `Register(r, route)` — mounts a public route
- `RegisterProtected(r, middleware, routes...)` — mounts routes under a middleware group

Request decoding:
- `POST`, `PUT`, `PATCH` → JSON body
- `GET` → query params via `gorilla/schema` (`schema` struct tags)
