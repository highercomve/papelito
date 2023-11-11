package middlewares

import (
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/highercomve/papelito/modules/helpers/authn"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const desactivateAuthentication = true

// CookieAuthentication middleware
func CookieAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if desactivateAuthentication {
			return next(c)
		}
		token := strings.TrimSpace(strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1))
		if token != "" {
			return next(c)
		}
		contentType := c.Request().Header.Get("Accept")

		sess, err := session.Get("session", c)
		if err != nil {
			return authn.ErrorOrRedirect(contentType, c)
		}

		if sess != nil {
			t, ok := sess.Values["token"]
			if ok && t != nil {
				token = t.(string)
			}
		}

		if token == "" {
			return authn.ErrorOrRedirect(contentType, c)
		}

		if err := validateToken(token); err != nil {
			return authn.ErrorOrRedirect(contentType, c)
		}

		if _, ok := sess.Values["id"]; !ok {
			return authn.ErrorOrRedirect(contentType, c)
		}

		return next(c)
	}
}

func getTokenClaims(tokenString string) (jwt.Claims, error) {
	var p = new(jwt.Parser)
	token, _, err := p.ParseUnverified(tokenString, &jwt.StandardClaims{})

	return token.Claims, err
}

func validateToken(tokenString string) error {
	claims, err := getTokenClaims(tokenString)
	if err != nil {
		return err
	}
	return claims.Valid()
}
