package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Username       string    `json:"username"`
	LocationOffset int       `json:"location_offset"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		dbUser.ID,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
		dbUser.Username,
		int(dbUser.LocationOffset),
	}
}

type Pokemon struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	Level     int32     `json:"level"`
	Shiny     bool      `json:"shiny"`
}

// func databasePokemonToPokemon(dbPokemon database.Pokemon) Pokemon {
// 	return Pokemon{
// 		dbPokemon.ID,
// 		dbPokemon.CreatedAt,
// 		dbPokemon.UpdatedAt,
// 		dbPokemon.OwnerID,
// 		dbPokemon.Name,
// 		dbPokemon.Level,
// 		dbPokemon.Shiny,
// 	}
// }

func dbMoveToMove(dbMove database.Move) pokeutils.Move {
	return pokeutils.Move{
		Name:         dbMove.Name,
		Accuracy:     int(dbMove.Accuracy),
		Power:        int(dbMove.Power),
		PP:           int(dbMove.Pp),
		Type:         dbMove.Type,
		DamageClass:  dbMove.DamageClass,
		EffectChance: int(dbMove.EffectChance),
		Effect:       dbMove.Effect,
	}
}
