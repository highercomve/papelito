package gameapi

import (
	"net/http"

	"github.com/highercomve/papelito/modules/game/gamemodels"
	"github.com/labstack/echo/v4"
)

// UpdateGame create game page
func UpdateGame(c echo.Context) error {
	payload := new(gamemodels.Configuration)
	if err := c.Bind(payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.Render(http.StatusOK, "games/game.html", payload)
}
