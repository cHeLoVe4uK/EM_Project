package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

/* Структура для хранения payload токенов */
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Service struct {
	salt      string // соль для создания токенов (хранится в .env)
	accessExp int    // время в часах, через которое протухнет access (также хранится в .env)
}

func NewService(salt string, accessExp int) *Service {
	return &Service{salt: salt, accessExp: accessExp}
}

/* Функция создает пару {Accesss, Refresh} для заданного User (временно только Accesss) */
func (s *Service) GetTokens(ctx context.Context, user models.User) (models.Tokens, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Chat App",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.accessExp))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessTokenString, err := s.makeAccessToken(claims)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens := models.Tokens{
		AccessToken: accessTokenString,
	}

	return tokens, nil
}

/* Делает рефреш пары {Access, Refresh} для заданного User (временно не рефрешит) */
func (s *Service) Refresh(ctx context.Context, user models.User) (models.Tokens, error) {
	return s.GetTokens(ctx, user)
}

/* Проверяет пару {Access, Refresh} на валидность */
func (s *Service) Authenticate(ctx context.Context, tokens models.Tokens) (models.Claims, error) {
	claims := Claims{}

	_, err := jwt.ParseWithClaims(tokens.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.salt), nil
	})

	if err != nil {
		return models.Claims{}, err
	}

	out := models.Claims{
		claims.UserID,
		claims.Username,
	}

	return out, nil
}

/* Создает Access */
func (s *Service) makeAccessToken(claims Claims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessTokenString, err := accessToken.SignedString([]byte(s.salt))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}
