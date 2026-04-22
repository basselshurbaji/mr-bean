# db

All database-related files: schema migrations, sqlc query inputs, and generated code.

```
db/
├── migrations/   versioned SQL files applied in order
├── queries/      sqlc input — named queries, one file per domain
└── sqlc/         sqlc output — never edit by hand
```

After changing a migration or query file, regenerate with:

```bash
$(go env GOPATH)/bin/sqlc generate
```
