package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Create a session key
func GenerateSessionKey(email string) string {
	RandomSeed := rand.Intn(1423)
	BcryptToken, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprint(RandomSeed*int(time.Nanosecond))), bcrypt.DefaultCost)
	RandomKey := string(BcryptToken)

	// Add session key to database
	res, _ := db.Exec("INSERT INTO useraccounts.sessions VALUES($1, $2, $3) ON CONFLICT DO NOTHING", RandomKey, time.Now().String(), email)

	// If data exists then generate a new session key and update existing
	dataRecived, _ := res.RowsAffected()
	if dataRecived < 1 {
		db.Exec("UPDATE useraccounts.sessions SET sessiontoken=$1, added=$2 WHERE idec=$3", RandomKey, time.Now().String(), email)
	}

	return RandomKey
}

// Check if token has not expired
func CheckIfExist(Email string) {
	rows, err := db.Query("")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	var TimeandKey []string
	for rows.Next() {
		if err := rows.Scan(&TimeandKey); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}

	fmt.Print(TimeandKey)
}
