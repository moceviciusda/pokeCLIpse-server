package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/internal/pokeapi"
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

func dbIvsToIvs(dbIvs database.Iv) pokeutils.IVs {
	return pokeutils.IVs{
		Hp:             int(dbIvs.Hp),
		Attack:         int(dbIvs.Attack),
		Defense:        int(dbIvs.Defense),
		SpecialAttack:  int(dbIvs.SpecialAttack),
		SpecialDefense: int(dbIvs.SpecialDefense),
		Speed:          int(dbIvs.Speed),
	}
}

func makePokemon(p pokeapi.PokemonResponse, dbPokemon database.Pokemon, moves []database.Move, dbIvs database.Iv) pokeutils.Pokemon {
	pokemonMoves := make([]pokeutils.Move, len(moves))
	for i, move := range moves {
		pokemonMoves[i] = dbMoveToMove(move)
	}

	types := make([]string, len(p.Types))
	for i, t := range p.Types {
		types[i] = t.Type.Name
	}

	ivs := dbIvsToIvs(dbIvs)

	baseStats := pokeutils.Stats{
		Hp:             p.Stats[0].BaseStat,
		Attack:         p.Stats[1].BaseStat,
		Defense:        p.Stats[2].BaseStat,
		SpecialAttack:  p.Stats[3].BaseStat,
		SpecialDefense: p.Stats[4].BaseStat,
		Speed:          p.Stats[5].BaseStat,
	}

	return pokeutils.Pokemon{
		Name:  p.Name,
		Types: types,
		Level: int(dbPokemon.Level),
		Shiny: dbPokemon.Shiny,
		Stats: pokeutils.CalculateStats(baseStats, ivs, int(dbPokemon.Level)),
		Moves: pokemonMoves,
	}
}
