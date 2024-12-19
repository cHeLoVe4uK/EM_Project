package chat_repository

import (
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type Chat struct {
	ChatID string `bson:"chat_id"`
	Name   string `bson:"name"`
}

func FromChat(chat models.Chat) Chat {
	return Chat{
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

func ToChatBatch(chats []Chat) []models.Chat {
	out := make([]models.Chat, len(chats))

	for i, chat := range chats {
		out[i] = ToChat(chat)
	}

	return out
}
