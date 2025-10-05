-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1
WHERE feeds.id = $2;