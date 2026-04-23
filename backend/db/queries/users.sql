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
