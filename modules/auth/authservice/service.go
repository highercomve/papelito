package authservice

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/highercomve/papelito/modules/auth/authmodels"
	"github.com/highercomve/papelito/utils/prnx"
	"github.com/labstack/echo/v4"
)

type AccountToken struct {
	jwt.StandardClaims `json:",inline" bson:",inline"`

	Nick   string                 `json:"nick"`
	Scopes string                 `json:"scopes"`
	Prn    prnx.Prn               `json:"prn"`
	Roles  authmodels.RolesType   `json:"roles"`
	Type   authmodels.AccountType `json:"type"`
}

// GetOwnerID get user id from jwt of cookie
func GetOwnerID(c echo.Context) (prnx.Prn, error) {
	userRaw := c.Get("user")

	if userRaw == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "user is not loaded")
	}

	user := userRaw.(*authmodels.AccountToken)
	return user.Prn, nil
}

// GetUser get user from jwt of cookie
func GetUser(c echo.Context) (*authmodels.AccountToken, error) {
	userRaw := c.Get("user")

	if userRaw == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "user is not loaded")
	}

	user := userRaw.(*authmodels.AccountToken)
	return user, nil
}
