package user_repository

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	usersCollection = "users"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrInvalidUserID  = errors.New("invalid user id")
	ErrUserNotFound   = errors.New("user with this id not found")
)

type UsersRepo struct {
	collection *mongo.Collection
}

func NewUsersRepo(ctx context.Context, db *mongo.Database) (*UsersRepo, error) {
	uc := db.Collection(usersCollection)

	_ = uc.Indexes().DropAll(ctx)

	return &UsersRepo{
		collection: uc,
	}, nil
}

func (r *UsersRepo) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicateEmail
	}

	return nil
}

func (r *UsersRepo) UpdateUser(ctx context.Context, user models.User) error {
	u := FromUser(user)

	updateFields := bson.M{}
	if user.Email != "" {
		updateFields["email"] = u.Email
	}
	if user.Username != "" {
		updateFields["username"] = u.Username
	}
	if user.Password != "" {
		updateFields["password"] = u.Password
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"id": u.UserID}, bson.M{"$set": updateFields})
	return err
}

func (r *UsersRepo) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrInvalidUserID
	}

	_, err := r.collection.DeleteOne(ctx, bson.M{"id": userID})
	return err
}

func (r *UsersRepo) CheckUserByID(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrInvalidUserID
	}

	if err := r.collection.FindOne(ctx, bson.M{"id": userID}).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (r *UsersRepo) CheckUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user User

	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	out := ToUser(user)

	return out, nil
}
