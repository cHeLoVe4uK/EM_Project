package v1

import "time"

type Chat struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateChatRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateChatResponse struct {
	ID string `json:"id"`
}

type Message struct {
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	ChatID     string    `json:"chat_id"`
	Content    string    `json:"content"`
	IsEdited   bool      `json:"is_edited"`
	CreatedAt  time.Time `json:"created_at"`
}
