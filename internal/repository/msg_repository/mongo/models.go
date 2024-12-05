package mongo

import (
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type Message struct {
	ID         string    `bson:"id"`
	AuthorID   string    `bson:"author_id"`
	AuthorName string    `bson:"author_name"`
	ChatID     string    `bson:"chat_id"`
	Content    string    `bson:"content"`
	IsEdited   bool      `bson:"is_edited"`
	CreatedAt  time.Time `bson:"created_at"`
}

func FromMessage(msg models.Message) Message {
	return Message{
		ID:         msg.ID,
		AuthorID:   msg.AuthorID,
		AuthorName: msg.AuthorName,
		ChatID:     msg.ChatID,
		Content:    msg.Content,
		IsEdited:   msg.IsEdited,
		CreatedAt:  msg.CreatedAt,
	}
}

func FromMessageBatch(msgs []models.Message) []Message {
	out := make([]Message, len(msgs))

	for i, msg := range msgs {
		out[i] = FromMessage(msg)
	}

	return out
}

func ToMessage(msg Message) models.Message {
	return models.Message{
		ID:         msg.ID,
		AuthorID:   msg.AuthorID,
		AuthorName: msg.AuthorName,
		ChatID:     msg.ChatID,
		Content:    msg.Content,
		IsEdited:   msg.IsEdited,
		CreatedAt:  msg.CreatedAt,
	}
}

func ToMessageBatch(msgs []Message) []models.Message {
	out := make([]models.Message, len(msgs))

	for i, msg := range msgs {
		out[i] = ToMessage(msg)
	}

	return out
}
