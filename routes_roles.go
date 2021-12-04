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
// I will not remove existing members but role won't be accessible anymore
func DeleteRole(c echo.Context) error {
	key := c.Param("key")
	role := c.Param("role")

	// check if token is valid
	message, data := CheckIfExist(key)
	if data == "" {
		fmt.Print("\n")
		fmt.Print(message)
		return c.JSON(http.StatusOK, returnData("Invalid login", ""))
	}

	// Check if user owns the role
	var UserRole string
	db.QueryRow("SELECT rolename FROM useraccounts.roles WHERE email=$1", data).Scan(&UserRole)

	// Is the role correct?
	if !strings.EqualFold(UserRole, role) {
		fmt.Print(UserRole, "\r", role)
		return c.JSON(http.StatusOK, returnData("Role does not exists or is not owned by you", ""))
	}

	// Delete from database delete nothing if role is not found under owners name/email
	res, _ := db.Exec("DELETE FROM useraccounts.roles WHERE rolename=$1", role)

	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Could not delete role, try agen later", ""))
	}

	return c.JSON(http.StatusOK, returnData("Role has been deleted", ""))
}
