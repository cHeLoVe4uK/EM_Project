package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
		ID:   primitive.NewObjectID().String(),
		Name: name,
	}, nil
}
