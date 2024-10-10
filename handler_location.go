package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/internal/pokeapi"
	"github.com/moceviciusda/pokeCLIpse-server/internal/pokebattle"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/ansiiutils"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

type locationInfo struct {
	Name     string `json:"name"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type message struct {
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Options []string `json:"options"`
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

	defer func() {
		conn.Close()
		log.Println("Websocket connection with user: " + user.Username + " closed")
	}()

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

	rMoves, err := cfg.pokeapiClient.SelectRandomMoves(p.Name, randomPokemon.Level)
	if err != nil {
		log.Println("Failed to get user: " + user.Username + " pokemon moves: " + err.Error())
		conn.WriteJSON(errResponse{Error: "Failed to get pokemon moves: " + err.Error()})
		return
	}

	moves := make([]pokeutils.Move, 0, 4)
	for _, move := range rMoves {
		dbMove, err := cfg.getMoveFromDbOrApi(r, move.Name)
		if err != nil {
			log.Println("Failed to get move: " + err.Error())
			conn.WriteJSON(errResponse{Error: "Failed to get move: " + err.Error()})
			return
		}

		moves = append(moves, dbMoveToMove(dbMove))
	}

	types := make([]string, 0, len(p.Types))
	for _, t := range p.Types {
		types = append(types, t.Type.Name)
	}
	ivs := pokeutils.GenerateIVs()
	pokemon := pokeutils.Pokemon{
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
	case "run":
		conn.WriteJSON(errResponse{Error: "You ran away!"})
		return
	case "battle":
		dbParty, err := cfg.DB.GetPokemonParty(r.Context(), user.ID)
		if err != nil {
			log.Println("Failed to get user: " + user.Username + " party: " + err.Error())
			conn.WriteJSON(errResponse{Error: "Failed to get pokemon party: " + err.Error()})
			return
		}

		pokemonParty := make([]pokebattle.Pokemon, 0, len(dbParty))
		for _, pokemon := range dbParty {
			p, err := cfg.pokeapiClient.GetPokemon(pokemon.Name)
			if err != nil {
				log.Println("Failed to get user: " + user.Username + " party pokemon: " + err.Error())
				conn.WriteJSON(errResponse{Error: "Failed to get party pokemon: " + err.Error()})
				return
			}

			dbMoves, err := cfg.DB.GetMovesByPokemonID(r.Context(), pokemon.ID)
			if err != nil {
				log.Println("Failed to get user: " + user.Username + " party pokemon moves: " + err.Error())
				conn.WriteJSON(errResponse{Error: "Failed to get party pokemon moves: " + err.Error()})
				return
			}

			dbIvs, err := cfg.DB.GetIVs(r.Context(), pokemon.IvsID)
			if err != nil {
				log.Println("Failed to get user: " + user.Username + " party pokemon IVs: " + err.Error())
				conn.WriteJSON(errResponse{Error: "Failed to get party pokemon IVs: " + err.Error()})
				return
			}

			pokemonParty = append(pokemonParty, pokebattle.Pokemon{
				Pokemon: makePokemon(p, pokemon, dbMoves, dbIvs),
				ExpGain: 0,
				BaseExp: p.BaseExperience,
			})
		}

		battle := pokebattle.NewBattle(
			pokebattle.Trainer{
				Name:    user.Username,
				Pokemon: pokemonParty,
			},
			pokebattle.Trainer{
				Name:    "Wild",
				Pokemon: []pokebattle.Pokemon{{Pokemon: pokemon, ExpGain: 0, BaseExp: p.BaseExperience}},
			},
			make(chan pokebattle.BattleMessage),
		)

		go battle.Run()

		for battleMsg := range battle.MsgChan {
			switch battleMsg.Type {
			case pokebattle.BattleMsgInfo:
				var color string
				if battleMsg.Subject == user.Username {
					color = ansiiutils.ColorGreen
				} else if battleMsg.Subject == "Wild" {
					color = ansiiutils.ColorRed
				} else {
					color = ansiiutils.ColorYellow
				}
				conn.WriteJSON(message{Message: color + battleMsg.Message + ansiiutils.Reset})

			case pokebattle.BattleMsgSelect:
				conn.WriteJSON(message{Message: "Select a pokemon", Options: battleMsg.Options})
				mt, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("Failed to read message: " + err.Error())
					return
				}
				if mt != websocket.TextMessage {
					log.Println("Invalid message type")
					return
				}

				battle.MsgChan <- pokebattle.BattleMessage{Type: pokebattle.BattleMsgAction, Message: string(msg), Subject: user.Username}
			}
		}

		if battle.Winner.Name != user.Username {
			conn.WriteJSON(message{Message: "You lost!"})
			return
		}

		for i, pokemon := range battle.Trainers[0].Pokemon {
			if pokemon.ExpGain <= 0 {
				continue
			}
			conn.WriteJSON(message{Message: ansiiutils.StyleItalic + pokemon.Pokemon.Name + " gained " + strconv.Itoa(pokemon.ExpGain) + " exp" + ansiiutils.Reset})

			dbPokemon := dbParty[i]
			dbPokemon, movesChanged := cfg.resolveExpGains(dbPokemon, &pokemon, conn)
			dbPokemon, err = cfg.DB.UpdatePokemonLvlAndExp(r.Context(), database.UpdatePokemonLvlAndExpParams{
				Level:      dbPokemon.Level,
				Experience: dbPokemon.Experience,
				ID:         dbPokemon.ID,
			})
			if err != nil {
				log.Println("Failed to LVL up pokemon: " + err.Error())
				break
			}

			if !movesChanged {
				continue
			}

			success := true
			for _, move := range pokemon.Moves {
				_, err := cfg.getMoveFromDbOrApi(r, move.Name)
				if err != nil {
					log.Println("Failed to get move: " + err.Error())
					success = false
				}
			}
			if !success {
				conn.WriteJSON(message{Message: "Failed to learn new moves"})
				continue
			}

			err = cfg.DB.RemoveAllMovesFromPokemon(r.Context(), dbPokemon.ID)
			if err != nil {
				log.Println("Failed to delete pokemon moves: " + err.Error())
				conn.WriteJSON(message{Message: "Failed to learn new moves"})
				continue
			}
			for _, move := range pokemon.Moves {
				_, err = cfg.DB.AddMoveToPokemon(r.Context(), database.AddMoveToPokemonParams{
					PokemonID: dbPokemon.ID,
					MoveName:  move.Name,
				})
				if err != nil {
					log.Println("Failed to add move to pokemon: " + err.Error())
					conn.WriteJSON(message{Message: "Failed to learn new moves"})
					continue
				}
			}
		}

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
		switch string(msg) {
		case "no":
			conn.WriteJSON(message{Message: "You left " + pokemon.Name + " in the wild.."})
			return
		case "yes":
			if ownedP.ID != uuid.Nil {
				_, err := cfg.DB.DeletePokemon(r.Context(), ownedP.ID)
				if err != nil {
					log.Println("Failed to delete pokemon: " + err.Error())
					conn.WriteJSON(message{Message: "Failed to release pokemon"})
					return
				}
			}

			dbIvs, err := cfg.DB.CreateIVs(r.Context(), database.CreateIVsParams{
				ID:             uuid.New(),
				CreatedAt:      time.Now().UTC(),
				UpdatedAt:      time.Now().UTC(),
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
				ID:         uuid.New(),
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
				OwnerID:    user.ID,
				Name:       pokemon.Name,
				Experience: int32(pokeutils.ExpAtLevel(pokemon.Level)),
				Level:      int32(pokemon.Level),
				Shiny:      pokemon.Shiny,
				IvsID:      dbIvs.ID,
			})
			if err != nil {
				log.Println("Failed to create pokemon: " + err.Error())
				conn.WriteJSON(message{Message: "Failed to create pokemon: " + err.Error()})
				return
			}
			for _, move := range pokemon.Moves {
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

			if len(dbParty) < 6 {
				_, err := cfg.DB.AddPokemonToParty(r.Context(), database.AddPokemonToPartyParams{
					UserID:    user.ID,
					PokemonID: dbPokemon.ID,
					Position:  int32(len(dbParty) + 1),
				})
				if err != nil {
					log.Println("Failed to add pokemon to party: " + err.Error())
					respondWithError(w, 500, "Failed to add pokemon to party: "+err.Error())
					return
				}
			}
		}

	default:
		log.Println("Invalid message: " + string(msg))
	}
}

func (cfg *apiConfig) resolveExpGains(dbPokemon database.Pokemon, pokemon *pokebattle.Pokemon, conn *websocket.Conn) (_ database.Pokemon, movesChanged bool) {
	dbPokemon.Experience += int32(pokemon.ExpGain)

	expectedLevel := pokeutils.LevelAtExp(int(dbPokemon.Experience))
	if dbPokemon.Level >= 100 || dbPokemon.Level >= int32(expectedLevel) {
		return dbPokemon, movesChanged
	}

	p, err := cfg.pokeapiClient.GetPokemon(dbPokemon.Name)
	if err != nil {
		log.Println("Failed to get pokemon: " + err.Error())
		conn.WriteJSON(errResponse{Error: "Failed to level up pokemon: " + err.Error()})
		return dbPokemon, movesChanged
	}
	pSpecies, err := cfg.pokeapiClient.GetPokemonSpecies(p.Species.Name, p.Species.URL)
	if err != nil {
		log.Println("Failed to get pokemon species: " + err.Error())
		conn.WriteJSON(errResponse{Error: "Failed to level up pokemon: " + err.Error()})
		return dbPokemon, movesChanged
	}
	evolutionChain, err := cfg.pokeapiClient.GetEvolutionChain(pSpecies.EvolutionChain.URL)
	if err != nil {
		log.Println("Failed to get evolution chain: " + err.Error())
		conn.WriteJSON(errResponse{Error: "Failed to level up pokemon: " + err.Error()})
		return dbPokemon, movesChanged
	}

	var evolvesAt int
	var evolvesTo string
	for _, e := range findEvolutionOptions(evolutionChain.Chain, dbPokemon.Name) {
		for _, opt := range e.EvolutionDetails {
			if opt.Trigger.Name != "level-up" {
				continue
			}
			evolvesAt = opt.MinLevel
			evolvesTo = e.Species.Name
		}
	}

	for dbPokemon.Level < int32(expectedLevel) {
		dbPokemon.Level++
		conn.WriteJSON(message{Message: ansiiutils.StyleItalic + dbPokemon.Name + ansiiutils.Reset + " leveled up and is now lvl " + strconv.Itoa(int(dbPokemon.Level)) + "!"})

		movesToLearn, err := cfg.pokeapiClient.GetMovesLearnedAtLvl(p.Name, int(dbPokemon.Level))
		if err != nil {
			log.Println("Failed to get moves learned at lvl: " + err.Error())
			conn.WriteJSON(errResponse{Error: "Failed to level up pokemon: " + err.Error()})
			return dbPokemon, movesChanged
		}
		movesChanged = moveLearnLoop(conn, movesToLearn, pokemon)

		if evolvesTo != "" && int(dbPokemon.Level) < evolvesAt {
			continue
		}

		conn.WriteJSON(message{Message: ansiiutils.StyleBlink + dbPokemon.Name + " is evolving!" + ansiiutils.Reset})
		// TODO: Implement evolution rejection
		dbPokemon.Name = evolvesTo
		p, err := cfg.pokeapiClient.GetPokemon(evolvesTo)
		if err != nil {
			log.Println("Failed to get pokemon: " + err.Error())
			return dbPokemon, movesChanged
		}
		conn.WriteJSON(message{Message: dbPokemon.Name + " evolved into " + evolvesTo + "!"})

		movesToLearn, err = cfg.pokeapiClient.GetMovesLearnedAtLvl(p.Name, int(dbPokemon.Level))
		if err != nil {
			log.Println("Failed to get moves learned at lvl: " + err.Error())
			return dbPokemon, movesChanged
		}
		movesChanged = moveLearnLoop(conn, movesToLearn, pokemon)
	}

	return dbPokemon, movesChanged
}

func findEvolutionOptions(evolutionChain pokeapi.EvolutionChainLink, pokemonName string) []pokeapi.EvolutionChainLink {
	if evolutionChain.Species.Name == pokemonName {
		return evolutionChain.EvolvesTo
	}

	for _, chain := range evolutionChain.EvolvesTo {
		if opts := findEvolutionOptions(chain, pokemonName); opts != nil {
			return opts
		}
	}

	return nil
}

func moveLearnLoop(conn *websocket.Conn, movesToLearn map[string]pokeapi.MoveResponse, pokemon *pokebattle.Pokemon) (movesChanged bool) {
	for moveName, m := range movesToLearn {
		typeIcon := pokeutils.TypeIcons[m.Type.Name]

		if len(pokemon.Moves) < 4 {
			conn.WriteJSON(message{Message: ansiiutils.StyleItalic + pokemon.Name + ansiiutils.Reset + " learned " + ansiiutils.StyleBold + moveName + typeIcon + ansiiutils.Reset + "\n"})
			pokemon.Moves = append(pokemon.Moves, pokeutils.Move{
				Name:         m.Name,
				Accuracy:     m.Accuracy,
				Power:        m.Power,
				PP:           m.Pp,
				Type:         m.Type.Name,
				DamageClass:  m.DamageClass.Name,
				EffectChance: m.EffectChance,
				Effect:       "",
			})
			movesChanged = true
			continue
		}

		forgetMoveOpts := make([]string, 0, len(pokemon.Moves))
		for _, pM := range pokemon.Moves {
			forgetMoveOpts = append(forgetMoveOpts, pM.Name+" "+pokeutils.TypeIcons[pM.Type])
		}
		forgetMoveOpts = append(forgetMoveOpts, "cancel")

		conn.WriteJSON(message{Message: ansiiutils.StyleItalic + pokemon.Name + ansiiutils.Reset + " is trying to learn " + ansiiutils.StyleBold + moveName + typeIcon + ansiiutils.Reset})
		conn.WriteJSON(message{Message: "Select a move to" + ansiiutils.ColorRed + " forget" + ansiiutils.Reset + ":", Options: forgetMoveOpts})
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message: " + err.Error())
			continue
		}
		if mt != websocket.TextMessage {
			log.Println("Invalid message type")
			continue
		}

		if string(msg) == "cancel" {
			conn.WriteJSON(message{Message: pokemon.Name + " did not learn " + ansiiutils.StyleBold + moveName + ansiiutils.Reset + "\n"})
			continue
		}

		for i, mName := range forgetMoveOpts {
			if mName == string(msg) {
				pokemon.Moves[i] = pokeutils.Move{
					Name:         m.Name,
					Accuracy:     m.Accuracy,
					Power:        m.Power,
					PP:           m.Pp,
					Type:         m.Type.Name,
					DamageClass:  m.DamageClass.Name,
					EffectChance: m.EffectChance,
					Effect:       "",
				}
				conn.WriteJSON(message{Message: ansiiutils.StyleItalic + pokemon.Name + ansiiutils.Reset + " forgot " + ansiiutils.StyleBold + pokemon.Moves[i].Name + pokeutils.TypeIcons[pokemon.Moves[i].Type] + ansiiutils.Reset + " and learned " + ansiiutils.StyleBold + moveName + typeIcon + ansiiutils.Reset + "\n"})
				movesChanged = true
				break
			}
		}
	}

	return movesChanged
}
