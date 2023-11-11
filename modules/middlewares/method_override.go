package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// MethodOverride validate authenticity on post
func MethodOverride(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.Request()
		isForm := strings.Contains(request.Header.Get(echo.HeaderContentType), "form")

		if isForm && request.Method == echo.POST {
			newMethod := c.FormValue("_method")
			if newMethod != "" {
				request.Method = newMethod
			}
		}

		return next(c)
	}
}
