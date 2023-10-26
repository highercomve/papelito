package gameapi

import (
	"fmt"
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

	fmt.Println(*payload)
	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	game, err := gameService.CreateGame(payload)
	if err != nil {
		return c.Render(http.StatusBadRequest, "games/create.html", err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/games/"+game.ID)
}
