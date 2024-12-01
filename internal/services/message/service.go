package message

import (
	"context"
	"log/slog"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type Repository interface {
	SaveMessages(ctx context.Context, msgs []models.Message) error
	GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SaveMessages(ctx context.Context, msgs []models.Message) error {
	if len(msgs) == 0 {
		slog.Debug(
			"no messages to save",
		)
		return nil
	}

	if err := s.repo.SaveMessages(ctx, msgs); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error) {

	msgs, err := s.repo.GetChatMessages(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
