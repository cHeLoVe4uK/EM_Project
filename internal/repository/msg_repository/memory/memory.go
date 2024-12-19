package memory

import (
	"context"
	"sync"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/repository/msg_repository"
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
		return nil, nil
	}

	msgsCount := len(chat)

	if msgsCount > 100 {
		return chat[msgsCount-100:], nil
	}

	return chat, nil
}

func (r *Repository) Update(ctx context.Context, msg models.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, ok := r.messages[msg.ChatID]
	if !ok {
		return nil
	}

	for i, m := range chat {
		if m.ID == msg.ID {
			r.messages[msg.ChatID][i] = msg
			return nil
		}
	}

	return msg_repository.ErrMessageNotFound
}

func (r *Repository) Delete(ctx context.Context, msg models.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, ok := r.messages[msg.ChatID]
	if !ok {
		return nil
	}

	for i, m := range chat {
		if m.ID == msg.ID {
			r.messages[msg.ChatID] = append(chat[:i], chat[i+1:]...)
			return nil
		}
	}

	return msg_repository.ErrMessageNotFound
}
