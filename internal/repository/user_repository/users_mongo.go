package user_repository

import (
	"context"
	"errors"
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

var (
	ErrUserNotFound = errors.New("user with this id not found")
)

type UsersRepo struct {
	collection *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		collection: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, bson.M{"email": user.Email, "username": user.Username, "password": user.Password})
	return err
}

func (r *UsersRepo) UpdateUser(ctx context.Context, user *models.User) error {
	objectID, _ := primitive.ObjectIDFromHex(user.ID)

	updateFields := bson.M{}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	if user.Username != "" {
		updateFields["username"] = user.Username
	}
	if user.Password != "" {
		updateFields["password"] = user.Password
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updateFields})
	return err
}

func (r *UsersRepo) DeleteUser(ctx context.Context, userID string) error {
	objectID, _ := primitive.ObjectIDFromHex(userID)
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UsersRepo) CheckUserByID(ctx context.Context, userID string) (bool, error) {
	objectID, _ := primitive.ObjectIDFromHex(userID)
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
