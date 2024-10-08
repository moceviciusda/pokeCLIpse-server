// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Iv struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Hp             int32
	Attack         int32
	Defense        int32
	SpecialAttack  int32
	SpecialDefense int32
	Speed          int32
}

type Move struct {
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

type MovesPokemon struct {
	MoveName  string
	PokemonID uuid.UUID
}

type Pokemon struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	Experience int32
	Level      int32
	Shiny      bool
	IvsID      uuid.UUID
	OwnerID    uuid.UUID
}

type PokemonParty struct {
	PokemonID uuid.UUID
	UserID    uuid.UUID
	Position  int32
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Username       string
	Password       string
	LocationOffset int32
}
