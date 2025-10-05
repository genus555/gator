-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT user_feed_pair UNIQUE (user_id, feed_id)
);

CREATE TRIGGER set_feed_follows_updated_at
BEFORE UPDATE ON feed_follows
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS set_feed_follows_updated_at ON feed_follows;
DROP TABLE IF EXISTS feed_follows;