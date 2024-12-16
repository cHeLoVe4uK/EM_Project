package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/repository/user_repository"
	"github.com/meraiku/logging"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
)

// Вход пользователя
func (us *UserService) Login(ctx context.Context, u models.User) (models.Tokens, error) {

	log := logging.L(ctx)

	log.Debug("check user in DB by email")

	// Проверка наличия пользователя в БД
	user, err := us.userRepo.CheckUserByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, user_repository.ErrUserNotFound) {
			log.Warn("user not found", logging.Err(err))

			return models.Tokens{}, ErrUserNotFound
		}
		log.Error("check user in DB by email", logging.Err(err))

		return models.Tokens{}, fmt.Errorf("check user in DB by email: %w", err)
	}

	log = logging.WithAttrs(
		ctx,
		logging.String("user_id", user.ID),
		logging.String("username", user.Username),
	)

	ctx = logging.ContextWithLogger(ctx, log)

	log.Debug("check user password")

	// Если есть, сверяем пароли и выбиваем токены
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return models.Tokens{}, ErrInvalidPassword
	}

	log.Debug("get tokens")

	tokens, err := us.authService.GetTokens(ctx, user)
	if err != nil {
		log.Error("get tokens", logging.Err(err))

		return models.Tokens{}, fmt.Errorf("get tokens: %w", err)
	}

	return tokens, nil
}

// Пока непонятно что тут будет происходить
func (us *UserService) Logout(ctx context.Context, u *models.User) error {
	return nil
}
