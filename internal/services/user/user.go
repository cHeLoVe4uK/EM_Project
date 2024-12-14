package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/repository/user_repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHashPassword    = errors.New("password hashing error")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserRepo interface {
	CreateUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, userID string) error
	CheckUserByEmail(ctx context.Context, email string) (models.User, error)
	CheckUserByID(ctx context.Context, userID string) error
}

type AuthService interface {
	GetTokens(context.Context, models.User) (models.Tokens, error)
	Refresh(context.Context, models.User) (models.Tokens, error)
	Authenticate(context.Context, models.Tokens) (models.Claims, error)
}

// Инстанс сервиса для работы с пользователями
type UserService struct {
	userRepo    UserRepo
	authService AuthService
}

// Конструктор сервиса
func NewUserService(repo UserRepo, service AuthService) *UserService {
	return &UserService{
		userRepo:    repo,
		authService: service,
	}
}

// Регистрация пользователя
func (us *UserService) Register(ctx context.Context, u models.User) (string, error) {
	// Проверка наличия пользователя в БД
	_, err := us.userRepo.CheckUserByEmail(ctx, u.Email)
	if err == nil {
		slog.Error("From user service - register user. Error occured while exist userRepo.CheckUserByEmail:", ErrUserExists)
		return "", ErrUserExists
	}
	if !errors.Is(err, user_repository.ErrUserNotFound) {
		slog.Error("From user service - register user. Error occured while exist userRepo.CheckUserByEmail:", err)
		return "", err
	}

	// Если пользователя не существует создаем его в БД, хэшируя пароль
	passHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("From user service - register user. Error occured while exist bcrypt.GenerateFromPassword:", ErrHashPassword)
		return "", ErrHashPassword
	}

	u.Password = string(passHash)

	err = us.userRepo.CreateUser(ctx, u)
	if err != nil {
		slog.Error("From user service - register user. Error occured while exist userRepo.CreateUser:", err)
		return "", err
	}

	return u.ID, nil
}

// Обновление пользователя
func (us *UserService) UpdateUser(ctx context.Context, u models.User) error {
	// Проверка наличия пользователя в БД
	err := us.userRepo.CheckUserByID(ctx, u.ID)
	if err != nil {
		slog.Error("From user service - update user. Error occured while exist userRepo.CheckUserByID:", err)
		return err
	}

	// Если найден обновляем
	err = us.userRepo.UpdateUser(ctx, u)
	if err != nil {
		slog.Error("From user service - update user. Error occured while exist userRepo.UpdateUser:", err)
		return err
	}

	return nil
}

// Удаление пользователя
func (us *UserService) DeleteUser(ctx context.Context, u models.User) error {
	// Проверка наличия пользователя в БД
	err := us.userRepo.CheckUserByID(ctx, u.ID)
	if err != nil {
		slog.Error("From user service - delete user. Error occured while exist userRepo.CheckUserByID:", err)
		return err
	}

	// Если найден удаляем
	err = us.userRepo.DeleteUser(ctx, u.ID)
	if err != nil {
		slog.Error("From user service - delete user. Error occured while exist userRepo.DeleteUser:", err)
		return err
	}

	return nil
}

// Вход пользователя
func (us *UserService) Login(ctx context.Context, u models.User) (models.Tokens, error) {
	// Проверка наличия пользователя в БД
	user, err := us.userRepo.CheckUserByEmail(ctx, u.Email)
	if err != nil {
		slog.Error("From user service - login user. Error occured while exist userRepo.CheckUserByEmail:", err)
		return models.Tokens{}, err
	}

	// Если есть, сверяем пароли и выбиваем токены
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		slog.Error("From user service - login user. Error occured while exist bcrypt.CompareHashAndPassword:", ErrInvalidPassword)
		return models.Tokens{}, ErrInvalidPassword
	}

	tokens, err := us.authService.GetTokens(ctx, user)
	if err != nil {
		slog.Error("From user service - login user. Error occured while exist authService.GetTokens:", err)
		return models.Tokens{}, err
	}

	return tokens, nil
}

// Пока непонятно что тут будет происходить
func (us *UserService) Logout(ctx context.Context, u models.User) error {
	return nil
}
