package services

import (
	"fmt"
	"os"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetToken(userID string, username string) (string, error) {
	claims := models.Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Chat App",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Refresh(claims models.Claims) (string, error) { // временно просто вызывает создание нового access токена, так как никакой логики проверки refresh токена пока нет
	return GetToken(claims.UserID, claims.Username)
}

func Authenticate(token string) (*models.Claims, error) {
	claims := &models.Claims{}
	secretKey := os.Getenv("SECRET_KEY")

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}
	return claims, nil
}
