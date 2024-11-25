package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrChatNameIsEmpty = errors.New("chat name is empty")
)

type Chat struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
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
