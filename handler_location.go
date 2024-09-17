package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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

	moveOptions := make(map[string]database.Move)
	for _, move := range p.Moves {
		if _, ok := moveOptions[move.Move.Name]; ok {
			continue
		}

		for _, details := range move.VersionGroupDetails {
			if details.LevelLearnedAt > randomPokemon.Level {
				continue
			}
			if !(details.MoveLearnMethod.Name == "level-up" || details.MoveLearnMethod.Name == "egg") {
				continue
			}

			dbm, err := cfg.DB.GetMoveByName(r.Context(), move.Move.Name)
			if err != nil {
				m, err := cfg.pokeapiClient.GetMove(move.Move.Name)
				if err != nil {
					conn.WriteJSON(errResponse{Error: "Failed to get move: " + err.Error()})
					return
				}

				dbm, err = cfg.DB.CreateMove(r.Context(), database.CreateMoveParams{
					ID:           uuid.New(),
					CreatedAt:    user.CreatedAt,
					UpdatedAt:    user.UpdatedAt,
					Name:         m.Name,
					Accuracy:     int32(m.Accuracy),
					Power:        int32(m.Power),
					Pp:           int32(m.Pp),
					Type:         m.Type.Name,
					DamageClass:  m.DamageClass.Name,
					EffectChance: int32(m.EffectChance),
					Effect:       m.EffectEntries[0].ShortEffect,
				})
				if err != nil {
					log.Println("Failed to create move: " + err.Error())
					conn.WriteJSON(errResponse{Error: "Failed to create move: " + err.Error()})
					return
				}
			}

			moveOptions[move.Move.Name] = dbm
			break
		}
	}

	moves := make([]pokeutils.Move, 0, 4)
	for i := 0; i < 4; i++ {
		if len(moveOptions) == 0 || len(moves) == 4 {
			break
		}

		moveOptKeys := make([]string, 0, len(moveOptions))
		for k := range moveOptions {
			moveOptKeys = append(moveOptKeys, k)
		}
		moveName := moveOptKeys[rand.Intn(len(moveOptKeys))]

		moves = append(moves, dbMoveToMove(moveOptions[moveName]))
		delete(moveOptions, moveName)
	}

	types := make([]string, 0, len(p.Types))
	for _, t := range p.Types {
		types = append(types, t.Type.Name)
	}

	ivs := pokeutils.GenerateIVs()

	pokemon := pokeutils.Pokemon{
		ID:    "",
		Name:  p.Name,
		Types: types,
		Level: randomPokemon.Level,
		Shiny: pokeutils.IsShiny(),
		Stats: pokeutils.CalculateStats(pBaseStats, ivs, randomPokemon.Level),
		Moves: moves,
	}

	conn.WriteJSON(pokemon)

	mt, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read message: " + err.Error())
		return
	}

	if mt != websocket.TextMessage {
		log.Println("Invalid message type")
		return
	}

	switch string(msg) {
	case "battle":
		battleMsgChan := make(chan string)

		battle := pokebattle.NewBattle(pokebattle.Trainer{
			Name:    user.Username,
			Pokemon: []pokeutils.Pokemon{pokeutils.Pikachu},
		}, pokebattle.Trainer{
			Name:    "Wild",
			Pokemon: []pokeutils.Pokemon{pokemon},
		},
			battleMsgChan,
		)

		type message struct {
			Error   string   `json:"error"`
			Message string   `json:"message"`
			Options []string `json:"options"`
		}

		go battle.Run()

		for battleMsg := range battleMsgChan {
			conn.WriteJSON(message{Message: battleMsg})
		}

		if battle.Winner.Name == user.Username {
			prompt := fmt.Sprintf(pokemon.Name + " is weakened, attempt to catch it?")
			ownedP, err := cfg.DB.CheckHasPokemon(r.Context(), database.CheckHasPokemonParams{
				OwnerID: user.ID,
				Name:    pokemon.Name,
				Shiny:   pokemon.Shiny,
			})
			if err == nil {
				prompt += fmt.Sprintf("\nYou already have a lvl %d %s. If you catch this one, the old one will be released", ownedP.Level, ownedP.Name)
			}

			conn.WriteJSON(message{Message: prompt, Options: []string{"yes", "no"}})

			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Failed to read message: " + err.Error())
				return
			}
			if string(msg) == "yes" {
				cfg.DB.DeletePokemon(r.Context(), ownedP.ID)

				dbIvs, err := cfg.DB.CreateIVs(r.Context(), database.CreateIVsParams{
					ID:             uuid.New(),
					CreatedAt:      user.CreatedAt,
					UpdatedAt:      user.UpdatedAt,
					Hp:             int32(ivs.Hp),
					Attack:         int32(ivs.Attack),
					Defense:        int32(ivs.Defense),
					SpecialAttack:  int32(ivs.SpecialAttack),
					SpecialDefense: int32(ivs.SpecialDefense),
					Speed:          int32(ivs.Speed),
				})
				if err != nil {
					log.Println("Failed to create IVs: " + err.Error())
					conn.WriteJSON(message{Message: "Failed to create IVs: " + err.Error()})
					return
				}

				dbPokemon, err := cfg.DB.CreatePokemon(r.Context(), database.CreatePokemonParams{
					ID:        uuid.New(),
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
					OwnerID:   user.ID,
					Name:      pokemon.Name,
					Level:     int32(pokemon.Level),
					Shiny:     pokemon.Shiny,
					IvsID:     dbIvs.ID,
				})
				if err != nil {
					log.Println("Failed to create pokemon: " + err.Error())
					conn.WriteJSON(message{Message: "Failed to create pokemon: " + err.Error()})
					return
				}

				for _, move := range pokemon.Moves {
					log.Println("Adding move to pokemon: " + move.Name)
					_, err = cfg.DB.AddMoveToPokemon(r.Context(), database.AddMoveToPokemonParams{
						PokemonID: dbPokemon.ID,
						MoveName:  move.Name,
					})
					if err != nil {
						log.Println("Failed to add move to pokemon: " + err.Error())
						conn.WriteJSON(message{Message: "Failed to add move to pokemon: " + err.Error()})
						return
					}
				}

				conn.WriteJSON(message{Message: "You caught " + pokemon.Name + "!"})
			}

		} else {
			conn.WriteJSON(message{Message: "You lost!"})
		}

	default:
		log.Println("Invalid message: " + string(msg))
	}

}
