-- +goose Up
CREATE TABLE pokemon (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    level INT NOT NULL,
    shiny BOOLEAN NOT NULL,
    
    ivs_id UUID NOT NULL UNIQUE REFERENCES ivs(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    UNIQUE (name, owner_id)
);

-- +goose Down
DROP TABLE pokemon;