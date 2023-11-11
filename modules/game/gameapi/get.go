package gameapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetGame create game page
func GetGame(c echo.Context) error {
	id := c.Param("id")

	game, err := gameService.GetGame(c.Request().Context(), id, nil)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "games/game.html", game)
}
