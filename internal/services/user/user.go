package user

import (
	"context"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type UserRepo interface {
	CreateUser(context.Context, *models.User) error
	UpdateUser(context.Context, *models.User) error
	DeleteUser(context.Context, *models.User) error
	CheckUserByUsername(context.Context, string) (bool, error)
	CheckUserByID(context.Context, string) (bool, error)
}

type AuthService interface {
	GetTokens(*models.User) (string, string, error)
	RefreshTokens(string, string) (string, string, error)
	Authorization(string) error
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
