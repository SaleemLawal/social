package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TestAuthenticator struct{}

const testSecret = "test-secret"

var testClaims = jwt.MapClaims{
	"sub": int64(6),
	"exp": time.Now().Add(1 * time.Hour).Unix(),
	"aud": "test-audience",
	"iat": time.Now().Unix(),
	"nbf": time.Now().Unix(),
	"iss": "test-issuer",
}

func (t *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)

	return token.SignedString([]byte(testSecret))
}

func (t *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	})
}
