package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidatorToken(token string) (*jwt.Token, error)
}
