-- name: GetPosts :many
SELECT * FROM posts
ORDER BY published_at ASC
LIMIT $1;