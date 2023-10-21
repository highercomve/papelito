package authapi

import (
	"github.com/labstack/echo/v4"
)

// Load Create new auth service
func Load(e *echo.Group) *echo.Group {
	g := e.Group("auth")

	g.GET("/login", LoginPage)
	g.POST("/login", Authorize)
	g.GET("/logout", CloseSession)
	g.GET("/register", RegisterPage)

	return g
}
