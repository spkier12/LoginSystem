package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Create a role in the database
func CreateRole(c echo.Context) error {
	key, err := c.Cookie("token")
	if err != nil {
		return c.JSON(http.StatusAccepted, returnData("Invalid login", "0"))
	}
	role := c.Param("role")

	// Check if token is valid
	_, data := CheckIfExist(key.Value)
	if data == "" {
		return c.JSON(http.StatusOK, returnData("Invalid login", "0"))
	}

	// Create role and if the role exist then do nothing and return back to user with error
	res, err := db.Exec("INSERT INTO useraccounts.roles VALUES ($1, 'YES', $2) ON CONFLICT DO NOTHING", role, data)
	if err != nil {
		return err
	}
	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusConflict, returnData("Role exists", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Role Created!", "1"))
}

// Delete the role from database
// It will not remove existing members but role won't be accessible anymore
func DeleteRole(c echo.Context) error {
	key0, err := c.Cookie("token")
	role := c.Param("role")

	if err != nil {
		return c.JSON(http.StatusAccepted, returnData("Invalid login", "0"))
	}

	// check if token is valid
	err0 := UserHasRole(key0.Value, role, true)
	if err0 != nil {
		fmt.Print(err0)
		return c.JSON(http.StatusOK, returnData("Could not delete role\nAre you the owner?", "0"))
	}

	// Delete from database delete nothing if role is not found under owners name/email
	res, _ := db.Exec("DELETE FROM useraccounts.roles WHERE rolename=$1", role)
	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Could not delete role, try agen later", "0"))
	}

	db.Exec("DELETE FROM useraccounts.invites WHERE rolename=$1", role)
	return c.JSON(http.StatusOK, returnData("Role has been deleted", "1"))
}

// Check if user has a certain role
// owned is true if you wanna check if user owns the role and false if you just wanna check if he has the role
func UserHasRole(key string, role string, owned bool) error {
	// check if token is valid
	message, data := CheckIfExist(key)
	if data == "" {
		fmt.Print("\n")
		fmt.Print(message)
		return fmt.Errorf("key invalid")
	}

	if owned {

		// Check if user owns the role
		var UserRole string
		db.QueryRow("SELECT rolename FROM useraccounts.roles WHERE email=$1 AND rolename=$2", data, role).Scan(&UserRole)

		// Is the role correct?
		if !strings.EqualFold(UserRole, role) {
			fmt.Print(UserRole + role)
			return fmt.Errorf("role is not owned by you")
		}
	}

	fmt.Print("Checking if user has role")
	var UserRole string
	db.QueryRow("SELECT rolename FROM useraccounts.invites WHERE email=$1 AND rolename=$2", data, role).Scan(&UserRole)

	var UserRole2 string
	db.QueryRow("SELECT rolename FROM useraccounts.roles WHERE email=$1 AND rolename=$2", data, role).Scan(&UserRole)

	// Is the role correct?
	if strings.EqualFold(UserRole, role) {
		return nil
	}

	// Is the role correct?
	if strings.EqualFold(UserRole2, role) {
		return nil
	}
	return fmt.Errorf("you are not a part of this role")
}

// If the user is the owner of role then he can invite peolpe to the role either my email.
// If user was invited then he has access to the website that requires the role.
func InviteRole(c echo.Context) error {
	key, err1 := c.Cookie("token")
	if err1 != nil {
		return c.JSON(http.StatusAccepted, returnData("Invalid login", "0"))
	}

	email := c.Param("email")
	role := c.Param("role")

	// Check if key is valid
	// Check if user owns the role
	err := UserHasRole(key.Value, role, false)
	if err != nil {
		return c.JSON(http.StatusLocked, returnData("You are not a part of this role", "0"))
	}

	sql, err := db.Exec("INSERT INTO useraccounts.invites VALUES ($1, $2)", email, role)
	if d, _ := sql.RowsAffected(); d < 1 {
		return c.JSON(http.StatusLocked, returnData("Error was risen trying to invite user", "0"))
	}

	if err != nil {
		return c.JSON(http.StatusLocked, returnData("Error was risen trying to invite user", "0"))
	}

	return c.JSON(http.StatusOK, returnData("User was invited to role/group", "1"))

}

// You can kick a specific user from the role he is in as long as you are the owner of the role.
func Kickfromrole() {

}

// Check what role user have
func CheckRole(c echo.Context) error {
	email := c.Param("email")
	role := c.Param("role")

	dbrole := ""
	db.QueryRow("SELECT rolename FROM useraccounts.invites WHERE email=$1 AND rolename=$2", email, role).Scan(&dbrole)

	if strings.EqualFold(role, dbrole) {
		return c.JSON(http.StatusOK, returnData("User has role", "1"))
	}
	return c.JSON(http.StatusOK, returnData("User does not have role", "0"))
}
