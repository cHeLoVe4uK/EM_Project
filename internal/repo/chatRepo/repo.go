package chatrepo

import "errors"

var (
	ErrInvalidChatID = errors.New("invalid chat id")
	ErrChatNotFound  = errors.New("chat not found")
)
