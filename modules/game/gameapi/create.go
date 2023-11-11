package gameapi

import (
	"net/http"

	"github.com/highercomve/papelito/modules/game/gamemodels"
	"github.com/labstack/echo/v4"
)

// CreateGame create game page
func CreateGame(c echo.Context) error {
	payload := new(gamemodels.Configuration)
	if err := c.Bind(payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	game, err := gameService.CreateGame(c.Request().Context(), payload)
	if err != nil {
		return c.Render(http.StatusBadRequest, "games/create.html", err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/games/"+game.ID)
}
