package mongo

import (
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID     primitive.ObjectID `bson:"_id"`
	ChatID string             `bson:"chat_id"`
	Name   string             `bson:"name"`
}

func FromChat(chat models.Chat) Chat {
	return Chat{
		ID:     primitive.NewObjectID(),
		ChatID: chat.ID,
		Name:   chat.Name,
	}
}

func ToChat(chat Chat) models.Chat {
	return models.Chat{
		ID:   chat.ChatID,
		Name: chat.Name,
	}
}
