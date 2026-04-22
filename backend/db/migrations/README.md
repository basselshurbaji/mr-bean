# db/migrations

Plain SQL migration files applied in ascending order.

Naming convention: `NNN_description.sql` (e.g. `001_create_users.sql`).

Rules:
- Never modify an existing migration — write a new one instead
- Each file should be self-contained and idempotent where possible
- Migrations are applied manually; there is no migration runner yet
