package mongo

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	userrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	userCollection = "users"
)

type Repository struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) (*Repository, error) {
	collection := db.Collection(userCollection)

	index := mongo.IndexModel{
		Keys:    bson.M{"email": "text"},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return nil, err
	}

	return &Repository{
		collection: collection,
	}, nil
}

func (r *Repository) CreateUser(ctx context.Context, user models.User) (string, error) {
	repoUser := FromUser(user)

	filter := bson.M{"username": repoUser.Username}

	if err := r.collection.FindOne(ctx, filter).Err(); err == nil {
		return "", userrepo.ErrUserAlreadyExists
	}

	_, err := r.collection.InsertOne(ctx, repoUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", userrepo.ErrUserAlreadyExists
		}
		return "", err
	}

	return user.ID, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user User

	filter := bson.M{"email": email}

	if err := r.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, userrepo.ErrUserNotFound
		}

		return models.User{}, err
	}

	return ToUser(user), nil
}
