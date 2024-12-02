package user

import (
	"context"
	"fmt"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/meraiku/logging"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = fmt.Errorf("invalid password")
)

type AuthService interface {
	GetTokens(ctx context.Context, user models.User) (models.Token, error)
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (string, error)
}

type Service struct {
	authService AuthService
	repo        Repository
}

func NewService(authService AuthService, repo Repository) *Service {
	return &Service{
		authService: authService,
		repo:        repo,
	}
}

func (s *Service) Login(ctx context.Context, user models.User) (models.Token, error) {
	log := logging.L(ctx)

	log.Debug("getting user by email")

	out, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return models.Token{}, fmt.Errorf("get user by email: %w", err)
	}

	log = logging.WithAttrs(
		ctx,
		logging.String("user_id", out.ID),
	)

	log.Debug("checking password")

	if err := bcrypt.CompareHashAndPassword([]byte(out.Password), []byte(user.Password)); err != nil {
		return models.Token{}, ErrInvalidPassword
	}

	log.Debug("getting tokens")

	token, err := s.authService.GetTokens(ctx, out)
	if err != nil {
		return models.Token{}, fmt.Errorf("get tokens: %w", err)
	}

	log.Debug("user logged in")

	return token, nil
}

func (s *Service) Create(ctx context.Context, user models.User) (string, error) {
	log := logging.L(ctx)

	log.Debug("hashing password")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user.Password = string(passHash)

	log.Debug("creating user")

	out, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	log.Debug("user created", logging.String("user_id", out))

	return out, nil
}
