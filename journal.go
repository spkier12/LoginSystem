package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create new journal in database
func CreateJournal(c echo.Context) error {
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusForbidden, returnData("Login failed", "0"))
	}

	err1 := UserHasRole(token.Value, "udips6")
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Login failed", fmt.Sprint(err1)))
	}

	return c.JSON(http.StatusOK, returnData("Journal var opprettet", "1"))
}
