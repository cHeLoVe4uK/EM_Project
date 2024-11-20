package app

import (
	"context"
	"fmt"
	"os"

	"github.com/cHeLoVe4uK/EM_Project/pkg/connect"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	mongoConnection *mongo.Client
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
		a.initMongo,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initMongo(ctx context.Context) error {
	mongoURL := os.Getenv("MONGO_DSN")
	if mongoURL == "" {
		return fmt.Errorf("MONGO_DSN is not set")
	}

	mongoConnection, err := connect.NewMongo(ctx, mongoURL)
	if err != nil {
		return err
	}

	a.mongoConnection = mongoConnection

	return nil
}

func (a *App) Run() error {

	fmt.Println("Running app...")

	return nil
}
