-- +goose Up
CREATE TABLE users_feeds (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES users
                    ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(user_id, feed_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (feed_id) REFERENCES feeds (id)
);

-- +goose Down
DROP TABLE users_feeds
