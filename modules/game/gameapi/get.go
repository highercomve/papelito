package gameapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetGame create game page
func GetGame(c echo.Context) error {
	id := c.Param("id")

	return c.Render(http.StatusOK, "games/game.html", map[string]string{"Id": id})
}
