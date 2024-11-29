package mongo

import (
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   string             `bson:"user_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func FromUser(user models.User) User {
	return User{
		ID:       primitive.NewObjectID(),
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func ToUser(user User) models.User {
	return models.User{
		ID:       user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}
