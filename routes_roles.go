package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Create a role in the database
func CreateRole(c echo.Context) error {
	key := c.Param("key")
	role := c.Param("role")

	// Check if token is valid
	message, data := CheckIfExist(key)
	if data == "" {
		fmt.Print("\n")
		fmt.Print(message)
		return c.JSON(http.StatusOK, returnData("Invalid login", ""))
	}

	// Create role and if the role exist then do nothing and return back to user with error
	res, err := db.Exec("INSERT INTO useraccounts.roles VALUES ($1, 'YES', $2) ON CONFLICT DO NOTHING", role, data)
	if err != nil {
		fmt.Print(err)
		return err
	}
	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusConflict, returnData("Role exists", ""))
	}

	return c.JSON(http.StatusOK, returnData("Role Created!", ""))
}

// Delete the role from database
// It will not remove existing members but role won't be accessible anymore
func DeleteRole(c echo.Context) error {
	key := c.Param("key")
	role := c.Param("role")

	// check if token is valid
	err := UserHasRole(key, role)
	if err != nil {
		return c.JSON(http.StatusOK, returnData("Could not delete role\nAre you the owner?", ""))
	}

	// Delete from database delete nothing if role is not found under owners name/email
	res, _ := db.Exec("DELETE FROM useraccounts.roles WHERE rolename=$1", role)

	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Could not delete role, try agen later", ""))
	}

	return c.JSON(http.StatusOK, returnData("Role has been deleted", ""))
}

// Check if user has a certain role
func UserHasRole(key string, role string) error {

	// check if token is valid
	message, data := CheckIfExist(key)
	if data == "" {
		fmt.Print("\n")
		fmt.Print(message)
		return fmt.Errorf("key invalid")
	}

	// Check if user owns the role
	var UserRole string
	db.QueryRow("SELECT rolename FROM useraccounts.roles WHERE email=$1", data).Scan(&UserRole)

	// Is the role correct?
	if !strings.EqualFold(UserRole, role) {
		fmt.Print(UserRole, "\r", role)
		return fmt.Errorf("role is not owned by you")
	}

	return nil
}

// If the user is the owner of role then he can invite peolpe to the role either my email.
// If user was invited then he has access to the website that requires the role.
func InviteRole(c echo.Context) error {
	key := c.Param("key")
	email := c.Param("email")
	role := c.Param("role")

	// Check if key is valid
	// Check if user owns the role
	err := UserHasRole(key, role)
	if err != nil {
		return c.JSON(http.StatusLocked, returnData("You are not the owner of the role or your login has expired", ""))
	}

	sql, err := db.Exec("INSERT INTO useraccounts.invites VALUES ($1, $2)", email, role)
	if d, _ := sql.RowsAffected(); d < 1 {
		return c.JSON(http.StatusLocked, returnData("Error was risen trying to invite user", ""))
	}

	if err != nil {
		return c.JSON(http.StatusLocked, returnData("Error was risen trying to invite user", ""))
	}

	return c.JSON(http.StatusOK, returnData("User was invited to role/group", ""))

}

// Check what role user have
func CheckRole(c echo.Context) error {
	email := c.Param("email")
	role := c.Param("role")

	dbrole := ""
	db.QueryRow("SELECT rolename FROM useraccounts.roles WHERE email=$1", email).Scan(&dbrole)

	if strings.EqualFold(role, dbrole) {
		return c.JSON(http.StatusOK, returnData("User has role", "OK"))
	}
	return c.JSON(http.StatusOK, returnData("User does not have role", "NO"))
}
