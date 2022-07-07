package handlers

import "net/http"

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SetupHandlers() {
	http.HandleFunc("/register", register)
}
