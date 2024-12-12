package user

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHashPassword = errors.New("server error")
	ErrUserExists   = errors.New("user already exists")
)

// Регистрация пользователя
func (us *UserService) Register(ctx context.Context, u models.User) (string, error) {
	// Проверка наличия пользователя в БД
	_, err := us.userRepo.CheckUserByEmail(ctx, u.Username)
	if err == nil {
		return "", ErrUserExists
	}

	// Если пользователя не существует создаем его в БД, хэшируя пароль
	passHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashPassword
	}

	u.Password = string(passHash)

	err = us.userRepo.CreateUser(ctx, u)
	if err != nil {
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
