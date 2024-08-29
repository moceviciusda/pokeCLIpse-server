package main

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

type locationInfo struct {
	Name     string `json:"name"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func (cfg *apiConfig) hadlerGetUserLocation(w http.ResponseWriter, r *http.Request, user database.User) {
	url := "https://pokeapi.co/api/v2/location-area?offset=" + strconv.Itoa(int(user.LocationOffset)) + "&limit=1"

	areas, err := cfg.pokeapiClient.GetLocationAreas(url)
	if err != nil {
		respondWithError(w, 500, "Failed to get location options: "+err.Error())
		return
	}

	var next, previous string

	if areas.Next != "" {
		nl, err := cfg.pokeapiClient.GetLocationAreas(areas.Next)
		if err != nil {
			respondWithError(w, 500, "Failed to get next location: "+err.Error())
			return
		}
		next = nl.Results[0].Name
	}

	if areas.Previous != "" {
		pl, err := cfg.pokeapiClient.GetLocationAreas(areas.Previous)
		if err != nil {
			respondWithError(w, 500, "Failed to get previous location: "+err.Error())
			return
		}
		previous = pl.Results[0].Name
	}

	location, err := cfg.pokeapiClient.GetLocationArea(areas.Results[0].Name)
	if err != nil {
		respondWithError(w, 500, "Failed to get user location: "+err.Error())
		return
	}

	respondWithJSON(w, 200, locationInfo{Name: location.Name, Next: next, Previous: previous})
}

func (cfg *apiConfig) handlerNextLocation(w http.ResponseWriter, r *http.Request, user database.User) {
	url := "https://pokeapi.co/api/v2/location-area?offset=" + strconv.Itoa(int(user.LocationOffset)) + "&limit=1"

	areas, err := cfg.pokeapiClient.GetLocationAreas(url)
	if err != nil {
		respondWithError(w, 500, "Failed to get location options: "+err.Error())
		return
	}

	if areas.Next == "" {
		respondWithError(w, 400, "No next location")
		return
	}

	location, err := cfg.pokeapiClient.GetLocationAreas(areas.Next)
	if err != nil {
		respondWithError(w, 500, "Failed to get next location: "+err.Error())
		return
	}

	user.LocationOffset++
	_, err = cfg.DB.UpdateUserLocation(r.Context(), database.UpdateUserLocationParams{
		ID:             user.ID,
		LocationOffset: user.LocationOffset,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to update user location: "+err.Error())
		return
	}

	respondWithJSON(w, 200, struct {
		Name string `json:"name"`
	}{Name: location.Results[0].Name})
}

func (cfg *apiConfig) handlerPreviousLocation(w http.ResponseWriter, r *http.Request, user database.User) {
	url := "https://pokeapi.co/api/v2/location-area?offset=" + strconv.Itoa(int(user.LocationOffset)) + "&limit=1"

	areas, err := cfg.pokeapiClient.GetLocationAreas(url)
	if err != nil {
		respondWithError(w, 500, "Failed to get location options: "+err.Error())
		return
	}

	if areas.Previous == "" {
		respondWithError(w, 400, "No previous location")
		return
	}

	location, err := cfg.pokeapiClient.GetLocationAreas(areas.Previous)
	if err != nil {
		respondWithError(w, 500, "Failed to get previous location: "+err.Error())
		return
	}

	user.LocationOffset--
	_, err = cfg.DB.UpdateUserLocation(r.Context(), database.UpdateUserLocationParams{
		ID:             user.ID,
		LocationOffset: user.LocationOffset,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to update user location: "+err.Error())
		return
	}

	respondWithJSON(w, 200, struct {
		Name string `json:"name"`
	}{Name: location.Results[0].Name})
}

func (cfg *apiConfig) handlerSearchForPokemon(w http.ResponseWriter, r *http.Request, user database.User) {
	url := "https://pokeapi.co/api/v2/location-area?offset=" + strconv.Itoa(int(user.LocationOffset)) + "&limit=1"

	areas, err := cfg.pokeapiClient.GetLocationAreas(url)
	if err != nil {
		respondWithError(w, 500, "Failed to get location options: "+err.Error())
		return
	}

	location, err := cfg.pokeapiClient.GetLocationArea(areas.Results[0].Name)
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

		var level int
		if maxLevel > minLevel {
			level = rand.Intn(maxLevel-minLevel) + minLevel
		} else {
			level = minLevel
		}

		encounterOptions = append(encounterOptions, encounterOption{
			Name:          encounter.Pokemon.Name,
			EncounterRate: details.MaxChance,
			Level:         level,
		})
		totalEncounterRate += details.MaxChance
	}

	randomPokemon := encounterOption{}

	rate := rand.Intn(totalEncounterRate)
	for _, option := range encounterOptions {
		rate -= option.EncounterRate
		if rate <= 0 {
			randomPokemon = option
			break
		}
	}

	pokemon, err := cfg.pokeapiClient.GetPokemon(randomPokemon.Name)
	if err != nil {
		respondWithError(w, 500, "Failed to get pokemon: "+err.Error())
		return
	}

	res := struct {
		Name           string `json:"name"`
		Level          int    `json:"level"`
		HP             int    `json:"hp"`
		Attack         int    `json:"attack"`
		Defense        int    `json:"defense"`
		SpecialAttack  int    `json:"special_attack"`
		SpecialDefense int    `json:"special_defense"`
		Speed          int    `json:"speed"`
	}{}

	res.Name = pokemon.Name
	res.Level = randomPokemon.Level

	ivs := pokeutils.GenerateIVs()
	res.HP = pokeutils.CalculateStat(pokemon.Stats[0].BaseStat, ivs.HP, randomPokemon.Level)
	res.Attack = pokeutils.CalculateStat(pokemon.Stats[1].BaseStat, ivs.Attack, randomPokemon.Level)
	res.Defense = pokeutils.CalculateStat(pokemon.Stats[2].BaseStat, ivs.Defense, randomPokemon.Level)
	res.SpecialAttack = pokeutils.CalculateStat(pokemon.Stats[3].BaseStat, ivs.SpecialAttack, randomPokemon.Level)
	res.SpecialDefense = pokeutils.CalculateStat(pokemon.Stats[4].BaseStat, ivs.SpecialDefense, randomPokemon.Level)
	res.Speed = pokeutils.CalculateStat(pokemon.Stats[5].BaseStat, ivs.Speed, randomPokemon.Level)

	respondWithJSON(w, 200, res)
}
