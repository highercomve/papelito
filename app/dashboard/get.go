package dashboard

import (
	"github.com/highercomve/papelito/modules/game/gameapi"
	"github.com/labstack/echo/v4"
)

// GetDashboard get user dashboard
func GetDashboard(c echo.Context) error {
	return gameapi.GetCreateGame(c)
}
