package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")
)

const (
	header     string = "Authorization"
	authScheme string = "Bearer"
)

// ReadRequestToken authenticate token
func ReadRequestToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := GetToken("token_", c)
			if err != nil {
				return err
			}

			// c.Set("user", user)
			c.Set("token", token)
			return next(c)
		}
	}
}

// GetToken
func GetToken(tokenKey string, c echo.Context) (string, error) {
	authCookie, err := c.Request().Cookie(tokenKey)
	token := ""
	if err != nil && err != http.ErrNoCookie {
		return "", ErrJWTMissing
	}
	if err == nil && authCookie.Value != "" {
		token = authCookie.Value
	}

	if token == "" {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if auth == "" || len(auth) < l+1 || auth[:l] != authScheme {
			return "", ErrJWTMissing
		}
		token = auth[l+1:]
	}

	if token == "" {
		return "", ErrJWTMissing
	}

	return token, nil
}
