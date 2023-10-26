package gameapi

import (
	"net/http"

	gamemachine "github.com/highercomve/papelito/modules/game/gameservice"
	"github.com/labstack/echo/v4"
)

// CreateGame create game page
func CreateGame(c echo.Context) error {
	payload := new(gamemachine.Configuration)
	if err := c.Bind(payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.Render(http.StatusOK, "games/game.html", payload)
}
