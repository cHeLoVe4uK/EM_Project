package v1

import "time"

type CreateChatResponse struct {
	ID string `json:"id"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type Chat struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	ChatID    string    `json:"chat_id"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"is_edited"`
	CreatedAt time.Time `json:"created_at"`
}
