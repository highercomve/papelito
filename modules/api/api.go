package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoadAPI load all API
func LoadAPI(e *echo.Echo) *echo.Group {
	api := e.Group("api/v1")
	api.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	return api
}
