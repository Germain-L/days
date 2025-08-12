-- name: CreateColorMeaning :one
INSERT INTO color_meanings (calendar_id, color_hex, meaning)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetColorMeaningsByCalendarID :many
SELECT * FROM color_meanings
WHERE calendar_id = $1
ORDER BY created_at;

-- name: GetColorMeaningByID :one
SELECT * FROM color_meanings
WHERE id = $1;

-- name: UpdateColorMeaning :one
UPDATE color_meanings
SET color_hex = $2, meaning = $3
WHERE id = $1
RETURNING *;

-- name: DeleteColorMeaning :exec
DELETE FROM color_meanings
WHERE id = $1;
