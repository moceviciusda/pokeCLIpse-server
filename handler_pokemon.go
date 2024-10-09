package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

func (cfg *apiConfig) handlerCreatePokemon(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name  string `json:"name"`
		Level int32  `json:"level"`
		Shiny bool   `json:"shiny"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body: " + err.Error())
		respondWithError(w, 400, "Error reading request body: "+err.Error())
		return
	}

	params := parameters{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		log.Println("Error parsing JSON: " + err.Error())
		respondWithError(w, 400, "Error parsing JSON: "+err.Error())
		return
	}

	p, err := cfg.pokeapiClient.GetPokemon(params.Name)
	if err != nil {
		log.Println("Error getting pokemon: " + err.Error())
		respondWithError(w, 400, "Invalid pokemon name: "+err.Error())
		return
	}

	ivs, err := cfg.DB.CreateIVs(r.Context(), database.CreateIVsParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Hp:             int32(rand.Intn(32)),
		Attack:         int32(rand.Intn(32)),
		Defense:        int32(rand.Intn(32)),
		SpecialAttack:  int32(rand.Intn(32)),
		SpecialDefense: int32(rand.Intn(32)),
		Speed:          int32(rand.Intn(32)),
	})
	if err != nil {
		log.Println("Error creating IVs: " + err.Error())
		respondWithError(w, 500, "Failed to create IVs: "+err.Error())
		return
	}

	dbPokemon, err := cfg.DB.CreatePokemon(r.Context(), database.CreatePokemonParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		OwnerID:    user.ID,
		Name:       p.Name,
		Experience: int32(pokeutils.ExpAtLevel(int(params.Level))),
		Level:      params.Level,
		Shiny:      params.Shiny,
		IvsID:      ivs.ID,
	})
	if err != nil {
		log.Println("Error creating pokemon: " + err.Error())
		respondWithError(w, 500, "Failed to create pokemon: "+err.Error())
		return
	}

	rMoves, err := cfg.pokeapiClient.SelectRandomMoves(p.Name, int(params.Level))
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " pokemon moves: " + err.Error())
		respondWithError(w, 500, "Failed to get pokemon moves: "+err.Error())
		return
	}

	for _, move := range rMoves {
		dbMove, err := cfg.DB.GetMoveByName(r.Context(), move.Name)
		if err != nil {
			dbMove, err = cfg.DB.CreateMove(r.Context(), database.CreateMoveParams{
				ID:           uuid.New(),
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
				Name:         move.Name,
				Accuracy:     int32(move.Accuracy),
				Power:        int32(move.Power),
				Pp:           int32(move.Pp),
				Type:         move.Type.Name,
				DamageClass:  move.DamageClass.Name,
				EffectChance: int32(move.EffectChance),
				Effect:       move.EffectEntries[0].ShortEffect,
			})
			if err != nil {
				log.Println("Failed to create move: " + err.Error())
				respondWithError(w, 500, "Failed to create move: "+err.Error())
				return
			}
		}

		_, err = cfg.DB.AddMoveToPokemon(r.Context(), database.AddMoveToPokemonParams{
			PokemonID: dbPokemon.ID,
			MoveName:  dbMove.Name,
		})
		if err != nil {
			log.Println("Failed to create pokemon move: " + err.Error())
			return
		}
	}

	party, err := cfg.DB.GetPokemonParty(r.Context(), user.ID)
	if err != nil {
		log.Println("Error getting pokemon: " + err.Error())
		respondWithError(w, 500, "Failed to get pokemon: "+err.Error())
		return
	}

	if len(party) < 6 {
		_, err := cfg.DB.AddPokemonToParty(r.Context(), database.AddPokemonToPartyParams{
			UserID:    user.ID,
			PokemonID: dbPokemon.ID,
			Position:  int32(len(party) + 1),
		})
		if err != nil {
			log.Println("Failed to add pokemon to party: " + err.Error())
			respondWithError(w, 500, "Failed to add pokemon to party: "+err.Error())
			return
		}
	}

	pBaseStats := pokeutils.Stats{
		Hp:             p.Stats[0].BaseStat,
		Attack:         p.Stats[1].BaseStat,
		Defense:        p.Stats[2].BaseStat,
		SpecialAttack:  p.Stats[3].BaseStat,
		SpecialDefense: p.Stats[4].BaseStat,
		Speed:          p.Stats[5].BaseStat,
	}

	pokemon := pokeutils.Pokemon{
		Name:  dbPokemon.Name,
		Level: int(dbPokemon.Level),
		Shiny: dbPokemon.Shiny,
		Stats: pokeutils.CalculateStats(pBaseStats, pokeutils.IVs{
			Hp:             int(ivs.Hp),
			Attack:         int(ivs.Attack),
			Defense:        int(ivs.Defense),
			SpecialAttack:  int(ivs.SpecialAttack),
			SpecialDefense: int(ivs.SpecialDefense),
			Speed:          int(ivs.Speed),
		}, int(dbPokemon.Level)),
	}

	respondWithJSON(w, 201, pokemon)
}

func (cfg *apiConfig) handlerGetPokemonParty(w http.ResponseWriter, r *http.Request, user database.User) {
	pokemons, err := cfg.DB.GetPokemonParty(r.Context(), user.ID)
	if err != nil {
		log.Println("Error getting pokemon: " + err.Error())
		respondWithError(w, 500, "Failed to get pokemon: "+err.Error())
		return
	}

	pokemonList := make([]pokeutils.Pokemon, 0, len(pokemons))
	for _, pokemon := range pokemons {
		p, err := cfg.pokeapiClient.GetPokemon(pokemon.Name)
		if err != nil {
			log.Println("Error getting pokemon: " + err.Error())
			respondWithError(w, 500, "Failed to get pokemon: "+err.Error())
			return
		}

		ivs, err := cfg.DB.GetIVs(r.Context(), pokemon.IvsID)
		if err != nil {
			log.Println("Error getting IVs: " + err.Error())
			respondWithError(w, 500, "Failed to get IVs: "+err.Error())
			return
		}

		moves, err := cfg.DB.GetMovesByPokemonID(r.Context(), pokemon.ID)
		if err != nil {
			log.Println("Error getting moves: " + err.Error())
			respondWithError(w, 500, "Failed to get moves: "+err.Error())
			return
		}

		pokemonList = append(pokemonList, makePokemon(p, pokemon, moves, ivs))
	}

	respondWithJSON(w, 200, pokemonList)
}
