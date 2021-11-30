package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Initialize database
func InitDB() (*sql.DB, error) {

	// Connection to database
	const (
		host     = "80.208.226.78"
		port     = "44320"
		user     = "zerorootzero"
		password = "zerorootzerozeroroorootzeroninehundredFiftyFive905327895642905327895642905327895642905327895642905327895642"
		dbname   = "zerorootzero"
	)

	con := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", con)
	return db, err
}

func main() {
	DB, err := InitDB()
	db = DB

	// Check if DB has errors
	if err != nil {
		fmt.Print("Database error exiting..")
		os.Exit(1)
	} else {
		fmt.Print("Database up and running.. \n")
	}

	// Mas creation test towards DB
	//CreateAccountTEST()

	e := echo.New()
	e.GET("/Create/:email/:password", Create)
	e.GET("/EnableMFA/:email/:key", EnableMFA)
	e.GET("/Login/:email/:key/:id", Login)
	e.Start(":5001")
}
