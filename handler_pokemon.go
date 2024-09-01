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
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		OwnerID:   user.ID,
		Name:      p.Name,
		Level:     params.Level,
		Shiny:     params.Shiny,
		IvsID:     ivs.ID,
	})
	if err != nil {
		log.Println("Error creating pokemon: " + err.Error())
		respondWithError(w, 500, "Failed to create pokemon: "+err.Error())
		return
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
		ID:    dbPokemon.ID.String(),
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
