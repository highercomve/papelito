package gameapi

import (
	"fmt"

	"github.com/highercomve/papelito/modules/game/gameservice"
	"github.com/labstack/echo/v4"
)

var gameService = gameservice.NewGameMachine()

// Load Create new auth service
func Load(e *echo.Group) *echo.Group {
	g := e.Group("/games")

	g.GET("", GetCreateGame)
	g.POST("", CreateGame)
	g.GET("/:id", GetGame)
	g.PUT("/:id", UpdateGame)

	fmt.Printf("%++v", g)
	return g
}
