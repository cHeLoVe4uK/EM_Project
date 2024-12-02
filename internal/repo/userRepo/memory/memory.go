package memory

import (
	"context"
	"sync"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	userrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo"
)

type Repository struct {
	users map[string]models.User
	mu    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		users: map[string]models.User{},
		mu:    sync.RWMutex{},
	}
}

func (r *Repository) GetUserByEmail(_ context.Context, email string) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return models.User{}, userrepo.ErrUserNotFound
}

func (r *Repository) CreateUser(_ context.Context, user models.User) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == user.Email {
			return "", userrepo.ErrUserAlreadyExists
		}
		if u.Username == user.Username {
			return "", userrepo.ErrUserAlreadyExists
		}
	}

	r.users[user.ID] = user

	return user.ID, nil
}
