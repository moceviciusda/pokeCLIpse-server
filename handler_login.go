package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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
	user, err := apiCfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(database.User{
		uuid.Max,
		time.Now().UTC(),
		time.Now(),
		params.Username,
		params.Password,
	}))
}
