package main

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
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

		encounterOptions = append(encounterOptions, encounterOption{
			Name:          encounter.Pokemon.Name,
			EncounterRate: details.MaxChance,
			Level:         rand.Intn(maxLevel-minLevel) + minLevel,
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
		HP             int    `json:"hp"`
		Attack         int    `json:"attack"`
		Defense        int    `json:"defense"`
		SpecialAttack  int    `json:"special_attack"`
		SpecialDefense int    `json:"special_defense"`
		Speed          int    `json:"speed"`
	}{}

	res.Name = pokemon.Name

	ivs := pokeutils.GenerateIVs()
	res.HP = pokeutils.CalculateStat(pokemon.Stats[0].BaseStat, ivs.HP, randomPokemon.Level)
	res.Attack = pokeutils.CalculateStat(pokemon.Stats[1].BaseStat, ivs.Attack, randomPokemon.Level)
	res.Defense = pokeutils.CalculateStat(pokemon.Stats[2].BaseStat, ivs.Defense, randomPokemon.Level)
	res.SpecialAttack = pokeutils.CalculateStat(pokemon.Stats[3].BaseStat, ivs.SpecialAttack, randomPokemon.Level)
	res.SpecialDefense = pokeutils.CalculateStat(pokemon.Stats[4].BaseStat, ivs.SpecialDefense, randomPokemon.Level)
	res.Speed = pokeutils.CalculateStat(pokemon.Stats[5].BaseStat, ivs.Speed, randomPokemon.Level)

	respondWithJSON(w, 200, res)
}
