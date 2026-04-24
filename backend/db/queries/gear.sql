-- Gear

-- name: ListGearByUserID :many
SELECT id, user_id, type_id, name, brand, model, year, notes, created_at, updated_at
FROM gear
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetGearByID :one
SELECT id, user_id, type_id, name, brand, model, year, notes, created_at, updated_at
FROM gear
WHERE id = $1 AND user_id = $2;

-- name: CreateGear :one
INSERT INTO gear (user_id, type_id, name, brand, model, year, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, type_id, name, brand, model, year, notes, created_at, updated_at;

-- name: UpdateGear :one
UPDATE gear
SET type_id = $1, name = $2, brand = $3, model = $4, year = $5, notes = $6, updated_at = NOW()
WHERE id = $7 AND user_id = $8
RETURNING id, user_id, type_id, name, brand, model, year, notes, created_at, updated_at;

-- name: GetStationIDsByGearID :many
SELECT station_id FROM station_gear WHERE gear_id = $1;

-- name: ListGearIDsInStationExcluding :many
SELECT gear_id FROM station_gear
WHERE station_id = $1 AND gear_id != $2
ORDER BY position ASC;

-- name: DeleteStationGearByStationID :exec
DELETE FROM station_gear WHERE station_id = $1;

-- name: DeleteGearByID :execrows
DELETE FROM gear WHERE id = $1 AND user_id = $2;

-- Stations

-- name: ListStationsByUserID :many
SELECT id, user_id, name, created_at, updated_at
FROM stations
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetStationByID :one
SELECT id, user_id, name, created_at, updated_at
FROM stations
WHERE id = $1 AND user_id = $2;

-- name: CreateStation :one
INSERT INTO stations (user_id, name)
VALUES ($1, $2)
RETURNING id, user_id, name, created_at, updated_at;

-- name: UpdateStation :one
UPDATE stations
SET name = $1, updated_at = NOW()
WHERE id = $2 AND user_id = $3
RETURNING id, user_id, name, created_at, updated_at;

-- name: DeleteStationByID :execrows
DELETE FROM stations WHERE id = $1 AND user_id = $2;

-- name: ListStationGearByUserID :many
SELECT sg.station_id,
       g.id, g.user_id, g.type_id, g.name, g.brand, g.model, g.year, g.notes,
       g.created_at, g.updated_at
FROM station_gear sg
JOIN gear g ON g.id = sg.gear_id
WHERE g.user_id = $1
ORDER BY sg.station_id, sg.position ASC;

-- name: InsertStationGear :exec
INSERT INTO station_gear (station_id, gear_id, position)
VALUES ($1, $2, $3);
