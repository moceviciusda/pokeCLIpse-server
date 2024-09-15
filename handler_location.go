package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/internal/pokebattle"
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
		log.Println("Failed to get user: " + user.Username + " location options: " + err.Error())
		respondWithError(w, 500, "Failed to get location options: "+err.Error())
		return
	}

	var next, previous string

	if areas.Next != "" {
		nl, err := cfg.pokeapiClient.GetLocationAreas(areas.Next)
		if err != nil {
			log.Println("Failed to get user: " + user.Username + " next location: " + err.Error())
			respondWithError(w, 500, "Failed to get next location: "+err.Error())
			return
		}
		next = nl.Results[0].Name
	}

	if areas.Previous != "" {
		pl, err := cfg.pokeapiClient.GetLocationAreas(areas.Previous)
		if err != nil {
			log.Println("Failed to get user: " + user.Username + " previous location: " + err.Error())
			respondWithError(w, 500, "Failed to get previous location: "+err.Error())
			return
		}
		previous = pl.Results[0].Name
	}

	location, err := cfg.pokeapiClient.GetLocationArea(areas.Results[0].Name)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " location: " + err.Error())
		respondWithError(w, 500, "Failed to get user location: "+err.Error())
		return
	}

	respondWithJSON(w, 200, locationInfo{Name: location.Name, Next: next, Previous: previous})
}

func (cfg *apiConfig) handlerNextLocation(w http.ResponseWriter, r *http.Request, user database.User) {
	url := "https://pokeapi.co/api/v2/location-area?offset=" + strconv.Itoa(int(user.LocationOffset)) + "&limit=1"

	areas, err := cfg.pokeapiClient.GetLocationAreas(url)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " location options: " + err.Error())
		respondWithError(w, 500, "Failed to get location options: "+err.Error())
		return
	}

	if areas.Next == "" {
		respondWithError(w, 400, "No next location")
		return
	}

	location, err := cfg.pokeapiClient.GetLocationAreas(areas.Next)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " next location: " + err.Error())
		respondWithError(w, 500, "Failed to get next location: "+err.Error())
		return
	}

	user.LocationOffset++
	_, err = cfg.DB.UpdateUserLocation(r.Context(), database.UpdateUserLocationParams{
		ID:             user.ID,
		LocationOffset: user.LocationOffset,
	})
	if err != nil {
		log.Println("Failed to update user: " + user.Username + " location: " + err.Error())
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
		log.Println("Failed to get user: " + user.Username + " location options: " + err.Error())
		respondWithError(w, 500, "Failed to get location options: "+err.Error())
		return
	}

	if areas.Previous == "" {
		respondWithError(w, 400, "No previous location")
		return
	}

	location, err := cfg.pokeapiClient.GetLocationAreas(areas.Previous)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " previous location: " + err.Error())
		respondWithError(w, 500, "Failed to get previous location: "+err.Error())
		return
	}

	user.LocationOffset--
	_, err = cfg.DB.UpdateUserLocation(r.Context(), database.UpdateUserLocationParams{
		ID:             user.ID,
		LocationOffset: user.LocationOffset,
	})
	if err != nil {
		log.Println("Failed to update user: " + user.Username + " location: " + err.Error())
		respondWithError(w, 500, "Failed to update user location: "+err.Error())
		return
	}

	respondWithJSON(w, 200, struct {
		Name string `json:"name"`
	}{Name: location.Results[0].Name})
}

func (cfg *apiConfig) handlerSearchForPokemon(w http.ResponseWriter, r *http.Request, user database.User) {
	log.Println("Websocket connection with user: " + user.Username + " established")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection: " + err.Error())
		return
	}

	defer conn.Close()

	url := "https://pokeapi.co/api/v2/location-area?offset=" + strconv.Itoa(int(user.LocationOffset)) + "&limit=1"

	areas, err := cfg.pokeapiClient.GetLocationAreas(url)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " location options: " + err.Error())
		conn.WriteJSON(errResponse{Error: "Failed to get location options: " + err.Error()})
		return
	}

	location, err := cfg.pokeapiClient.GetLocationArea(areas.Results[0].Name)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " location: " + err.Error())
		conn.WriteJSON(errResponse{Error: "Failed to get user location: " + err.Error()})
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

	p, err := cfg.pokeapiClient.GetPokemon(randomPokemon.Name)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " pokemon: " + err.Error())
		conn.WriteJSON(errResponse{Error: "Failed to get pokemon: " + err.Error()})
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

	moves := make([]pokeutils.Move, 0, len(p.Moves))
	for _, move := range p.Moves {
		m, err := cfg.pokeapiClient.GetMove(move.Move.Name)
		if err != nil {
			log.Println("Failed to get user: " + user.Username + " move: " + err.Error())
			conn.WriteJSON(errResponse{Error: "Failed to get move: " + err.Error()})
			return
		}

		moves = append(moves, pokeutils.Move{
			Name:         m.Name,
			Accuracy:     m.Accuracy,
			Power:        m.Power,
			Type:         m.Type.Name,
			DamageClass:  m.DamageClass.Name,
			EffectChance: m.EffectChance,
			Effect:       "",
		})
	}

	types := make([]string, 0, len(p.Types))
	for _, t := range p.Types {
		types = append(types, t.Type.Name)
	}

	pokemon := pokeutils.Pokemon{
		ID:    "",
		Name:  p.Name,
		Types: types,
		Level: randomPokemon.Level,
		Shiny: pokeutils.IsShiny(),
		Stats: pokeutils.CalculateStats(pBaseStats, pokeutils.GenerateIVs(), randomPokemon.Level),
		Moves: moves,
	}

	battle := pokebattle.NewBattle(pokebattle.Trainer{
		Name:    user.Username,
		Pokemon: []pokeutils.Pokemon{pokemon},
	}, pokebattle.Trainer{
		Name:    "Wild " + p.Name,
		Pokemon: []pokeutils.Pokemon{pokemon},
	})

	conn.WriteJSON(pokemon)

	battle.Run()

	conn.Close()
}
