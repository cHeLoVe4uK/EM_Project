package v1

type CreateChatRequest struct {
	Name string `json:"name"`
}

type JoinChatRequest struct {
	ChatID string `json:"chat_id"`
}
