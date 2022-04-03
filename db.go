package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var connection_string = os.Getenv("CONN_STRING")

func OpendbConnection() *sql.DB {
	db, err := sql.Open("mysql", connection_string)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func dbAddUser(username string, password string) bool {
	db := OpendbConnection()

	row := db.QueryRow("SELECT EXISTS(SELECT * FROM user WHERE username = ?)", username)
	var exist bool
	err := row.Scan(&exist)

	if err != nil {
		return false
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if !exist {
		db.Query("INSERT INTO user (username, password_hash) VALUES (?, ?)", username, password_hash)
		db.Query("INSERT INTO authentication_tokens (user_id) VALUES ((SELECT id FROM user WHERE username = ?));", username)
	}

	defer db.Close()

	return !exist
}

func dbAddFriend(user_id int, friend_name string) bool {
	db := OpendbConnection()

	row := db.QueryRow("SELECT id FROM user WHERE username = ?;", friend_name)
	var friend_id int
	err := row.Scan(&friend_id)
	if err != nil {
		return false
	}

	if user_id == friend_id {
		return false
	}

	row = db.QueryRow("SELECT EXISTS(SELECT * FROM friend_requests WHERE user_id = ? and friend_id = ?)", friend_id, user_id)
	exist := false
	err = row.Scan(&exist)
	defer db.Close()

	if exist {
		_, err = db.Query(`
			DELETE FROM friend_requests 
			WHERE user_id = ? and friend_id = ?;`,
			friend_id, user_id)

		_, err = db.Query(`INSERT INTO friend (user_id, friend_id) 
			VALUES (?, ?), (?, ?);`,
			friend_id, user_id, user_id, friend_id)

		username := ""
		db.QueryRow("SELECT username FROM user where id = ?", user_id).Scan(&username)

		room_name := fmt.Sprintf("%s, %s", username, friend_name)
		var room_id int
		db.Query("INSERT INTO room (name) value (?)", room_name)
		db.QueryRow("SELECT id FROM room WHERE name = ?", room_name).Scan(&room_id)

		db.Query("INSERT INTO participants (user_id, room_id) VALUES (?, ?), (?, ?);", user_id, room_id, friend_id, room_id)

		return true
	} else {
		_, err = db.Query("INSERT INTO friend_requests (user_id, friend_id) VALUES (?, ?);", user_id, friend_id)
		log.Println(err)
		return true
	}

}

type friend struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func dbFriendsList(user_id int) []friend {

	var friends_list []friend

	db := OpendbConnection()

	rows, err := db.Query(`
		SELECT
		    user.id,
		    user.username
		FROM friend
		LEFT JOIN user ON friend.friend_id = user.id
		WHERE friend.user_id = ?;`, user_id)

	if err != nil {
		log.Print(err)
		return nil
	}

	for rows.Next() {
		var f friend
		if err := rows.Scan(&f.Id, &f.Username); err != nil {
			log.Print(err)
			return nil
		}
		friends_list = append(friends_list, f)
	}

	defer db.Close()
	return friends_list
}
