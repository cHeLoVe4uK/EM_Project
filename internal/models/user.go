package models

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

var (
	ErrInvalidData = errors.New("invalid data")
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func NewUser(
	email string,
	username string,
	password string,
) (User, error) {

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return User{}, fmt.Errorf("invalid email: %w", err)
	}

	if username == "" || password == "" {
		return User{}, ErrInvalidData
	}

	u := User{
		ID:       uuid.NewString(),
		Email:    addr.String(),
		Username: username,
		Password: password,
	}

	return u, nil
}
