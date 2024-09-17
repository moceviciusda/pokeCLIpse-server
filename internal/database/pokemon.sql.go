// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: pokemon.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPokemon = `-- name: CreatePokemon :one
INSERT INTO pokemon (id, created_at, updated_at, name, level, shiny, ivs_id, owner_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at, name, level, shiny, ivs_id, owner_id
`

type CreatePokemonParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Level     int32
	Shiny     bool
	IvsID     uuid.UUID
	OwnerID   uuid.UUID
}

func (q *Queries) CreatePokemon(ctx context.Context, arg CreatePokemonParams) (Pokemon, error) {
	row := q.db.QueryRowContext(ctx, createPokemon,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Level,
		arg.Shiny,
		arg.IvsID,
		arg.OwnerID,
	)
	var i Pokemon
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Level,
		&i.Shiny,
		&i.IvsID,
		&i.OwnerID,
	)
	return i, err
}

const deletePokemon = `-- name: DeletePokemon :one
DELETE FROM pokemon WHERE id = $1 RETURNING id, created_at, updated_at, name, level, shiny, ivs_id, owner_id
`

func (q *Queries) DeletePokemon(ctx context.Context, id uuid.UUID) (Pokemon, error) {
	row := q.db.QueryRowContext(ctx, deletePokemon, id)
	var i Pokemon
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Level,
		&i.Shiny,
		&i.IvsID,
		&i.OwnerID,
	)
	return i, err
}

const getPokemon = `-- name: GetPokemon :one
SELECT id, created_at, updated_at, name, level, shiny, ivs_id, owner_id FROM pokemon WHERE id = $1
`

func (q *Queries) GetPokemon(ctx context.Context, id uuid.UUID) (Pokemon, error) {
	row := q.db.QueryRowContext(ctx, getPokemon, id)
	var i Pokemon
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Level,
		&i.Shiny,
		&i.IvsID,
		&i.OwnerID,
	)
	return i, err
}

const getPokemonPartyByOwnerID = `-- name: GetPokemonPartyByOwnerID :many
SELECT id, created_at, updated_at, name, level, shiny, ivs_id, owner_id FROM pokemon WHERE owner_id = $1 ORDER BY updated_at DESC
`

func (q *Queries) GetPokemonPartyByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]Pokemon, error) {
	rows, err := q.db.QueryContext(ctx, getPokemonPartyByOwnerID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pokemon
	for rows.Next() {
		var i Pokemon
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Level,
			&i.Shiny,
			&i.IvsID,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPokemonWithIvsByOwnerID = `-- name: GetPokemonWithIvsByOwnerID :many
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
WHERE p.owner_id = $1
`

type GetPokemonWithIvsByOwnerIDRow struct {
	ID                uuid.UUID
	Name              string
	Level             int32
	Shiny             bool
	IvsHp             int32
	IvsAttack         int32
	IvsDefense        int32
	IvsSpecialAttack  int32
	IvsSpecialDefense int32
	IvsSpeed          int32
}

func (q *Queries) GetPokemonWithIvsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]GetPokemonWithIvsByOwnerIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getPokemonWithIvsByOwnerID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPokemonWithIvsByOwnerIDRow
	for rows.Next() {
		var i GetPokemonWithIvsByOwnerIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Level,
			&i.Shiny,
			&i.IvsHp,
			&i.IvsAttack,
			&i.IvsDefense,
			&i.IvsSpecialAttack,
			&i.IvsSpecialDefense,
			&i.IvsSpeed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
