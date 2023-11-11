package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/highercomve/papelito/modules/auth/authmodels"
	"github.com/highercomve/papelito/modules/helpers/authn"
	"github.com/highercomve/papelito/modules/helpers/prnx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func cookieToToken(c echo.Context) {
	sess, err := session.Get("session", c)
	if err != nil {
		return
	}

	token := ""
	if sess != nil {
		t, ok := sess.Values["token"]
		if ok && t != nil {
			token = t.(string)
		}
	}

	if token == "" {
		return
	}

	c.Request().Header.Set("Authorization", "Bearer "+token)
}

func putPrnIntoToken(c echo.Context) {
	token := c.Get("user").(*jwt.Token)

	if token != nil {
		claims := token.Claims.(jwt.MapClaims)
		_, ok := claims["prn"]
		if !ok {
			id := claims["id"]
			claims["prn"] = prnx.GetPrn(id.(string), "accounts")
			token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
			tokenString, err := token.SignedString(Env.JWTAuthSecret)
			if err == nil {
				c.Request().Header.Set("Authorization", "Bearer "+tokenString)
			}
		}
	}
}

// GetJWTConfig return jwt config
func GetJWTConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		SigningMethod:           "RS256",
		SigningKey:              &Env.JWTAuthSecret.PublicKey,
		Claims:                  jwt.MapClaims{},
		BeforeFunc:              cookieToToken,
		ErrorHandlerWithContext: errorOrRedirect,
		SuccessHandler:          putPrnIntoToken,
	}
}

func errorOrRedirect(err error, c echo.Context) error {
	contentType := c.Request().Header.Get("Accept")

	return authn.ErrorOrRedirect(contentType, c)
}

// CreateAccountToken create account token
func CreateAccountToken(acc *authmodels.Account, scopes []string) (*authmodels.TokenPayload, error) {
	timeoutStr := Env.JWTTimeoutMinutes
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return nil, err
	}

	claimScopes := "prn:papelito:apis:/api/all"
	if scopes != nil {
		claimScopes = ""
		for index, scope := range scopes {
			if index == 0 {
				claimScopes = scope
			}
			claimScopes = fmt.Sprintf("%s,%s", claimScopes, scope)
		}
	}

	claims := &authmodels.AccountToken{}
	claims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(timeout)).Unix()
	claims.ID = acc.ID
	claims.Nick = acc.Nick
	claims.Roles = acc.Role
	claims.Type = acc.Type
	claims.Prn = prnx.GetPrn(acc.Identification.ID, "accounts")
	claims.Scopes = claimScopes

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	tokenString, err := token.SignedString(Env.JWTAuthSecret)

	return &authmodels.TokenPayload{
		Token:     tokenString,
		TokenType: "bearer",
		Scopes:    claimScopes,
	}, err
}
