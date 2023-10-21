package authapi

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/highercomve/papelito/modules/auth/authmodels"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// ReqPayload payload for authentication request
type ReqPayload struct {
	User     string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// Authorize process login form
func Authorize(c echo.Context) error {
	payload := new(ReqPayload)
	if err := c.Bind(payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// token, acc, err := authservice.LoginWithPass(c.Request().Context(), payload.User, payload.Password)
	// if err != nil {
	// 	return err
	// }

	token := authmodels.TokenPayload{}
	acc := authmodels.Account{}

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
	}

	sess.Values["token"] = token.Token
	sess.Values["nick"] = acc.Nick
	sess.Values["id"] = acc.ID
	sess.Values["prn"] = acc.PRN

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't save session cookie")
	}

	return c.Redirect(http.StatusFound, "/dashboard")
}
