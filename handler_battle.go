package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/moceviciusda/pokeCLIpse-server/internal/pokebattle"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/ansiiutils"
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
		Guest struct {
			Name         string               `json:"name"`
			PokemonParty []pokebattle.Pokemon `json:"pokemonParty"`
		} `json:"guest"`
		Opponent struct {
			Name         string               `json:"name"`
			PokemonParty []pokebattle.Pokemon `json:"pokemonParty"`
		} `json:"opponent"`
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

	battle := pokebattle.NewBattle(
		pokebattle.Trainer{
			Name:    initMsg.Guest.Name,
			Pokemon: initMsg.Guest.PokemonParty,
		},
		pokebattle.Trainer{
			Name:    initMsg.Opponent.Name,
			Pokemon: initMsg.Opponent.PokemonParty,
		},
		make(chan pokebattle.BattleMessage),
	)

	go battle.Run()

	for battleMsg := range battle.MsgChan {
		switch battleMsg.Type {
		case pokebattle.BattleMsgInfo:
			var color string
			if battleMsg.Subject == initMsg.Guest.Name {
				color = ansiiutils.ColorGreen
			} else if battleMsg.Subject == initMsg.Opponent.Name {
				color = ansiiutils.ColorRed
			} else {
				color = ansiiutils.ColorYellow
			}
			conn.WriteJSON(wsMessage{Message: color + battleMsg.Message + ansiiutils.Reset})

		case pokebattle.BattleMsgSelect:
			if battleMsg.Subject != initMsg.Guest.Name {
				// randomly select a pokemon for the opponent
				randomPokemon := rand.Intn(len(battleMsg.Options))

				battle.MsgChan <- pokebattle.BattleMessage{
					Type:    pokebattle.BattleMsgAction,
					Message: battleMsg.Options[randomPokemon],
					Subject: initMsg.Opponent.Name,
				}
				continue
			}

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

			battle.MsgChan <- pokebattle.BattleMessage{Type: pokebattle.BattleMsgAction, Message: string(msg), Subject: initMsg.Guest.Name}
		}
	}

	if battle.Winner.Name != initMsg.Guest.Name {
		conn.WriteJSON(wsMessage{Message: "You lost!"})
		return
	}
}
