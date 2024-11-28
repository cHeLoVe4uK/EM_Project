package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type ChatService interface {
	ConnectByID(w http.ResponseWriter, r *http.Request, chatID string, user *models.User) error
	CreateChat(ctx context.Context, chat models.Chat) (string, error)
}

type UserService interface {
	Create(ctx context.Context, user models.User) (string, error)
	Login(ctx context.Context, user models.User) (models.Token, error)
}

type API struct {
	chatService ChatService
	userService UserService
}

func NewAPI(chatService ChatService, userService UserService) *API {
	return &API{
		chatService: chatService,
		userService: userService,
	}
}

func (a *API) Run() error {
	srv := http.Server{
		Addr:    ":8080",
		Handler: a.routes(),
	}

	slog.Info("server started", slog.String("addr", ":8080"))

	return srv.ListenAndServe()
}
