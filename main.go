package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	DB, _ := InitDB()
	db = DB

	err2 := db.Ping()
	if err2 != nil {
		fmt.Print("Database failed ping test...")
	}

	// Mas creation test towards DB
	//CreateAccountTEST()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// User management
	e.GET("/Create/:email/:password", Create)  // Create a new user account unless it exists
	e.GET("/EnableMFA/:email/:key", EnableMFA) // Enable MFA so that user can login, this is mandatory.
	e.GET("/Login/:email/:key/:id", Login)     // Login to the user account and get a token in return Valid until next day
	e.GET("/Validate", Validate)               // Check if token is valid

	// Role management
	e.GET("/Createrole/:role/:key", CreateRole)        // Create a role where the user who creates it is the owner
	e.GET("/Deleterole/:role/:key", DeleteRole)        // Delete a role if owner owns it
	e.GET("/Inviterole/:role/:email/:key", InviteRole) // Invite the user to join role
	e.GET("/Checkrole/:role/:email", CheckRole)        // Check if user is part of a role
	e.Start(":5001")                                   // Start the server
}
