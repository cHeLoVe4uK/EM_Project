package models

import "time"

type Message struct {
	ID        string
	Author    string
	ChatID    string
	Content   string
	IsEdited  bool
	CreatedAt time.Time
}
