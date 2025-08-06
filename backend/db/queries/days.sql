-- name: CreateDay :one
INSERT INTO days (user_id, date, title, description, mood)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDay :one
SELECT * FROM days WHERE id = $1;

-- name: GetDayByUserAndDate :one
SELECT * FROM days WHERE user_id = $1 AND date = $2;

-- name: ListDaysByUser :many
SELECT * FROM days
WHERE user_id = $1
ORDER BY date DESC;

-- name: ListDaysByUserAndDateRange :many
SELECT * FROM days
WHERE user_id = $1
  AND date >= $2
  AND date <= $3
ORDER BY date DESC;

-- name: UpdateDay :one
UPDATE days
SET title = $2, description = $3, mood = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDay :exec
DELETE FROM days WHERE id = $1;

-- name: GetUserDaysCount :one
SELECT COUNT(*) FROM days WHERE user_id = $1;
