package v1

type CreateChatRequest struct {
	Name string `json:"name"`
}

type JoinChatRequest struct {
	ChatID string `json:"chat_id"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
