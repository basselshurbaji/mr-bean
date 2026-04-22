# internal/user

User domain: database access, business logic, and the `/user/me` endpoint.

| File         | Purpose                                                                |
|--------------|------------------------------------------------------------------------|
| `repo.go`    | `UserRepo` interface, `User` domain struct, sqlc-backed implementation |
| `service.go` | `UserService` interface and implementation                             |
| `me.go`      | `GET /user/me` handler — returns the authenticated user's profile      |

The `User` struct is the canonical user domain type. Handlers in other packages that need a minimal user view (e.g. auth) define their own narrow struct rather than importing this one.

Future handlers for user management (create, update, list) go here.
