package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrChatNameIsEmpty = errors.New("chat name is empty")
)

type Chat struct {
	ID   string
	Name string
}

func NewChat(name string) (Chat, error) {
	if name == "" {
		return Chat{}, ErrChatNameIsEmpty
	}

	return Chat{
		ID:   uuid.NewString(),
		Name: name,
	}, nil
}
