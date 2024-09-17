// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: moves.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addMoveToPokemon = `-- name: AddMoveToPokemon :one
INSERT INTO moves_pokemon (move_name, pokemon_id)
VALUES ($1, $2)
RETURNING move_name, pokemon_id
`

type AddMoveToPokemonParams struct {
	MoveName  string
	PokemonID uuid.UUID
}

func (q *Queries) AddMoveToPokemon(ctx context.Context, arg AddMoveToPokemonParams) (MovesPokemon, error) {
	row := q.db.QueryRowContext(ctx, addMoveToPokemon, arg.MoveName, arg.PokemonID)
	var i MovesPokemon
	err := row.Scan(&i.MoveName, &i.PokemonID)
	return i, err
}

const createMove = `-- name: CreateMove :one
INSERT INTO moves (id, created_at, updated_at, name, accuracy, power, pp, type, damage_class, effect_chance, effect)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, created_at, updated_at, name, accuracy, power, pp, type, damage_class, effect_chance, effect
`

type CreateMoveParams struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Accuracy     int32
	Power        int32
	Pp           int32
	Type         string
	DamageClass  string
	EffectChance int32
	Effect       string
}

func (q *Queries) CreateMove(ctx context.Context, arg CreateMoveParams) (Move, error) {
	row := q.db.QueryRowContext(ctx, createMove,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Accuracy,
		arg.Power,
		arg.Pp,
		arg.Type,
		arg.DamageClass,
		arg.EffectChance,
		arg.Effect,
	)
	var i Move
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Accuracy,
		&i.Power,
		&i.Pp,
		&i.Type,
		&i.DamageClass,
		&i.EffectChance,
		&i.Effect,
	)
	return i, err
}

const getMoveByName = `-- name: GetMoveByName :one
SELECT id, created_at, updated_at, name, accuracy, power, pp, type, damage_class, effect_chance, effect FROM moves WHERE name = $1
`

func (q *Queries) GetMoveByName(ctx context.Context, name string) (Move, error) {
	row := q.db.QueryRowContext(ctx, getMoveByName, name)
	var i Move
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Accuracy,
		&i.Power,
		&i.Pp,
		&i.Type,
		&i.DamageClass,
		&i.EffectChance,
		&i.Effect,
	)
	return i, err
}

const getMovesByPokemonID = `-- name: GetMovesByPokemonID :many
SELECT m.id, m.created_at, m.updated_at, m.name, m.accuracy, m.power, m.pp, m.type, m.damage_class, m.effect_chance, m.effect
FROM moves m
JOIN moves_pokemon mp ON m.name = mp.move_name
WHERE mp.pokemon_id = $1
`

func (q *Queries) GetMovesByPokemonID(ctx context.Context, pokemonID uuid.UUID) ([]Move, error) {
	rows, err := q.db.QueryContext(ctx, getMovesByPokemonID, pokemonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Move
	for rows.Next() {
		var i Move
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Accuracy,
			&i.Power,
			&i.Pp,
			&i.Type,
			&i.DamageClass,
			&i.EffectChance,
			&i.Effect,
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