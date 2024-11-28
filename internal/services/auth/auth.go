package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/pkg/tokens"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetTokens(ctx context.Context, user models.User) (models.Token, error) {

	accessToken, err := tokens.GenerateJWT(
		user.ID,
		user.Username,
		time.Hour,
		[]byte("secret"),
	)
	if err != nil {
		return models.Token{}, fmt.Errorf("generate access token: %w", err)
	}

	t := models.Token{
		Token: accessToken,
	}

	return t, nil
}
