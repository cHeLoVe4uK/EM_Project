package user

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

// Вход пользователя
func (us *UserService) Login(ctx context.Context, u *models.User) (*models.Token, error) {
	// Проверка наличия пользователя в БД
	user, err := us.userRepo.CheckUserByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	// Если есть, сверяем пароли и выбиваем токены
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	tokens, err := us.authService.GetTokens(u)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Пока непонятно что тут будет происходить
func (us *UserService) Logout(ctx context.Context, u *models.User) error {
	return nil
}
