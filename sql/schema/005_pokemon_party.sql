-- +goose Up
CREATE TABLE pokemon_party (
    pokemon_id UUID NOT NULL,
    user_id UUID NOT NULL,
    position INT NOT NULL,

    PRIMARY KEY (pokemon_id, user_id),
    FOREIGN KEY (pokemon_id) REFERENCES pokemon(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, position)
);

-- +goose Down
DROP TABLE pokemon_party;