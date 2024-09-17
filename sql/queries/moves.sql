-- name: CreateMove :one
INSERT INTO moves (id, created_at, updated_at, name, accuracy, power, pp, type, damage_class, effect_chance, effect)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetMoveByNameOrId :one
SELECT * FROM moves WHERE name = $1 OR id = $1;

-- name: GetMovesByPokemonID :many
SELECT m.*
FROM moves m
JOIN moves_pokemon mp ON m.name = mp.move_name
WHERE mp.pokemon_id = $1;

-- name: AddMoveToPokemon :one
INSERT INTO moves_pokemon (move_name, pokemon_id)
VALUES ($1, $2)
RETURNING *;