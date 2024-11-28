package schemas

type RequestConnectToChat struct {
	ChatID string `json:"chat_id"`
}

type RequestCreateChat struct {
	Name string `json:"name"`
}

type ResponseCreateChat struct {
	ChatID string `json:"chat_id"`
}

type RequestDeleteChat struct {
	ChatID string `json:"chat_id"`
}
