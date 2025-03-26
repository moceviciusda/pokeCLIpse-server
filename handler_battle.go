package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/moceviciusda/pokeCLIpse-server/internal/pokebattle"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/ansiiutils"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

func (cfg *apiConfig) handlerSimulateBattle(w http.ResponseWriter, r *http.Request) {
	log.Println("Established simulate-battle websocket connection with client: " + r.RemoteAddr)
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

	type initMessage struct {
		PokemonParty []pokebattle.Pokemon `json:"pokemonParty"`
		Opponent     string               `json:"opponent"`
	}

	var initMsg initMessage
	err = conn.ReadJSON(&initMsg)
	if err != nil {
		log.Println("Failed to read initial message: " + err.Error())
		return
	}

	defer func() {
		conn.Close()
		log.Println("Closed simulate-battle websocket connection with client: " + r.RemoteAddr)
	}()

	pokemonParty := make([]pokebattle.Pokemon, 0, 1)

	p := pokeutils.Pokemon{
		Name:  "Pikachu",
		Types: []string{"electric"},
		Level: 5,
		Stats: pokeutils.Stats{
			Hp:             35,
			Attack:         55,
			Defense:        40,
			SpecialAttack:  50,
			SpecialDefense: 50,
			Speed:          90,
		},
		Moves: []pokeutils.Move{
			{
				Name:         "Thunder Shock",
				Accuracy:     100,
				Power:        40,
				PP:           30,
				Type:         "electric",
				DamageClass:  "special",
				EffectChance: 10,
				Effect:       "paralyze",
			},
			{
				Name:         "Quick Attack",
				Accuracy:     100,
				Power:        40,
				PP:           30,
				Type:         "normal",
				DamageClass:  "physical",
				EffectChance: 0,
				Effect:       "",
			},
		},
	}

	pokemon := pokebattle.Pokemon{
		Pokemon: p,
		ExpGain: 0,
		BaseExp: 0,
	}

	pokemonParty = append(pokemonParty, pokemon)

	battle := pokebattle.NewBattle(
		pokebattle.Trainer{
			Name:    "Guest",
			Pokemon: initMsg.PokemonParty,
		},
		pokebattle.Trainer{
			Name:    "Wild",
			Pokemon: pokemonParty,
		},
		make(chan pokebattle.BattleMessage),
	)

	go battle.Run()

	for battleMsg := range battle.MsgChan {
		switch battleMsg.Type {
		case pokebattle.BattleMsgInfo:
			var color string
			if battleMsg.Subject == "Guest" {
				color = ansiiutils.ColorGreen
			} else if battleMsg.Subject == "Wild" {
				color = ansiiutils.ColorRed
			} else {
				color = ansiiutils.ColorYellow
			}
			conn.WriteJSON(wsMessage{Message: color + battleMsg.Message + ansiiutils.Reset})

		case pokebattle.BattleMsgSelect:
			conn.WriteJSON(wsMessage{Message: "Select a pokemon", Options: battleMsg.Options})
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Failed to read message: " + err.Error())
				return
			}
			if mt != websocket.TextMessage {
				log.Println("Invalid message type")
				return
			}

			battle.MsgChan <- pokebattle.BattleMessage{Type: pokebattle.BattleMsgAction, Message: string(msg), Subject: "Guest"}
		}
	}

	if battle.Winner.Name != "Guest" {
		conn.WriteJSON(wsMessage{Message: "You lost!"})
		return
	}
}
