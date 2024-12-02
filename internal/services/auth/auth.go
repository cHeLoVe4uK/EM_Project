package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/pkg/tokens"
	"github.com/meraiku/logging"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetTokens(ctx context.Context, user models.User) (models.Token, error) {
	log := logging.L(ctx)

	log.Debug("generating access token")

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
