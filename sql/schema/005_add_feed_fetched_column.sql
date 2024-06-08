-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE users_feeds
DROP COLUMN last_fetched_at;
