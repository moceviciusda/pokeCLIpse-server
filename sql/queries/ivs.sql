-- name: CreateIVs :one
INSERT INTO ivs (id, created_at, updated_at, hp, attack, defense, special_attack, special_defense, speed)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetIVs :one
SELECT * FROM ivs WHERE id = $1;