package main

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
)

func (cfg *apiConfig) hadlerGetUserLocation(w http.ResponseWriter, r *http.Request, user database.User) {
	location, err := cfg.pokeapiClient.GetLocationArea(strconv.Itoa(int(user.LocationID)))
	if err != nil {
		respondWithError(w, 500, "Failed to get user location: "+err.Error())
		return
	}

	type response struct {
		Name string `json:"name"`
	}

	respondWithJSON(w, 200, response{Name: location.Name})
}

func (cfg *apiConfig) handlerSearchForPokemon(w http.ResponseWriter, r *http.Request, user database.User) {
	location, err := cfg.pokeapiClient.GetLocationArea(strconv.Itoa(int(user.LocationID)))
	if err != nil {
		respondWithError(w, 500, "Failed to get user location: "+err.Error())
		return
	}

	type encounterOption struct {
		Name          string `json:"name"`
		EncounterRate int    `json:"encounter_rate"`
		Level         int    `json:"level"`
	}

	encounterOptions := make([]encounterOption, 0, len(location.PokemonEncounters))
	totalEncounterRate := 0

	for _, encounter := range location.PokemonEncounters {
		details := encounter.VersionDetails[0]

		minLevel := details.EncounterDetails[0].MinLevel
		maxLevel := details.EncounterDetails[0].MaxLevel

		level := rand.Intn(maxLevel-minLevel) + minLevel

		encounterOptions = append(encounterOptions, encounterOption{
			Name:          encounter.Pokemon.Name,
			EncounterRate: details.MaxChance,
			Level:         level,
		})
		totalEncounterRate += details.MaxChance
	}

	res := encounterOption{}

	rate := rand.Intn(totalEncounterRate)
	for _, option := range encounterOptions {
		rate -= option.EncounterRate
		if rate <= 0 {
			res.Name = option.Name
			res.Level = option.Level
			res.EncounterRate = option.EncounterRate
			break
		}
	}

	respondWithJSON(w, 200, res)
}
