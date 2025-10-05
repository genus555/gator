-- name: GetFeedFromFeedID :one
SELECT * FROM feeds
    WHERE feeds.id = $1;