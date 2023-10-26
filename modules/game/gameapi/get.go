package gameapi

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetGame create game page
func GetGame(c echo.Context) error {
	id := c.Param("id")

	fmt.Printf("%++v", gameService.Games)
	game, ok := gameService.Games[id]
	if !ok {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "games/game.html", game)
}
