-- name: PostsByFeedID :many
SELECT * FROM posts
WHERE feed_id = $1
ORDER BY published_at ASC
LIMIT $2;