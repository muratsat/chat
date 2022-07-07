package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/muratsat/chat/backend/db"
)

func register(w http.ResponseWriter, r *http.Request) {
	// marshal request body to struct
	var request authRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate request body
	err = request.validate()
	if err != nil {
		// return error message
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// add user to database
	err = db.AddUser(request.Username, request.Password)
	if err != nil {
		// return error message if user already exists
		if err.Error() == "User already exists" {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// return success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User added"))
}
