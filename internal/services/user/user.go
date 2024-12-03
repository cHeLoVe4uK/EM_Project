package user

import (
	"context"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type UserRepo interface {
	CreateUser(context.Context, *models.User) error
	UpdateUser(context.Context, *models.User) error
	DeleteUser(context.Context, *models.User) error
	CheckUserByEmail(context.Context, string) (*models.User, error)
	CheckUserByID(context.Context, string) error
}

type AuthService interface {
	GetTokens(*models.User) (*models.Token, error)
	RefreshTokens(*models.Token) (*models.Token, error)
	Authorization(*models.Token) error
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
