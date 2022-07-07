package db

import (
	"fmt"
	"log"
)

// Add user to database
func AddUser(username string, password_hash string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("Adding user ", username)

	// check if user already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("User already exists")
		return fmt.Errorf("User already exists")
	}

	// add user to database
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, password_hash)
	if err != nil {
		return err
	}

	log.Println("User added")

	return nil
}
