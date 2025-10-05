-- name: GetNextFeedToFetch :one
SELECT id, url, name FROM feeds
ORDER BY last_fetched_at NULLS FIRST ,created_at ASC
LIMIT 1;