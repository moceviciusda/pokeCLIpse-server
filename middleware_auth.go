package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		vals := strings.Split(authorization, " ")
		if len(vals) != 2 || strings.ToLower(vals[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Malformed Authorization header")
			return
		}

		token, err := jwt.ParseWithClaims(vals[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
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

		user, err := cfg.DB.GetUserById(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		next(w, r, user)
	}
}
