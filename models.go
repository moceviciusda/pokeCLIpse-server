package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
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
