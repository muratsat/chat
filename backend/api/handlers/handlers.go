package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SetupHandlers() {
	http.HandleFunc("/register", register)
}

func setupCORS(w *http.ResponseWriter, r *http.Request) {

	// get client origin from environment variable
	origin := os.Getenv("CLIENT_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	header := (*w).Header()
	header.Set("Access-Control-Allow-Origin", origin)
	if r.Method == "OPTIONS" {
		header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func writeJSON(w *http.ResponseWriter, data interface{}) {
	(*w).Header().Set("Content-Type", "application/json")
	json.NewEncoder(*w).Encode(data)
}
