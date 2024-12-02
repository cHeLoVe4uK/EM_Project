package v1

import "time"

type Chat struct {
	ID   string `json:"id" example:"UUID"`
	Name string `json:"name" example:"Best chat name!"`
}

type CreateChatRequest struct {
	Name string `json:"name" example:"Best chat name!" validate:"required"`
}

type CreateChatResponse struct {
	ID string `json:"id" example:"UUID"`
}

type Message struct {
	ID         string    `json:"id" example:"UUID"`
	AuthorID   string    `json:"author_id" example:"UUID"`
	AuthorName string    `json:"author_name" example:"Username"`
	ChatID     string    `json:"chat_id" example:"UUID"`
	Content    string    `json:"content" example:"Hello world!"`
	IsEdited   bool      `json:"is_edited" example:"false"`
	CreatedAt  time.Time `json:"created_at" example:"2022-05-01T00:00:00Z"`
}
