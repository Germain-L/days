-- name: CreateDayEntry :one
INSERT INTO day_entries (calendar_id, date, color_meaning_id, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDayEntriesByCalendarID :many
SELECT de.*, cm.color_hex, cm.meaning
FROM day_entries de
JOIN color_meanings cm ON de.color_meaning_id = cm.id
WHERE de.calendar_id = $1
ORDER BY de.date DESC;

-- name: GetDayEntriesByDateRange :many
SELECT de.*, cm.color_hex, cm.meaning, c.name as calendar_name
FROM day_entries de
JOIN color_meanings cm ON de.color_meaning_id = cm.id
JOIN calendars c ON de.calendar_id = c.id
WHERE c.user_id = $1 
  AND de.date >= $2 
  AND de.date <= $3
ORDER BY de.date DESC, c.name;

-- name: GetDayEntryByCalendarAndDate :one
SELECT de.*, cm.color_hex, cm.meaning
FROM day_entries de
JOIN color_meanings cm ON de.color_meaning_id = cm.id
WHERE de.calendar_id = $1 AND de.date = $2;

-- name: UpdateDayEntry :one
UPDATE day_entries
SET color_meaning_id = $2, notes = $3, updated_at = NOW()
WHERE calendar_id = $1 AND date = $4
RETURNING *;

-- name: DeleteDayEntry :exec
DELETE FROM day_entries
WHERE calendar_id = $1 AND date = $2;
