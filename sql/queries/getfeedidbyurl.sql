-- name: GetFeedIDByUrl :one
SELECT id FROM feeds
WHERE url = $1;