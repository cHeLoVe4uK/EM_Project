package chat

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageDTO struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	ChatID    string    `json:"chat_id"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"is_edited"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(client *Client, text string) *MessageDTO {
	return &MessageDTO{
		ID:        uuid.NewString(),
		Author:    client.ID,
		ChatID:    client.ChatRoom.ID,
		Content:   text,
		IsEdited:  false,
		Timestamp: time.Now(),
	}
}

func (msg *MessageDTO) Render() ([]byte, error) {

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return data, nil
}
