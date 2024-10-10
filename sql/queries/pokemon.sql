-- name: CreatePokemon :one
INSERT INTO pokemon (id, created_at, updated_at, name, experience, level, shiny, ivs_id, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdatePokemonLvlAndExp :one
UPDATE pokemon SET level = $1, experience = $2 WHERE id = $3 RETURNING *;

-- name: UpdatePokemonName :one
UPDATE pokemon SET name = $1 WHERE id = $2 RETURNING *;

-- name: GetPokemonWithIvsByOwnerID :many
SELECT 
    p.id,
    p.name,
    p.level,
    p.shiny,
    i.hp AS ivs_hp,
    i.attack AS ivs_attack,
    i.defense AS ivs_defense,
    i.special_attack AS ivs_special_attack,
    i.special_defense AS ivs_special_defense,
    i.speed AS ivs_speed
FROM pokemon p
JOIN ivs i ON p.ivs_id = i.id
WHERE p.owner_id = $1;

-- name: GetPokemon :one
SELECT * FROM pokemon WHERE id = $1;

-- name: DeletePokemon :one
DELETE FROM pokemon WHERE id = $1 RETURNING *;