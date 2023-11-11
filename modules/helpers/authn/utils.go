package authn

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func ErrorOrRedirect(contentType string, c echo.Context) error {
	err := Logout(c)
	if err != nil {
		if strings.Contains(contentType, "html") {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		}

		return echo.ErrUnauthorized
	}

	if strings.Contains(contentType, "html") {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
	}

	return echo.ErrUnauthorized
}

func Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}

	err = sess.Save(c.Request(), c.Response())

	return err
}
