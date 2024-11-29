package mongo

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	userrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	userCollection = "users"
)

type Repository struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) *Repository {
	collection := db.Collection(userCollection)

	return &Repository{
		collection: collection,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user models.User) (string, error) {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user User

	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, userrepo.ErrUserNotFound
		}

		return models.User{}, err
	}

	return ToUser(user), nil
}
