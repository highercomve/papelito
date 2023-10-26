package gameapi

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

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
