package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/muratsat/chat/backend/db"
)

func register(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)

	// if request method is OPTIONS, return ok
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// else if request method is not POST, method not allowed
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// marshal request body to struct
	var request authRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Error decoding request body: ", err)
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// validate username and password
	err = request.validate()
	if err != nil {
		log.Println("Error validating request body: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// generate password hash
	hash, err := generatePasswordHash(request.Password)
	if err != nil {
		log.Println("Error generating password hash: ", err)
		http.Error(w, "Error generating password hash", http.StatusInternalServerError)
		return
	}

	// add user to database
	err = db.AddUser(request.Username, hash)
	if err != nil {
		if err.Error() == "User already exists" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Println("Error adding user to database: ", err)
			http.Error(w, "Error adding user to database", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(&w, "User added")
}
