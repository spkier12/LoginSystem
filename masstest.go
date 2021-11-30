package main

import (
	"fmt"
	"math/rand"
	"time"
)

var spam int = 666

func CreateAccountTEST() {
	spam++
	fmt.Print("Starting mass test...")
	for spam < 900000 {
		elapsed := time.Now()
		db.Exec("INSERT INTO useraccounts.useracc VALUES($1, $2, $3, $4, $5)", spam*len(time.Now().String())*rand.Intn(12323), time.Now().String()+string(spam), time.Now().String(), "TEST", time.Now().String()+string(spam*123))
		sub := time.Since(elapsed)
		fmt.Print("\nIt took: " + sub.String() + " To create new user accounts")
	}
}
