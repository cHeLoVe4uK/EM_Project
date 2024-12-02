package message

import (
	"context"
	"fmt"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/meraiku/logging"
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
	log := logging.L(ctx)

	if len(msgs) == 0 {
		log.Debug("no messages to save, skipping")
		return nil
	}

	log.Debug(
		"saving messages",
		logging.Int("messages_count", len(msgs)),
	)

	if err := s.repo.SaveMessages(ctx, msgs); err != nil {
		return fmt.Errorf("save messages: %w", err)
	}

	log.Debug("messages saved")

	return nil
}

func (s *Service) GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error) {
	log := logging.L(ctx)

	log.Debug("getting messages from repository")

	msgs, err := s.repo.GetChatMessages(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat messages: %w", err)
	}

	log.Debug(
		"got messages",
		logging.Int("message_count", len(msgs)),
	)

	return msgs, nil
}
