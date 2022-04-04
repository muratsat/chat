package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var hub = newHub()

func main() {
	go hub.run()

	http.HandleFunc("/", home)
	http.HandleFunc("/friends", friends)
	http.HandleFunc("/friends/add", addFriend)
	http.HandleFunc("/messages", messages)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "register.html")
		} else {
			register(w, r)
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "login.html")
		} else {
			login(w, r)
		}
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

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
	log.Print(newAuth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("%s %s", newAuth.Username, newAuth.Password)

	if !dbAddUser(newAuth.Username, newAuth.Password) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("such user already exists or invalid username")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("succesfully registered")
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

	expires, _ := time.Parse(timeLayout, fmt.Sprintf("%s", token["expires_at"]))

	cookie := http.Cookie{
		Name:    "auth_token",
		Value:   fmt.Sprintf("%s", token["auth_token"]),
		Expires: expires,
	}

	http.SetCookie(w, &cookie)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(token)
}

func friends(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authToken := c.Value
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
	c, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authToken := c.Value
	user_id, err := ValidateToken(authToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var f friend
	err = json.NewDecoder(r.Body).Decode(&f)
	log.Println(user_id, f.Username)

	found := dbAddFriend(user_id, f.Username)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("username not found")
		return
	}
}

func messages(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authToken := c.Value
	user_id, err := ValidateToken(authToken)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var f friend
	json.NewDecoder(r.Body).Decode(&f)
	friend_id := f.Id

	message_list := dbGetMessages(user_id, friend_id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message_list)
}

func home(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	authToken := c.Value
	user_id, err := ValidateToken(authToken)

	if err != nil {
		log.Println("invalid token")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "home.html")
	log.Print(user_id)
}
