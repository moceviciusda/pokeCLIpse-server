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