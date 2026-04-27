-- name: ListExtractionsByUserID :many
SELECT
    e.id, e.user_id, e.bean_id,
    e.dose_in, e.yield_out, e.time, e.target_time, e.grind_size,
    e.pre_infusion, e.tasting_note, e.created_at, e.updated_at,
    b.name AS bean_name, b.roaster AS bean_roaster, b.roast_level AS bean_roast
FROM extractions e
JOIN beans b ON b.id = e.bean_id
WHERE e.user_id = $1
ORDER BY e.created_at DESC
LIMIT $2::bigint OFFSET $3::bigint;

-- name: GetExtractionByID :one
SELECT
    e.id, e.user_id, e.bean_id,
    e.dose_in, e.yield_out, e.time, e.target_time, e.grind_size,
    e.pre_infusion, e.tasting_note, e.created_at, e.updated_at,
    b.name AS bean_name, b.roaster AS bean_roaster, b.roast_level AS bean_roast
FROM extractions e
JOIN beans b ON b.id = e.bean_id
WHERE e.id = $1 AND e.user_id = $2;

-- name: CreateExtraction :one
INSERT INTO extractions (user_id, bean_id, dose_in, yield_out, time, target_time, grind_size, pre_infusion, tasting_note)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, bean_id, dose_in, yield_out, time, target_time, grind_size, pre_infusion, tasting_note, created_at, updated_at;

-- name: UpdateExtraction :one
UPDATE extractions
SET bean_id = $1, dose_in = $2, yield_out = $3, time = $4, target_time = $5,
    grind_size = $6, pre_infusion = $7, tasting_note = $8, updated_at = NOW()
WHERE id = $9 AND user_id = $10
RETURNING id, user_id, bean_id, dose_in, yield_out, time, target_time, grind_size, pre_infusion, tasting_note, created_at, updated_at;

-- name: DeleteExtractionByID :execrows
DELETE FROM extractions WHERE id = $1 AND user_id = $2;

-- name: GetExtractionGear :many
SELECT g.id, g.type_id, g.name
FROM extraction_gear eg
JOIN gear g ON g.id = eg.gear_id
WHERE eg.extraction_id = $1;

-- name: ListExtractionGearByUserID :many
SELECT eg.extraction_id, g.id, g.type_id, g.name
FROM extraction_gear eg
JOIN gear g ON g.id = eg.gear_id
JOIN extractions e ON e.id = eg.extraction_id
WHERE e.user_id = $1;

-- name: InsertExtractionGear :exec
INSERT INTO extraction_gear (extraction_id, gear_id)
VALUES ($1, $2);

-- name: DeleteExtractionGearByExtractionID :exec
DELETE FROM extraction_gear WHERE extraction_id = $1;
