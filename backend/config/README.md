# config

Environment-based configuration loaded once at startup via `config.Load()`.

- All values come from environment variables with sensible defaults
- Time durations are integers in minutes (e.g. `JWT_EXPIRY=15`)
- `JWT_SECRET` is required — server fatals if missing
- `Config` is passed explicitly to anything that needs it; no global state

Add new config values here when a package needs an env-configurable setting. Do not read `os.Getenv` anywhere else in the codebase.
