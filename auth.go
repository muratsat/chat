package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const timeLayout = "2006-01-02 15:04:05"

func CheckCredentials(username string, password string) bool {
	db := OpendbConnection()

	row := db.QueryRow("SELECT password_hash FROM user WHERE username = ?", username)
	var password_hash string
	err := row.Scan(&password_hash)

	if err != nil {
		return false
	}

	defer db.Close()

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	return err == nil
}

func GenerateToken(username string) (map[string]interface{}, error) {
	db := OpendbConnection()
	row := db.QueryRow("SELECT id FROM user WHERE username = ?", username)
	var user_id int
	err := row.Scan(&user_id)
	if err != nil {
		log.Print(1)
		return nil, err
	}

	randomToken := make([]byte, 32)
	_, err = rand.Read(randomToken)
	if err != nil {
		log.Print(2)
		return nil, err
	}

	authToken := base64.URLEncoding.EncodeToString(randomToken)
	dt := time.Now()
	exipiryTime := dt.Add(time.Minute * 60)

	generatedAt := dt.Format(timeLayout)
	expiresAt := exipiryTime.Format(timeLayout)

	tokenDetails := map[string]interface{}{
		"token_type":   "Bearer",
		"auth_token":   authToken,
		"generated_at": generatedAt,
		"expires_at":   expiresAt,
	}

	log.Print(tokenDetails)
	defer db.Close()
	return tokenDetails, nil
}

func UpdateToken(username string) (map[string]interface{}, error) {
	token, err := GenerateToken(username)
	if err != nil {
		return nil, err
	}

	db := OpendbConnection()
	db.Query("UPDATE authentication_tokens SET generated_at = ?, expires_at = ?, auth_token = ? WHERE user_id = (SELECT id FROM user WHERE username = ?);", token["generated_at"], token["expires_at"], token["auth_token"], username)

	defer db.Close()
	return token, nil
}

func ValidateToken(token string) (int, error) {
	db := OpendbConnection()

	row := db.QueryRow(`SELECT
    	user.id,
    	username,
    	generated_at,
    	expires_at
	FROM authentication_tokens
	LEFT JOIN user
	ON authentication_tokens.user_id = user.id
    WHERE auth_token = ?`, token)

	defer db.Close()

	var user_id int
	var username string
	var generated_at string
	var expires_at string

	err := row.Scan(&user_id, &username, &generated_at, &expires_at)
	if err != nil {
		// invalid token
		return 0, errors.New("invalid token")
	}

	expiry_time, _ := time.Parse(timeLayout, expires_at)
	current_time, _ := time.Parse(timeLayout, time.Now().Format(timeLayout))

	if expiry_time.Before(current_time) {
		// token has expired
		return 0, errors.New("token has expired")
	}

	return user_id, nil
}
