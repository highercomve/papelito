package authmodels

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/highercomve/papelito/modules/helper/helpermodels"
	"github.com/highercomve/papelito/utils/prnx"
)

const TokenServicePrn = "tokens"

// AccountToken Token payload
type AccountToken struct {
	jwt.StandardClaims `json:",inline" bson:",inline"`

	ID     string      `json:"id"`
	Nick   string      `json:"nick"`
	Scopes string      `json:"scopes"`
	Prn    prnx.Prn    `json:"prn"`
	Roles  RolesType   `json:"roles"`
	Type   AccountType `json:"type"`
}

// TokenPayload login token payload
type TokenPayload struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type,omitempty"`
	Scopes    string `json:"scopes,omitempty"`
}

// Token database token
type Token struct {
	helpermodels.Identificable `json:",inline" bson:",inline"`
	helpermodels.Timestamp     `json:",inline" bson:",inline"`
	helpermodels.Ownership     `json:",inline" bson:",inline"`

	ExpiresAt *time.Time  `json:"expires_at" bson:"expires_at"`
	Nick      string      `json:"nick" bson:"nick"`
	Scopes    string      `json:"scopes" bson:"scopes"`
	Roles     RolesType   `json:"roles" bson:"roles"`
	Type      AccountType `json:"type" bson:"type"`
	Raw       string      `json:"token" bson:"raw"`
}

func (t *Token) GetServicePrn() string {
	return TokenServicePrn
}
