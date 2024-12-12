package user

import (
	"context"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
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
