package app

import (
	"github.com/highercomve/papelito/modules/auth/authapi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoadApp load application
func LoadApp(e *echo.Echo) *echo.Group {
	app := e.Group("")
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	authapi.Load(app)

	return app
}
