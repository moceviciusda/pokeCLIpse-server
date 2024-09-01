-- name: CreatePokemon :one
INSERT INTO pokemon (id, created_at, updated_at, name, level, shiny, ivs_id, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUsersPokemon :many
SELECT * FROM pokemon WHERE owner_id = $1;

-- name: GetPokemon :one
SELECT * FROM pokemon WHERE id = $1;
