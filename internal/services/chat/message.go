package chat

import (
	"encoding/json"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"
)

// Модель сообщения в рамках чата
type MessageDTO struct {
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	ChatID     string    `json:"chat_id"`
	Content    string    `json:"content"`
	IsEdited   bool      `json:"is_edited"`
	CreatedAt  time.Time `json:"created_at"`
}

// Создаёт новую модель сообщения, наполняет поля значениями
func NewMessage(client *Client, text string) *MessageDTO {
	return &MessageDTO{
		ID:         uuid.NewString(),
		AuthorID:   client.ID,
		AuthorName: client.Username,
		ChatID:     client.ChatRoom.ID,
		Content:    text,
		IsEdited:   false,
		CreatedAt:  time.Now().UTC(),
	}
}

// Преобразует модель сообщения в json
func (msg *MessageDTO) Render() ([]byte, error) {

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Преобразует сервис модель в модель чата
func FromMessage(msg models.Message) MessageDTO {
	return MessageDTO{
		ID:         msg.ID,
		AuthorID:   msg.AuthorID,
		AuthorName: msg.AuthorName,
		ChatID:     msg.ChatID,
		Content:    msg.Content,
		IsEdited:   msg.IsEdited,
		CreatedAt:  msg.CreatedAt,
	}
}

// Преобразует сервис модели в модели чата
func FromMessageBatch(msgs []models.Message) []MessageDTO {
	messages := make([]MessageDTO, len(msgs))

	for i, msg := range msgs {
		messages[i] = FromMessage(msg)
	}

	return messages
}

// Преобразует модель чата в модель сервиса
func ToMessage(msg MessageDTO) models.Message {
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

// Преобразует модели чата в модели сервиса
func ToMessageBatch(msgs []MessageDTO) []models.Message {
	messages := make([]models.Message, len(msgs))

	for i, msg := range msgs {
		messages[i] = ToMessage(msg)
	}

	return messages
}
