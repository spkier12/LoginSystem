package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	test()
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
	e.Use(middleware.CORS())

	// User management
	e.GET("/Create/:email/:password", Create)  // Create a new user account unless it exists
	e.GET("/EnableMFA/:email/:key", EnableMFA) // Enable MFA so that user can login, this is mandatory.
	e.GET("/Login/:email/:key/:id", Login)     // Login to the user account and get a token in return Valid until next day
	e.GET("/Validate/:id", Validate)           // Check if token is valid

	// Role management
	e.GET("/Createrole/:role/:key", CreateRole) // Create a role where the user who creates it is the owner
	e.GET("/Deleterole/:role/:key", DeleteRole) // Delete a role if owner owns it
	e.Start(":5001")                            // Start the server
}

// Test response time
func test() {
	a := time.Now()
	pass, _ := bcrypt.GenerateFromPassword([]byte("password"), 5)
	b := time.Now()
	bcrypt.CompareHashAndPassword(pass, []byte("password"))
	c := time.Now()
	encrypt := b.Sub(a)
	decrypt := c.Sub(b)
	fmt.Println("encrypt: ", encrypt)
	fmt.Println("decrypt: ", decrypt)
}
