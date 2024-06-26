-- name: CreateFeedFollow :one
INSERT INTO users_feeds(id, user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: AllUserFeedFollows :many
SELECT * FROM users_feeds
WHERE user_id = $1;

-- name: UserFeedFollows :many
SELECT * FROM users_feeds
WHERE user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM users_feeds
WHERE id = $1
AND user_id = $2;
