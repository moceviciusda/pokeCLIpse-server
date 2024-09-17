-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, password, location_offset)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserLocation :one
UPDATE users
    set location_offset = $2
WHERE id = $1
RETURNING *;

-- name: CheckHasPokemon :one
SELECT * FROM pokemon WHERE owner_id = $1 AND name = $2 AND shiny = $3 LIMIT 1;

-- name: AddPokemonToParty :one
INSERT INTO pokemon_party (pokemon_id, user_id, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPokemonParty :many
SELECT p.*
FROM pokemon p
JOIN pokemon_party pp ON p.id = pp.pokemon_id
WHERE pp.user_id = $1
ORDER BY pp.position;

-- name: GetPokemonInPartyPosition :one
SELECT p.*
FROM pokemon p
JOIN pokemon_party pp ON p.id = pp.pokemon_id
WHERE pp.user_id = $1 AND pp.position = $2;

-- name: RemovePokemonFromParty :exec
DELETE FROM pokemon_party
WHERE pokemon_id = $1 AND user_id = $2;

-- name: GetUserPokemon :many
SELECT * FROM pokemon WHERE owner_id = $1;