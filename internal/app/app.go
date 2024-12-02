package app

import (
	"context"
	"os"

	"github.com/cHeLoVe4uK/EM_Project/internal/config"
	v1 "github.com/cHeLoVe4uK/EM_Project/internal/controllers/http/rest/v1"
	chatRepoMemory "github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo/memory"
	chatRepoMongo "github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo/mongo"
	msgRepoMemory "github.com/cHeLoVe4uK/EM_Project/internal/repo/msgRepo/memory"
	msgRepoMongo "github.com/cHeLoVe4uK/EM_Project/internal/repo/msgRepo/mongo"
	userRepoMemory "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo/memory"
	userRepoMongo "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo/mongo"

	"github.com/cHeLoVe4uK/EM_Project/internal/services/auth"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/chat"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/message"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/user"
	"github.com/cHeLoVe4uK/EM_Project/pkg/connect"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	mongo    *mongo.Client
	chatRepo chat.ChatRepository
	msgRepo  message.Repository
	userRepo user.Repository
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {

	deps := []func(context.Context) error{
		a.initConfig,
		a.initRepos,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {

	config.Load()

	return nil
}

func (a *App) initRepos(ctx context.Context) error {

	repoType := os.Getenv("REPO_TYPE")

	switch repoType {
	case "mongo":

		if err := a.initMongo(ctx); err != nil {
			return err

		}
		db := a.mongo.Database("em_chat")

		chatRepo, err := chatRepoMongo.New(db)
		if err != nil {
			return err
		}
		a.chatRepo = chatRepo

		userRepo, err := userRepoMongo.New(db)
		if err != nil {
			return err
		}
		a.userRepo = userRepo

		msgRepo, err := msgRepoMongo.New(db)
		if err != nil {
			return err
		}
		a.msgRepo = msgRepo

	default:
		a.chatRepo = chatRepoMemory.New()
		a.userRepo = userRepoMemory.New()
		a.msgRepo = msgRepoMemory.New()
	}

	return nil
}

func (a *App) initMongo(ctx context.Context) error {

	mongoDSN := os.Getenv("MONGO_DSN")

	client, err := connect.NewMongo(ctx, mongoDSN)
	if err != nil {
		return err
	}

	a.mongo = client

	return nil
}

func (a *App) Run() error {

	msgService := message.New(a.msgRepo)

	chatService := chat.NewService(context.Background(), msgService, a.chatRepo)

	authService := auth.NewService()

	userService := user.NewService(authService, a.userRepo)

	api := v1.NewAPI(chatService, userService)

	return api.Run()
}
