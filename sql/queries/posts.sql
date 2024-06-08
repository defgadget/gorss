-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_date, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
JOIN users_feeds
ON posts.feed_id = users_feeds.feed_id
JOIN users
ON users_feeds.user_id = users.id
WHERE users_feeds.user_id = $1
ORDER BY posts.published_date DESC
LIMIT $2;
