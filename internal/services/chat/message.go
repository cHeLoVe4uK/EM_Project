package chat

type MessageDTO struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	ChatID    string `json:"chat_id"`
	Content   string `json:"content"`
	IsEdited  bool   `json:"is_edited"`
	Timestamp string `json:"timestamp"`
}
