package memory

import (
	"context"
	"sync"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type Repository struct {
	messages map[string][]models.Message
	mu       sync.RWMutex
}

func New() *Repository {
	return &Repository{
		messages: map[string][]models.Message{},
		mu:       sync.RWMutex{},
	}
}

func (r *Repository) SaveMessages(ctx context.Context, msgs []models.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	chatID := msgs[0].ChatID

	chat, ok := r.messages[chatID]
	if !ok {
		r.messages[chatID] = msgs
		return nil
	}

	r.messages[chatID] = append(chat, msgs...)

	return nil
}

func (r *Repository) GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chat, ok := r.messages[chatID]
	if !ok {
	}

	msgsCount := len(chat)

	if msgsCount > 100 {
		return chat[msgsCount-100:], nil
	}

	return chat, nil
}
