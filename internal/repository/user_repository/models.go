package user_repository

import (
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type User struct {
	UserID   string `bson:"id"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func FromUser(user models.User) User {
	return User{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}
}

func ToUser(user User) models.User {
	return models.User{
		ID:       user.UserID,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}
}
