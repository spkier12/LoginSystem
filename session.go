package main

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

// Create a session key
func GenerateSessionKey() string {
	RandomSeed := rand.Intn(123445323423423)
	BcryptToken, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprint(RandomSeed)), 62)
	RandomKey := string(BcryptToken)
	return RandomKey
}

// Check if There is a key within current date and time in database before generating new key
func CheckIfExist() {

}
