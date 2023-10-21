package authapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterPage login page
func RegisterPage(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/register.html", nil)
}
