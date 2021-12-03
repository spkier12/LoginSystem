package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

// Create a session key
func GenerateSessionKey(email string) string {
	RandomKey2 := base64.StdEncoding.EncodeToString(sha256.New().Sum([]byte(fmt.Sprint(rand.Intn(123543)))))

	// Get date
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()

	// Add session key to database
	res, _ := db.Exec("INSERT INTO useraccounts.sessions VALUES($1, $2, $3) ON CONFLICT DO NOTHING", RandomKey2, fmt.Sprint(year)+" "+fmt.Sprint(month)+" "+fmt.Sprint(day), email)

	// If data exists then generate a new session key and update existing
	dataRecived, _ := res.RowsAffected()
	if dataRecived < 1 {
		db.Exec("UPDATE useraccounts.sessions SET sessiontoken=$1, added=$2 WHERE idec=$3", RandomKey2, fmt.Sprint(year)+" "+fmt.Sprint(month)+" "+fmt.Sprint(day), email)
	}

	return RandomKey2
}

// Check if token has not expired
func CheckIfExist(key string) string {

	// Get date
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()

	rows, err := db.Query("SELECT idec FROM useraccounts.sessions WHERE sessiontoken=$1", key)
	if err != nil {
		fmt.Print("There was a error connecting to database...")
		return ""
	}

	var email string
	for rows.Next() {
		if err := rows.Scan(&email); err != nil {
			return ""
		}
	}

	// check if date is correct
	if email != "" {
		rows, _ := db.Query("SELECT added FROM useraccounts.sessions WHERE sessiontoken=$1", key)

		var timer string
		for rows.Next() {
			if err := rows.Scan(&timer); err != nil {
				return ""
			}
		}

		if timer == fmt.Sprint(year)+" "+fmt.Sprint(month)+" "+fmt.Sprint(day) {
			return email
		}
	}

	return "Here is your email but token is invalid\r" + email
}
