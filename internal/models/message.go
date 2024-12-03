package models

import "time"

type Message struct {
	ID         string
	AuthorID   string
	AuthorName string
	ChatID     string
	Content    string
	IsEdited   bool
	CreatedAt  time.Time
}
