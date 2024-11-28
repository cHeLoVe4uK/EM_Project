package memory

import (
	"context"
	"sync"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	chatrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo"
)

type Repository struct {
	chats map[string]models.Chat
	mu    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		chats: map[string]models.Chat{},
		mu:    sync.RWMutex{},
	}
}

func (r *Repository) GetChatByID(_ context.Context, chatID string) (models.Chat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chat, ok := r.chats[chatID]
	if !ok {
		return models.Chat{}, chatrepo.ErrChatNotFound
	}

	return chat, nil
}

func (r *Repository) CreateChat(_ context.Context, chat models.Chat) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.chats[chat.ID] = chat

	return chat.ID, nil
}

func (r *Repository) UpdateChat(_ context.Context, chat models.Chat) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.chats[chat.ID] = chat

	return nil
}

func (r *Repository) DeleteChat(_ context.Context, chatID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.chats, chatID)

	return nil
}
