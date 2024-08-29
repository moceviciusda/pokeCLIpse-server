-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,

    location_offset INT NOT NULL DEFAULT 0,

    UNIQUE (username)
);

-- +goose Down
DROP TABLE users;