package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthtenticator(secret string, aud string, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud}

}

func (auth JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(auth.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
func (auth JWTAuthenticator) ValidatorToken(claims string) (*jwt.Token, error) {
	return nil, nil

}
