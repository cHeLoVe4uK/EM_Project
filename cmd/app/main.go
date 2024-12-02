package main

import (
	"context"
	"log"

	_ "github.com/cHeLoVe4uK/EM_Project/api/swagger"
	"github.com/cHeLoVe4uK/EM_Project/internal/app"
)

func main() {
	ctx := context.TODO()

	app, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
