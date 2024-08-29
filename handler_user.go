package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
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
	user, err := apiCfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err == nil {
		respondWithError(w, 400, "Username is already taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, 400, "Could not create user: "+err.Error())
	}

	user, err = apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
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
