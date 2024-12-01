package chat

import (
	"encoding/json"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"
)

type MessageDTO struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	ChatID    string    `json:"chat_id"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"is_edited"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessage(client *Client, text string) *MessageDTO {
	return &MessageDTO{
		ID:        uuid.NewString(),
		Author:    client.ID,
		ChatID:    client.ChatRoom.ID,
		Content:   text,
		IsEdited:  false,
		CreatedAt: time.Now().UTC(),
	}
}

func (msg *MessageDTO) Render() ([]byte, error) {

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func FromMessage(msg models.Message) MessageDTO {
	return MessageDTO{
		ID:        msg.ID,
		Author:    msg.Author,
		ChatID:    msg.ChatID,
		Content:   msg.Content,
		IsEdited:  msg.IsEdited,
		CreatedAt: msg.CreatedAt,
	}
}

func FromMessageBatch(msgs []models.Message) []MessageDTO {
	messages := make([]MessageDTO, len(msgs))

	for i, msg := range msgs {
		messages[i] = FromMessage(msg)
	}

	return messages
}

func ToMessage(msg MessageDTO) models.Message {
	return models.Message{
		ID:        msg.ID,
		Author:    msg.Author,
		ChatID:    msg.ChatID,
		Content:   msg.Content,
		IsEdited:  msg.IsEdited,
		CreatedAt: msg.CreatedAt,
	}
}

func ToMessageBatch(msgs []MessageDTO) []models.Message {
	messages := make([]models.Message, len(msgs))

	for i, msg := range msgs {
		messages[i] = ToMessage(msg)
	}

	return messages
}
