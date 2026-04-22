# internal/auth

JWT authentication: token generation/validation, login and refresh handlers, and the auth middleware.

| File            | Purpose                                                              |
|-----------------|----------------------------------------------------------------------|
| `token.go`      | `TokenService` — generate and validate access/refresh JWTs           |
| `service.go`    | `AuthService` — login (bcrypt check) and refresh logic               |
| `login.go`      | `POST /auth/login` handler                                           |
| `refresh.go`    | `POST /auth/refresh` handler                                         |
| `middleware.go` | chi middleware — validates Bearer token, sets user ID in `principal` |
| `errors.go`     | internal validation error helper                                     |

`auth` does not import `internal/user`. It defines its own `UserStore` interface with a minimal `StoredUser` struct. The adapter is in `cmd/server/main.go`.

Future: OAuth2 provider login will be added here.
