package connect

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongo(ctx context.Context, mongoDSN string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoDSN))
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	go func() {
		<-ctx.Done()

		slog.Debug("Disconnecting mongo")

		if err := client.Disconnect(ctx); err != nil {
			slog.Warn("failed to disconnect mongo", "error", err)
		}

		slog.Debug("Disconnected from mongo successfully")
	}()

	return client, nil
}
