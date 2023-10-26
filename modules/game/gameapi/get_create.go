package gameapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCreateGame create game page
func GetCreateGame(c echo.Context) error {
	return c.Render(http.StatusOK, "games/create.html", nil)
}
