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
