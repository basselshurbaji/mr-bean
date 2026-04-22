-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, password_hash, is_active, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, first_name, last_name, email, password_hash, is_active, created_at, updated_at
FROM users
WHERE id = $1;
