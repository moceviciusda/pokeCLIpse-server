package main

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerparseJWT(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString := authorization[len("Bearer "):]

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiCfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Printf("Error parsing UUID: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	user, err := apiCfg.DB.GetUserById(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))

}
