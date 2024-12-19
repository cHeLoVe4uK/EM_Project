package message

import (
	"context"
	"errors"
	"fmt"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/repository/msg_repository"
	"github.com/meraiku/logging"
)

var (
	ErrMessageNotFound = errors.New("message not found")
	ErrNotAllowed      = errors.New("not allowed")
)

// Интерфейс репозитория сообщений
type Repository interface {
	SaveMessages(ctx context.Context, msgs []models.Message) error
	GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error)
	Update(ctx context.Context, msg models.Message) error
	Delete(ctx context.Context, msg models.Message) error
}

// Сервис сообщений
type Service struct {
	repo Repository
}

// Конструктор сервиса
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Сохранение массива сообщений в репозиторий
func (s *Service) SaveMessages(ctx context.Context, msgs []models.Message) error {
	log := logging.L(ctx)

	if len(msgs) == 0 {
		log.Debug("no messages to save, skipping")

		return nil
	}

	log.Debug("save messages", logging.Int("messages_count", len(msgs)))

	if err := s.repo.SaveMessages(ctx, msgs); err != nil {
		log.Error("save messages", logging.Err(err))

		return fmt.Errorf("save messages: %w", err)
	}

	log.Debug("messages saved")

	return nil
}

// Получение массива сообщений из репозитория по айди чата. Возвращает слайс до 100 последних сообщений
func (s *Service) GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error) {
	log := logging.L(ctx)

	log.Debug("get chat messages from repository")

	msgs, err := s.repo.GetChatMessages(ctx, chatID)
	if err != nil {
		log.Error("get chat messages", logging.Err(err))

		return nil, fmt.Errorf("get chat messages: %w", err)
	}

	log.Debug(
		"got messages",
		logging.Int("message_count", len(msgs)),
	)

	return msgs, nil
}

func (s *Service) UpdateMessageContent(ctx context.Context, msg models.Message) error {
	log := logging.L(ctx)

	log.Debug("update message status")

	msg.IsEdited = true

	log.Debug("update message content")

	err := s.repo.Update(ctx, msg)
	switch err {
	case msg_repository.ErrMessageNotFound:
		log.Warn("update message", logging.Err(err))

		return ErrMessageNotFound
	case msg_repository.ErrNotAllowed:
		log.Warn("update message", logging.Err(err))

		return ErrNotAllowed
	case nil:
		log.Debug("message updated")

		return nil
	default:
		log.Error("update message", logging.Err(err))

		return err
	}
}

func (s *Service) DeleteMessage(ctx context.Context, msg models.Message) error {
	log := logging.L(ctx)

	log.Debug("delete message")

	err := s.repo.Delete(ctx, msg)
	switch err {
	case msg_repository.ErrMessageNotFound:
		log.Warn("delete message", logging.Err(err))

		return ErrMessageNotFound
	case msg_repository.ErrNotAllowed:
		log.Warn("delete message", logging.Err(err))

		return ErrNotAllowed
	case nil:
		log.Debug("message deleted")

		return nil
	default:
		log.Error("delete message", logging.Err(err))

		return err
	}
}
