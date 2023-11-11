package middlewares

import (
	"strings"

	"github.com/highercomve/papelito/modules/helpers"
	"github.com/labstack/echo/v4"
)

// ValidateAuthenticity validate authenticity on post
func ValidateAuthenticity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isForm := strings.Contains(c.Request().Header.Get(echo.HeaderContentType), "form")

		if isForm && c.Request().Method == echo.POST {
			authenticityToken := c.FormValue("authenticity_token")
			err := helpers.ValidateAuthenticityToken(authenticityToken)
			if err != nil {
				return echo.NewHTTPError(400, "Form authenticiy is invalid")
			}
		}

		return next(c)
	}
}
