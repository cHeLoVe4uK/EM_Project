package user

import (
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

var (
	UserNotFound          = errors.New("User not found")
	UserAlreadyRegistered = errors.New("User already registered")
)

// Метод для регистрации пользователя
func (us *UserService) Register(u *models.User) error {
	// Для начала нужно проверить есть ли пользователь в БД
	ok, err := us.userRepo.CheckUserByID(u.ID)
	if err != nil {
		return err
	}
	if ok {
		return UserAlreadyRegistered
	}

	// Если нет создаем
	err = us.userRepo.Create(u)
	if err != nil {
		return err
	}
	return nil
}

// Метод для обновления пользователя
func (us *UserService) UpdateUser(u *models.User) error {
	// Для начала нужно проверить есть ли пользователь в БД
	ok, err := us.userRepo.CheckUserByID(u.ID)
	if err != nil {
		return err
	}
	if !ok {
		return UserNotFound
	}

	// Если найден обновляем
	err = us.userRepo.Create(u)
	if err != nil {
		return err
	}
	return nil
}

// Метод для удаления пользователя
func (us *UserService) DeleteUser(u *models.User) error {
	// Для начала нужно проверить есть ли пользователь в БД
	ok, err := us.userRepo.CheckUserByID(u.ID)
	if err != nil {
		return err
	}
	if !ok {
		return UserNotFound
	}

	// Если найден удаляем
	err = us.userRepo.Delete(u)
	if err != nil {
		return err
	}
	return nil
}
