# internal/principal

Shared context helpers. Currently stores the authenticated user ID.

- `WithUserID(ctx, id)` — set by `auth.Middleware` after token validation
- `UserIDFromContext(ctx)` — read by any handler that needs the caller's identity

This is the only package that puts values into `context.Context`. If a new value needs to be shared via context across packages, add it here rather than creating another context package.
