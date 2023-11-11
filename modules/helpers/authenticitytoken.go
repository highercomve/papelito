package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// CreateAuthenticityToken create authenticity token for forms
func CreateAuthenticityToken() (string, error) {
	JWTTimeoutMinutes, _ := strconv.Atoi(Env.AuthenticityTokenTTLMinutes)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().UTC().Add(time.Duration(JWTTimeoutMinutes) * time.Minute).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(Env.JWTSigningSecret))
}

// ValidateAuthenticityToken validate authenticity token for forms
func ValidateAuthenticityToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(Env.JWTSigningSecret), nil
	})

	if err != nil {
		return err
	}

	return token.Claims.Valid()
}
