package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type jwtAuthenticator struct {
	secret   string
	audience string
	iss      string
}

func NewJWTAuthenticator(secret string, audience string, iss string) *jwtAuthenticator {
	return &jwtAuthenticator{
		secret:   secret,
		audience: audience,
		iss:      iss,
	}
}

func (a *jwtAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *jwtAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %s", t.Header["alg"])
		}
		return []byte(a.secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
}
