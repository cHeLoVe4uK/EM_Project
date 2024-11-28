package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/cHeLoVe4uK/EM_Project/internal/app"
)

func main() {
	ctx := context.TODO()

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(l)

	app, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
