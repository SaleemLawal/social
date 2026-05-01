package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type Authenticator interface {
	GenerateToken(jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
