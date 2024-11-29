package v1

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")

	addr := net.JoinHostPort(host, port)

	srv := http.Server{
		Addr:    addr,
		Handler: a.routes(),
	}

	slog.Info("server started", slog.String("addr", addr))

	return srv.ListenAndServe()
}
