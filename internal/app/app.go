package app

import (
	"context"

	v1 "github.com/cHeLoVe4uK/EM_Project/internal/controllers/http/ws/v1"
	chatrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo/memory"
	userrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo/memory"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/auth"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/chat"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/user"
)

type App struct {
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	// Init App

	return a, nil
}

func (a *App) Run() error {

	chatRepo := chatrepo.New()
	chatService := chat.NewService(context.Background(), nil, chatRepo)

	authService := auth.NewService()

	userRepo := userrepo.New()
	userService := user.NewService(authService, userRepo)

	api := v1.NewAPI(chatService, userService)

	return api.Run()
}
