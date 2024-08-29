-- +goose Up
CREATE TABLE pokemon (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    level INT NOT NULL,
    
    ivs_id UUID NOT NULL REFERENCES ivs(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ON pokemon (name, owner_id);
CREATE INDEX ON pokemon (owner_id);

-- +goose Down
DROP TABLE pokemon;