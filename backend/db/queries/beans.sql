-- name: ListBeansByUserID :many
SELECT id, user_id, name, roaster, origin, process, roast_level, tasting_notes, notes, created_at, updated_at
FROM beans
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: CreateBean :one
INSERT INTO beans (user_id, name, roaster, origin, process, roast_level, tasting_notes, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, name, roaster, origin, process, roast_level, tasting_notes, notes, created_at, updated_at;

-- name: UpdateBean :one
UPDATE beans
SET name = $1, roaster = $2, origin = $3, process = $4, roast_level = $5,
    tasting_notes = $6, notes = $7, updated_at = NOW()
WHERE id = $8 AND user_id = $9
RETURNING id, user_id, name, roaster, origin, process, roast_level, tasting_notes, notes, created_at, updated_at;

-- name: DeleteBeanByID :execrows
DELETE FROM beans WHERE id = $1 AND user_id = $2;
