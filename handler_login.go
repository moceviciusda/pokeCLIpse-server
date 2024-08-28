package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		ID       uuid.UUID `json:"id"`
		Token    string    `json:"token"`
		Username string    `json:"username"`
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
		respondWithError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "pokeclipse",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24).UTC()),
		Subject:   user.ID.String(),
	})

	signedT, err := token.SignedString([]byte(apiCfg.jwtSecret))
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	respondWithJSON(w, 200, response{
		ID:       user.ID,
		Token:    signedT,
		Username: user.Username,
	})
}
