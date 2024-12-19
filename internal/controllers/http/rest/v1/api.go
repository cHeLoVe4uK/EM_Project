package v1

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/labstack/echo/v4"
)

type ChatService interface {
	GetAllChats(ctx context.Context) ([]models.Chat, error)
	GetActiveChats(ctx context.Context) ([]models.Chat, error)
	ConnectByID(w http.ResponseWriter, r *http.Request, chatID string, user *models.User) error
	CreateChat(ctx context.Context, chat models.Chat) (string, error)

	GetMessages(ctx context.Context, chatID string) ([]models.Message, error)
	UpdateMessage(ctx context.Context, msg models.Message) error
	DeleteMessage(ctx context.Context, msg models.Message) error
}

type UserService interface {
	Register(ctx context.Context, user models.User) (string, error)
	Login(ctx context.Context, user models.User) (models.Tokens, error)
}

type AuthService interface {
	GetTokens(ctx context.Context, user models.User) (models.Tokens, error)
	Authenticate(ctx context.Context, tokens models.Tokens) (models.Claims, error)
}

type API struct {
	chatService ChatService
	userService UserService
	authService AuthService
}

func NewAPI(
	chatService ChatService,
	userService UserService,
	authService AuthService,
) *API {
	return &API{
		chatService: chatService,
		userService: userService,
		authService: authService,
	}
}

func (a *API) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")

	e := echo.New()

	addr := net.JoinHostPort(host, port)

	a.routes(e)

	srv := http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	slog.Info("server started", slog.String("addr", addr))

	return srv.ListenAndServe()
}
