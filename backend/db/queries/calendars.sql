-- name: CreateCalendar :one
INSERT INTO calendars (user_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCalendarsByUserID :many
SELECT * FROM calendars
WHERE user_id = $1
ORDER BY created_at;

-- name: GetCalendarByID :one
SELECT * FROM calendars
WHERE id = $1;

-- name: UpdateCalendar :one
UPDATE calendars
SET name = $2, description = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCalendar :exec
DELETE FROM calendars
WHERE id = $1;
