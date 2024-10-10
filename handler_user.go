package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error parsing JSON: "+err.Error())
		return
	}
	if params.Username == "" || params.Password == "" {
		respondWithError(w, 400, "Username and password are required")
		return
	}
	_, err = apiCfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err == nil {
		respondWithError(w, 400, "Username is already taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, 400, "Could not create user: "+err.Error())
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Username:  params.Username,
		Password:  string(hash),
	})
	if err != nil {
		respondWithError(w, 400, "Could not create user: "+err.Error())
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerSelectStarterPokemon(w http.ResponseWriter, r *http.Request, user database.User) {
	log.Println("select-starter websocket connection with user: " + user.Username + " established")
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
		log.Println("select-starter websocket connection with user: " + user.Username + " closed")
	}()

	userPokemon, err := apiCfg.DB.GetUserPokemon(r.Context(), user.ID)
	if err != nil {
		conn.WriteJSON(wsMessage{Error: "Failed to get user pokemon: " + err.Error()})
		return
	}
	if len(userPokemon) > 0 {
		conn.WriteJSON(wsMessage{Error: "User already has pokemon"})
		return
	}

	conn.WriteJSON(wsMessage{Message: "Select your starter pokemon", Options: pokeutils.Starters})
	mt, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read message: " + err.Error())
		return
	}
	if mt != websocket.TextMessage {
		log.Println("Invalid message type")
		return
	}

	for _, starter := range pokeutils.Starters {
		if string(msg) == starter {
			type parameters struct {
				Name  string `json:"name"`
				Level int32  `json:"level"`
				Shiny bool   `json:"shiny"`
			}

			pokemon, err := apiCfg.handlerCreatePokemon(parameters{
				Name:  starter,
				Level: 5,
				Shiny: pokeutils.IsShiny(),
			}, r, user)
			if err != nil {
				conn.WriteJSON(wsMessage{Error: "Failed to create pokemon: " + err.Error()})
				return
			}

			conn.WriteJSON(pokemon)
			return
		}
	}

	conn.WriteJSON(wsMessage{Error: "Invalid starter pokemon"})
}
