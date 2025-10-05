-- name: GetFeedFollowerForUser :many
SELECT * from feed_follows
    WHERE feed_follows.user_id = $1;