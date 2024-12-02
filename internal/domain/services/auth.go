package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetTokens(ctx context.Context, user models.User) (models.Tokens, error) {
	claims := models.Claims{
		UserID:   user.UserID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Chat App",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return models.Tokens{}, err
	}

	tokens := models.Tokens{
		AccessToken: accessTokenString,
	}

	return tokens, nil
}

func (s *Service) Refresh(ctx context.Context, user models.User) (models.Tokens, error) { // временно просто вызывает создание нового access токена, так как никакой логики проверки refresh токена пока нет
	return s.GetTokens(ctx, user)
}

func (s *Service) Authenticate(ctx context.Context, tokens models.Tokens) (models.Claims, error) {
	claims := models.Claims{}
	secretKey := os.Getenv("SECRET_KEY")

	_, err := jwt.ParseWithClaims(tokens.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return models.Claims{}, err
	}
	return claims, nil
}
