package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateUserIndexes(ctx context.Context, db *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": "text"},
		Options: options.Index().SetUnique(true).SetName("email_idx"),
	}

	_, err := db.Indexes().CreateOne(ctx, indexModel)
	return err
}
