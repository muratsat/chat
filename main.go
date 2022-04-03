package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func serveHttp(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	switch r.URL.Path {
	case "/":
		http.ServeFile(w, r, "home.html")
		return

	case "/register":
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "register.html")
		}
		if r.Method == http.MethodPost {
			register(w, r)
		}

	case "/login":
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "login.html")
		}
		if r.Method == http.MethodPost {
			login(w, r)
		}

	case "/friends":
		if r.Method == http.MethodGet {
			friends(w, r)
		}

	case "/friends/add":
		if r.Method == http.MethodPost {
			addFriend(w, r)
		}

	default:
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
}

func main() {
	http.HandleFunc("/", serveHttp)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newAuth auth
	err := decoder.Decode(&newAuth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("%s %s", newAuth.Username, newAuth.Password)

	if !dbAddUser(newAuth.Username, newAuth.Password) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("such user already exists"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newAuth auth
	err := decoder.Decode(&newAuth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !CheckCredentials(newAuth.Username, newAuth.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid username or password"))
		return
	}

	token, err := UpdateToken(newAuth.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(token)
}

func friends(w http.ResponseWriter, r *http.Request) {
	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	user_id, err := ValidateToken(authToken)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	friends_list := dbFriendsList(user_id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends_list)
}

func addFriend(w http.ResponseWriter, r *http.Request) {
	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	user_id, err := ValidateToken(authToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var f friend
	err = json.NewDecoder(r.Body).Decode(&f)

	found := dbAddFriend(user_id, f.Username)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
