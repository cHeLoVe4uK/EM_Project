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
	ErrHashPassword = errors.New("server error")
	ErrUserExists   = errors.New("user already exists")
)

// Регистрация пользователя
func (us *UserService) Register(ctx context.Context, u models.User) (string, error) {
	log := logging.WithAttrs(
		ctx,
		logging.String("user_id", u.ID),
	)

	ctx = logging.ContextWithLogger(ctx, log)

	log.Debug("check user in DB by email")
	// Проверка наличия пользователя в БД
	_, err := us.userRepo.CheckUserByEmail(ctx, u.Email)
	switch err {
	case nil:
		return "", ErrUserExists
	case user_repository.ErrUserNotFound:
		break
	default:
		log.Error("check user in DB by email", logging.Err(err))
		return "", fmt.Errorf("check user in DB by email: %w", err)
	}

	log.Debug("hash user password")

	// Если пользователя не существует создаем его в БД, хэшируя пароль
	passHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("hash user password", logging.Err(err))

		return "", ErrHashPassword
	}

	u.Password = string(passHash)

	log.Debug("create user in DB")

	err = us.userRepo.CreateUser(ctx, u)
	if err != nil {
		log.Error("create user in DB", logging.Err(err))

		return "", err
	}

	return u.ID, nil
}

// Обновление пользователя
func (us *UserService) UpdateUser(ctx context.Context, u models.User) error {
	// Проверка наличия пользователя в БД
	err := us.userRepo.CheckUserByID(ctx, u.ID)
	if err != nil {
		return err
	}

	// Если найден обновляем
	err = us.userRepo.UpdateUser(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

// Удаление пользователя
func (us *UserService) DeleteUser(ctx context.Context, u *models.User) error {
	// Проверка наличия пользователя в БД
	err := us.userRepo.CheckUserByID(ctx, u.ID)
	if err != nil {
		return err
	}

	// Если найден удаляем
	err = us.userRepo.DeleteUser(ctx, u.ID)
	if err != nil {
		return err
	}
	return nil
}
