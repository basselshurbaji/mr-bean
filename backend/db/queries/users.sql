-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, password_hash, is_active, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, first_name, last_name, email, password_hash, is_active, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id, first_name, last_name, email, password_hash, is_active, created_at, updated_at;

-- name: UpdateUserProfile :one
UPDATE users SET first_name = $1, last_name = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, first_name, last_name, email, password_hash, is_active, created_at, updated_at;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $1, updated_at = NOW()
WHERE id = $2;
