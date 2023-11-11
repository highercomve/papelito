package helpers

import (
	"crypto/rsa"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// Method define methods for encrypt supported
type Method string

const (
	// BCryptMethod method
	BCryptMethod Method = "bcrypt"
	// SCryptMethod method
	SCryptMethod Method = "scrypt"
)

type crytoMethods struct {
	BCrypt Method
	SCrypt Method
}

// CryptoMethods kind a enum for cryptography methods supported
var (
	CryptoMethods = &crytoMethods{
		BCrypt: BCryptMethod,
		SCrypt: SCryptMethod,
	}
	errMethodNotFound = errors.New("the only encrypt method supported are bcrypt and scrypt")
)

// JwtRsaKeys Public and Private keys for Jwt
type JwtRsaKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// HashPassword create a hashed version of a string
func HashPassword(password string, method Method) (string, error) {
	switch method {
	case BCryptMethod:
		return bcryptHashPassword(password)
	case SCryptMethod:
		return scryptHashPassword(password)
	default:
		return "", errMethodNotFound
	}
}

// CheckPasswordHash validate password agains a given hash
func CheckPasswordHash(password, hash string, method Method) bool {
	switch method {
	case BCryptMethod:
		return bcryptCheckPasswordHash(password, hash)
	case SCryptMethod:
		return scryptCheckPasswordHash(password, hash)
	default:
		return false
	}
}

func bcryptHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func bcryptCheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func scryptHashPassword(password string) (string, error) {
	bytes, err := scrypt.Key([]byte(password), []byte(Env.ScryptSecret), 32768, 8, 1, 32)
	return hex.EncodeToString(bytes), err
}

func scryptCheckPasswordHash(password, hash string) bool {
	passwordHash, err := scryptHashPassword(password)
	return err == nil && passwordHash == hash
}
