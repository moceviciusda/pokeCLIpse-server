package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal JSON response:", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
