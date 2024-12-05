package user

import (
	"context"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/auth"
)

type UserRepo interface {
	CreateUser(context.Context, *models.User) error
	UpdateUser(context.Context, *models.User) error
	DeleteUser(context.Context, *models.User) error
	CheckUserByEmail(context.Context, string) (*models.User, error)
	CheckUserByID(context.Context, string) error
}

type AuthService interface {
	GetTokens(context.Context, models.User) (models.Tokens, error)
	Refresh(context.Context, models.User) (models.Tokens, error)
	Authenticate(context.Context, models.Tokens) (auth.Claims, error)
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
