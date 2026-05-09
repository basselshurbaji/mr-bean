-- name: CreateAppToken :one
INSERT INTO app_tokens (user_id, app_name)
VALUES ($1, $2)
RETURNING id, user_id, app_name, revoked, created_at;

-- name: GetAppTokenByID :one
SELECT id, user_id, app_name, revoked, created_at
FROM app_tokens
WHERE id = $1;

-- name: RevokeAppToken :exec
UPDATE app_tokens
SET revoked = TRUE
WHERE id = $1 AND user_id = $2;