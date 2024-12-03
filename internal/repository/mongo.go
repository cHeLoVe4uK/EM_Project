package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createIndexes(ctx context.Context, db *mongo.Database) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("email_"),
	}

	_, err := db.Collection("users").Indexes().CreateOne(ctx, indexModel)
	return err
}
