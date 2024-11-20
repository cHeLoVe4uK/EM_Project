package connect

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongo(ctx context.Context, mongoDSN string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoDSN))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("failed to disconnect mongodb: %v", err)
		}
	}()

	return client, nil
}
