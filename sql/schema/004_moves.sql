-- +goose Up
CREATE TABLE moves (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    accuracy INT NOT NULL,
    power INT NOT NULL,
    pp INT NOT NULL,
    type TEXT NOT NULL,
    damage_class TEXT NOT NULL,
    effect_chance INT NOT NULL,
    effect TEXT NOT NULL
);

CREATE TABLE moves_pokemon (
    move_name TEXT NOT NULL,
    pokemon_id UUID NOT NULL,
    PRIMARY KEY (move_name, pokemon_id),
    FOREIGN KEY (move_name) REFERENCES moves(name) ON DELETE CASCADE,
    FOREIGN KEY (pokemon_id) REFERENCES pokemon(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE moves_pokemon;
DROP TABLE moves;
