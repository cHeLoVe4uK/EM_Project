package v1

type CreateChatResponse struct {
	ID string `json:"id"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type GetActiveChatsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
