package app

import (
	"context"
	"fmt"
)

type App struct {
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	// Init App

	return a, nil
}

func (a *App) Run() error {

	fmt.Println("Running app...")

	return nil
}
