-- name: CreateFeedFollow :one
WITH insert_feed_follow AS (
    INSERT INTO feed_follows (feed_id, user_id)
    VALUES (
        $1,
        $2
    )
    RETURNING *
)

SELECT 
    insert_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM insert_feed_follow
INNER JOIN feeds ON insert_feed_follow.feed_id = feeds.id
INNER JOIN users ON insert_feed_follow.user_id = users.id;