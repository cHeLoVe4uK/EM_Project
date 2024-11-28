package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoSecret = errors.New("empty jwt secret")
)

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(
	id string,
	username string,
	ttl time.Duration,
	secret []byte,
) (string, error) {

	if len(secret) == 0 {
		return "", ErrNoSecret
	}

	c := &Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl).UTC()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	token, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseJWT(tokenStr string, secret []byte) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if token != nil {
		if token.Valid {
			return claims, nil
		}
	}

	return nil, err
}
