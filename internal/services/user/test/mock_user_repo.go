package user_test

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/repository/user_repository"
)

var ErrWithDB = errors.New("trouble with db")

type userRepo map[string]models.User

// Имитация репозитория (работы с бд) (UserRepo)
type Mock struct {
	userRepo
}

// Конструктор
func New() *Mock {
	return &Mock{userRepo: userRepo{}}
}

// Создать пользователя
func (m *Mock) CreateUser(ctx context.Context, user models.User) error {
	if m.userRepo == nil {
		return ErrWithDB
	}

	m.userRepo[user.ID] = user
	return nil
}

// Обновить пользователя
func (m *Mock) UpdateUser(ctx context.Context, user models.User) error {
	if _, ok := m.userRepo[user.ID]; !ok {
		return user_repository.ErrUserNotFound
	}

	m.userRepo[user.ID] = user
	return nil
}

// Удалить пользователя
func (m *Mock) DeleteUser(ctx context.Context, userID string) error {
	if _, ok := m.userRepo[userID]; !ok {
		return user_repository.ErrUserNotFound
	}

	delete(m.userRepo, userID)
	return nil
}

// Найти пользователя по email
func (m *Mock) CheckUserByEmail(ctx context.Context, email string) (models.User, error) {
	if m.userRepo == nil {
		return models.User{}, ErrWithDB
	}

	for _, u := range m.userRepo {
		if u.Email == email {
			return u, nil
		}
	}

	return models.User{}, user_repository.ErrUserNotFound
}

// Найти пользователя по ID
func (m *Mock) CheckUserByID(ctx context.Context, userID string) error {
	if m.userRepo == nil {
		return ErrWithDB
	}

	if _, ok := m.userRepo[userID]; !ok {
		return user_repository.ErrUserNotFound
	}

	return nil
}
