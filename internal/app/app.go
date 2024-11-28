package app

import (
	"context"

	v1 "github.com/cHeLoVe4uK/EM_Project/internal/controllers/http/ws/v1"
	"github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo/memory"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/chat"
)

type App struct {
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	// Init App

	return a, nil
}

func (a *App) Run() error {

	chatRepo := memory.New()
	chatService := chat.NewService(context.Background(), nil, chatRepo)

	api := v1.NewAPI(chatService)

	return api.Run()
}
