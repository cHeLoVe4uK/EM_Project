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

type API struct {
	chatService ChatService
}

func NewAPI(chatService ChatService) *API {
	return &API{
		chatService: chatService,
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
