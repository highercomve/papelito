package authapi

import (
	"net/http"

	"github.com/highercomve/papelito/modules/helpers/authn"
	echo "github.com/labstack/echo/v4"
)

// CloseSession Remove session variables a logout user
func CloseSession(c echo.Context) error {
	err := authn.Logout(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't save session cookie")
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}
