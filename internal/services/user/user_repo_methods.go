package user

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

var (
	ErrUserNotFound          = errors.New("User not found")
	ErrUserAlreadyRegistered = errors.New("User already registered")
)

// Регистрация пользователя
func (us *UserService) Register(ctx context.Context, u *models.User) error {
	// Проверка наличия пользователя в БД
	ok, err := us.userRepo.CheckUserByUsername(ctx, u.Username)
	if err != nil {
		return err
	}
	if ok {
		return ErrUserAlreadyRegistered
	}

	// Если нет создаем
	err = us.userRepo.CreateUser(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

// Обновление пользователя
func (us *UserService) UpdateUser(ctx context.Context, u *models.User) error {
	// Проверка наличия пользователя в БД
	ok, err := us.userRepo.CheckUserByID(ctx, u.Username)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
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
	ok, err := us.userRepo.CheckUserByID(ctx, u.Username)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	// Если найден удаляем
	err = us.userRepo.DeleteUser(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
