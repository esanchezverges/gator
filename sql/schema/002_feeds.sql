-- +goose Up
CREATE TABLE feeds (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	url TEXT NOT NULL UNIQUE,
	user_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	CONSTRAINT fk_feeds_user
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE
);

CREATE INDEX idx_feeds_user_id ON feeds(user_id);

-- +goose Down
DROP TABLE IF EXISTS feeds;

