package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidData = errors.New("invalid data")
)

type User struct {
	ID       string
	Email    string
	Username string
	Password string
}

func NewUser(
	email string,
	username string,
	password string,
) (User, error) {

	if email == "" || username == "" || password == "" {
		return User{}, ErrInvalidData
	}

	u := User{
		ID:       uuid.NewString(),
		Email:    email,
		Username: username,
		Password: password,
	}

	return u, nil
}
