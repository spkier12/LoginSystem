package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// Random Seed used to create UID at a later date
var storedUser int = rand.Intn(9000)

// Create a new user account in database and check if not exists
func Create(c echo.Context) error {
	storedUser++
	Email := c.Param("email")
	Pass := c.Param("password")

	// Generate Authenticator keys
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "UB-Systems",
		AccountName: Email,
	})

	// Create password
	hash, _ := bcrypt.GenerateFromPassword([]byte(Pass), 5)

	// Create user account
	res, _ := db.Exec("INSERT INTO useraccounts.useracc VALUES($1, $2, $3, $4, $5) ON CONFLICT (email) DO NOTHING", storedUser*len(Email), string(hash), key.Secret()+"-"+fmt.Sprint(storedUser*len(Email)*900), "NO", Email)
	if d, _ := res.RowsAffected(); d < 1 {
		return c.String(http.StatusConflict, "-")
	}

	// Create QR code
	image, _ := key.Image(128, 128)
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, image, nil)
	return c.String(http.StatusCreated, buf.String())
}

// If MFA is disabled then account is locked and signin will return wrong username/password
// To enable please scan QR code and write the key back to use
func EnableMFA(c echo.Context) error {
	// Get the values from parameter
	Email := c.Param("email")
	Key := c.Param("key")

	// Get the Key in database
	var keysfound string
	db.QueryRow("SELECT mfakey FROM useraccounts.useracc WHERE email = $1", Email).Scan(&keysfound)

	// Verify if key is valid and update the database with the correct value for mfaenabled
	if totp.Validate(Key, strings.Split(keysfound, "-")[0]) {
		res, err := db.Exec("UPDATE useraccounts.useracc SET mfaenabled = 'YES' WHERE email=$1", Email)
		if err != nil {
			return c.JSON(http.StatusLocked, returnData("Error in database try agen later!", ""))
		} else if d, _ := res.RowsAffected(); d > 0 {
			return c.JSON(http.StatusOK, returnData("MFA Is now enabled and you can login to your account!", ""))
		}
	}

	// If code is not valid
	return c.JSON(http.StatusNotAcceptable, returnData("Invalid Auth code", ""))
}

// Check if user login is valid and create a session token valid for 24 hours
func Login(c echo.Context) error {

	// Get the values from parameter
	Email := c.Param("email")
	Key := c.Param("key")
	pass := c.Param("id")

	// Get the Key in database
	var keysfound string
	var passfromdb string
	var mfaEnabled2 string

	db.QueryRow("SELECT mfakey, uid, mfaenabled FROM useraccounts.useracc WHERE email = $1", Email).Scan(&keysfound, &passfromdb, &mfaEnabled2)

	// If mfaenabled == "NO" then don't check for key
	if mfaEnabled2 == "YES" {

		// Verify if key is valid and update the database with the correct value for mfaenabled
		if !totp.Validate(Key, strings.Split(keysfound, "-")[0]) {
			return c.JSON(http.StatusLocked, returnData("Login failed", ""))
		}
	}

	// time from function start
	timestart := time.Now()

	// Check agenst hashed password in db and generate a session token if valid
	if err := bcrypt.CompareHashAndPassword([]byte(passfromdb), []byte(pass)); err != nil {
		return c.JSON(http.StatusForbidden, returnData("Login failed", ""))
	}

	// Time elapsed
	fmt.Print("\nLogin took: \r" + time.Since(timestart).String())

	// Everything went fine then continue
	return c.JSON(http.StatusOK, returnData("Login OK", GenerateSessionKey(Email)))
}

// Check if sessionkey is valid
func Validate(c echo.Context) error {
	key := c.Param("id")
	return c.JSON(http.StatusAccepted, returnData(CheckIfExist(key)))
}

// Easy function to generate data in return-
func returnData(message string, data string) string {
	type MyData struct {
		Message string
		Data    string
	}
	var mydata MyData
	mydata.Message = message
	mydata.Data = data
	d, _ := json.Marshal(mydata)
	return string(d)
}
