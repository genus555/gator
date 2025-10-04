-- name: GetUserFromID :one
SELECT * FROM users
    WHERE users.id = $1;